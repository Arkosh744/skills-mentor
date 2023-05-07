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

	timeUnit           = time.Second
	timeDuration       = 5
	rate         int64 = 10
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
	rate         int64
	counter      int64
	lock         sync.Mutex
	resetTimer   *time.Timer
	resetChannel chan struct{}
}

func newRateLimiter(unit time.Duration, duration int, rate int64) *rateLimiter {
	resetChannel := make(chan struct{})
	resetTimer := time.NewTimer(unit * time.Duration(duration))

	go func() {
		for {
			<-resetTimer.C
			resetChannel <- struct{}{}
			resetTimer.Reset(unit * time.Duration(duration))
		}
	}()

	return &rateLimiter{rate: rate, resetTimer: resetTimer, resetChannel: resetChannel}
}

func (rl *rateLimiter) wait() {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	if rl.counter >= rl.rate {
		<-rl.resetChannel
		atomic.SwapInt64(&rl.counter, 0)
	}

	atomic.AddInt64(&rl.counter, 1)
}
