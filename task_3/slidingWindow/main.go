package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	limit       = 10
	durInterval = 1 * time.Second
	count       = 100
)

func main() {
	now := time.Now()
	resultCh := make(chan int)
	var wg sync.WaitGroup

	go func() {
		for value := range resultCh {
			fmt.Println(value)
		}
	}()

	rLimiter := NewSlidingWindow(limit, durInterval)

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

			for !rLimiter.Allow() {
				time.Sleep(10 * time.Millisecond)
			}

			resultCh <- RPCCall(i)
		}

		reqCh <- task
	}

	wg.Wait()
	close(resultCh)
	close(reqCh)

	fmt.Println(time.Since(now))
}

func RPCCall(i int) int {
	fmt.Println(i)
	return rand.Int()
}

type rateLimiter struct {
	limit    int
	interval time.Duration
	mu       sync.Mutex

	currentTime time.Time

	prevCount    int
	currentCount int
}

func NewSlidingWindow(limit int, interval time.Duration) *rateLimiter {
	return &rateLimiter{
		limit:    limit,
		interval: interval,
	}

}

func (l *rateLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	newPeriod := l.currentTime.Add(l.interval)
	if time.Now().After(newPeriod) {
		l.currentTime = time.Now()
		l.prevCount = l.currentCount
		l.currentCount = 0
	}

	elapsed := time.Now().Sub(l.currentTime).Seconds()
	prevCount := float64(l.prevCount)
	currentCount := float64(l.currentCount)
	interval := l.interval.Seconds()

	curCount := (prevCount*(interval-elapsed) + currentCount) / interval
	if curCount > float64(l.limit) {
		return false
	}

	l.currentCount++

	return true
}
