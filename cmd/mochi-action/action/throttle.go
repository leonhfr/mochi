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

func NewThrottledTransport(requests int, period time.Duration) *ThrottledTransport {
	return &ThrottledTransport{
		transporter: http.DefaultTransport,
		limiter:     rate.NewLimiter(rate.Every(period), requests),
	}
}

func (tt *ThrottledTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if err := tt.limiter.Wait(r.Context()); err != nil {
		return nil, err
	}
	return tt.transporter.RoundTrip(r)
}
