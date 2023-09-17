package sync

import (
	"context"
	"sync"

	"github.com/leonhfr/mochi/filesystem"
)

func processJobMap(ctx context.Context, jobs jobMap, numHandlers int, client Client, fs filesystem.Interface) (CardResult, error) {
	done := make(chan struct{})
	defer close(done)

	jobc := newJobChannel(jobs, done)
	reqc := make(chan cardRequest)
	errc := make(chan error)
	resc := make(chan CardResult, numHandlers)

	var wgJobs, wgReqs sync.WaitGroup
	wgJobs.Add(numHandlers)
	wgReqs.Add(numHandlers)

	for i := 0; i < numHandlers; i++ {
		go func() {
			defer wgJobs.Done()
			jobHandler(ctx, client, fs, done, jobc, reqc, errc)
		}()

		go func() {
			defer wgReqs.Done()
			resc <- reqHandler(ctx, client, done, reqc, errc)
		}()
	}

	go func() {
		wgJobs.Wait()
		close(reqc)
	}()

	go func() {
		wgReqs.Wait()
		close(errc)
		close(resc)
	}()

	for err := range errc {
		if err != nil {
			return CardResult{}, err
		}
	}

	var cr CardResult
	for res := range resc {
		cr.Created += res.Created
		cr.Updated += res.Updated
		cr.Archived += res.Archived
	}
	return cr, nil
}

func jobHandler(ctx context.Context, client Client, fs filesystem.Interface, done <-chan struct{}, jobc <-chan *deckJob, reqc chan<- cardRequest, errc chan<- error) {
	for job := range jobc {
		reqs, err := generateCardRequests(ctx, job, client, fs)
		select {
		case errc <- err:
		default:
		}

		for _, req := range reqs {
			select {
			case reqc <- req:
			case <-done:
				return
			}
		}
	}
}

func reqHandler(ctx context.Context, client Client, done <-chan struct{}, reqc <-chan cardRequest, errc chan<- error) CardResult {
	var cr CardResult
	for req := range reqc {
		select {
		case errc <- processCardRequest(ctx, req, client):
			cr.increment(req.kind)
		case <-done:
			return cr
		}
	}
	return cr
}

func newJobChannel(jobs jobMap, done <-chan struct{}) <-chan *deckJob {
	jobc := make(chan *deckJob)
	go func(jobs jobMap) {
		defer close(jobc)
		for _, job := range jobs {
			select {
			case jobc <- job:
			case <-done:
				return
			}
		}
	}(jobs)
	return jobc
}
