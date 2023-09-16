package sync

import (
	"context"
	"sync"

	"github.com/leonhfr/mochi/filesystem"
)

func processJobMap(ctx context.Context, jobs jobMap, numHandlers int, scr *syncCardResult, client Client, fs filesystem.Interface) error {
	done := make(chan struct{})
	defer close(done)

	jobc := newJobChannel(jobs, done)
	reqc := make(chan cardRequest)
	errc := make(chan error)

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
			reqHandler(ctx, scr, client, done, reqc, errc)
		}()
	}

	go func() {
		wgJobs.Wait()
		close(reqc)
	}()

	go func() {
		wgReqs.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			return err
		}
	}

	return nil
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

func reqHandler(ctx context.Context, scr *syncCardResult, client Client, done <-chan struct{}, reqc <-chan cardRequest, errc chan<- error) {
	for req := range reqc {
		select {
		case errc <- req.do(ctx, scr, client):
		case <-done:
			return
		}
	}
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
