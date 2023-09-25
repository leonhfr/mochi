package action

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type ThrottledTransport struct {
	transporter http.RoundTripper
	limiter     *rate.Limiter
}

// NewThrottledTransport allows events up to rate r and permits bursts of at most t tokens.
func NewThrottledTransport(r time.Duration, t int) *ThrottledTransport {
	return &ThrottledTransport{
		transporter: http.DefaultTransport,
		limiter:     rate.NewLimiter(rate.Every(r), t),
	}
}

func (tt *ThrottledTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if err := tt.limiter.Wait(r.Context()); err != nil {
		return nil, err
	}
	return tt.transporter.RoundTrip(r)
}
