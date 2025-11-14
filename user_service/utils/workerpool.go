package utils

import (
	"context"
	"sync"
)

func Workerpool[T, R any](ctx context.Context, in <-chan T, workersNum int, f func(e T) R) <-chan R {
	result := make(chan R)

	wg := &sync.WaitGroup{}

	for range workersNum {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case value, ok := <-in:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case result <- f(value):
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	return result
}
