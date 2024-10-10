package throttle

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

var _ http.RoundTripper = &Transport{}

// Transport is a throttled transport that implements the http.RoundTripper interface.
type Transport struct {
	rt      http.RoundTripper
	limiter *rate.Limiter
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

// RoundTrip implements the http.RoundTripper interface.
func (tt *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	if err := tt.limiter.Wait(r.Context()); err != nil {
		return nil, err
	}
	return tt.rt.RoundTrip(r)
}
