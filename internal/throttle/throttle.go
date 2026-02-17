package throttle

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

var _ http.RoundTripper = &Transport{}

// Transport is a throttled transport that implements the http.RoundTripper interface.
type Transport struct {
	rt      http.RoundTripper
	limiter *rate.Limiter
	sem     chan struct{}
}

// Option is a Transport option.
type Option func(*Transport)

// New allows events up to rate 1 event/interval and permits bursts of at most burst tokens.
func New(interval time.Duration, burst int, options ...Option) *Transport {
	transport := &Transport{
		rt:      http.DefaultTransport,
		limiter: rate.NewLimiter(rate.Every(interval), burst),
	}
	for _, option := range options {
		option(transport)
	}
	return transport
}

// WithTransport sets a http.RoundTripper that replaces http.DefaultTransport.
func WithTransport(rt http.RoundTripper) Option {
	return func(t *Transport) {
		t.rt = rt
	}
}

// WithConcurrencyLimit sets the maximum number of concurrent in-flight requests.
// A value lower than 1 disables concurrency limiting.
func WithConcurrencyLimit(limit int) Option {
	return func(t *Transport) {
		if limit < 1 {
			t.sem = nil
			return
		}
		t.sem = make(chan struct{}, limit)
	}
}

// RoundTrip implements the http.RoundTripper interface.
func (tt *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	if err := tt.limiter.Wait(r.Context()); err != nil {
		return nil, err
	}

	release, err := tt.acquire(r.Context())
	if err != nil {
		return nil, err
	}
	defer release()

	return tt.rt.RoundTrip(r)
}

func (tt *Transport) acquire(ctx context.Context) (func(), error) {
	if tt.sem == nil {
		return func() {}, nil
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case tt.sem <- struct{}{}:
		return func() { <-tt.sem }, nil
	}
}
