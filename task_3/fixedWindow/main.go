package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	count = 100

	timeUnit     = time.Second
	timeDuration = 5
	rate         = 2
)

func main() {
	resultCh := make(chan int)
	var wg sync.WaitGroup

	go func() {
		for value := range resultCh {
			fmt.Println(value)
		}
	}()

	rLimiter := newRateLimiter(timeUnit, timeDuration, rate)

	// store requests in order and then execute them
	reqCh := make(chan func())
	go func() {
		for task := range reqCh {
			task()
		}
	}()

	for i := 0; i < count; i++ {
		wg.Add(1)
		task := func() {
			defer wg.Done()

			rLimiter.wait()
			resultCh <- RPCCall()
		}

		reqCh <- task
	}

	wg.Wait()
	close(resultCh)
	close(reqCh)
}

func RPCCall() int {
	return rand.Int()
}

type rateLimiter struct {
	rate        int64
	counter     int64
	lock        sync.Mutex
	resetWindow time.Time
	windowSize  time.Duration
	done        chan struct{}
}

func newRateLimiter(unit time.Duration, duration int, rate int64) *rateLimiter {
	windowSize := time.Duration(duration) * unit
	rl := &rateLimiter{
		rate:       rate,
		windowSize: windowSize,
		done:       make(chan struct{}),
	}

	go func() {
		ticker := time.NewTicker(windowSize)
		for range ticker.C {
			rl.lock.Lock()
			atomic.StoreInt64(&rl.counter, 0)
			rl.lock.Unlock()

			close(rl.done)
			rl.done = make(chan struct{})
		}
	}()

	return rl
}

func (rl *rateLimiter) wait() {
	rl.lock.Lock()

	if rl.counter >= rl.rate {
		rl.lock.Unlock()
		<-rl.done
		rl.lock.Lock()
	}

	atomic.AddInt64(&rl.counter, 1)
	rl.lock.Unlock()
}
