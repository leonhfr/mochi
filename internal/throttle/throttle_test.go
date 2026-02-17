package throttle

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testRoundTripper struct {
	started chan struct{}
	release chan struct{}
	active  atomic.Int32
	max     atomic.Int32
}

func newTestRoundTripper() *testRoundTripper {
	return &testRoundTripper{
		started: make(chan struct{}, 8),
		release: make(chan struct{}),
	}
}

func (tt *testRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	current := tt.active.Add(1)
	defer tt.active.Add(-1)

	for {
		max := tt.max.Load()
		if current <= max || tt.max.CompareAndSwap(max, current) {
			break
		}
	}

	tt.started <- struct{}{}
	<-tt.release

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("ok")),
	}, nil
}

func Test_WithConcurrencyLimit(t *testing.T) {
	rt := newTestRoundTripper()
	transport := New(time.Nanosecond, 10,
		WithTransport(rt),
		WithConcurrencyLimit(1),
	)

	var wg sync.WaitGroup
	wg.Add(2)
	for range 2 {
		go func() {
			defer wg.Done()
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com", nil)
			require.NoError(t, err)
			resp, err := transport.RoundTrip(req)
			require.NoError(t, err)
			assert.NoError(t, resp.Body.Close())
		}()
	}

	select {
	case <-rt.started:
	case <-time.After(time.Second):
		t.Fatal("first request did not start")
	}

	select {
	case <-rt.started:
		t.Fatal("second request started before first completed")
	case <-time.After(50 * time.Millisecond):
	}

	rt.release <- struct{}{}
	select {
	case <-rt.started:
	case <-time.After(time.Second):
		t.Fatal("second request did not start after first completed")
	}
	rt.release <- struct{}{}

	wg.Wait()
	assert.Equal(t, int32(1), rt.max.Load())
}
