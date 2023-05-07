package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	count = 100

	timeUnit     = time.Second
	timeDuration = 5
	rate         = 10
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
	rate   int
	bucket chan struct{}
	ticker *time.Ticker
}

func newRateLimiter(unit time.Duration, duration int, rate int) *rateLimiter {
	bucket := make(chan struct{}, rate)
	rl := &rateLimiter{rate: rate, bucket: bucket}
	rl.refill()

	go func() {
		ticker := time.NewTicker(time.Duration(duration) * unit)
		for range ticker.C {
			rl.refill()
		}
	}()

	return rl
}

func (rl *rateLimiter) refill() {
	for i := 0; i < rl.rate; i++ {
		rl.bucket <- struct{}{}
	}
}

func (rl *rateLimiter) wait() {
	<-rl.bucket
}
