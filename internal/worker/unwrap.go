package worker

import "sync"

// Result wraps a Data with an error.
type Result[Data any] struct {
	data Data
	err  error
}

// Unwrap forwards all errors to the errC channel. If the error is nil,
// the value is written to the returned channel.
//
// When in is closed, the function releases its counter in the wait group.
func Unwrap[Data any](wg *sync.WaitGroup, in <-chan Result[Data], errC chan<- error) <-chan Data {
	out := make(chan Data, cap(in))
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(out)

		for result := range in {
			if result.err != nil {
				errC <- result.err
			} else {
				out <- result.data
			}
		}
	}()
	return out
}
