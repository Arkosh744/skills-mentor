package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	mainRace()
}

func mainRace() {
	cache := New(3)
	wg := sync.WaitGroup{}

	cache.Set("name", "Alex", time.Millisecond*500)
	cache.Set("hobby", "BJJ", time.Second*1)
	cache.Set("job", "developer", time.Second*5)

	wg.Add(4)

	go func() {
		defer wg.Done()
		fmt.Println(cache.Get("hobby"))
		fmt.Println(cache.Get("name"))
		fmt.Println(cache.Get("job"))
	}()

	go func() {
		defer wg.Done()
		cache.Delete("hobby")
		fmt.Println(cache.Get("hobby"))
		cache.Set("hobby", "DoDo", time.Second*5)
		fmt.Println(cache.Get("hobby"))
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 1100)
		fmt.Println("_____________________________")
		fmt.Println("Testing ttl:")
		fmt.Println(cache.Get("name"))
		fmt.Println(cache.Get("hobby"))
		fmt.Println(cache.Get("job"))
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 2)
		cache.Set("city", "Moscow", time.Second*5)
		cache.Set("country", "Canada", time.Second*5)
		cache.Set("language", "Go", time.Second*5)
		cache.Set("system", "linux", time.Second*5)
	}()

	wg.Wait()
	fmt.Println("_____________________________")
	fmt.Println("Testing LRU:")
	fmt.Println(cache.Get("name"))
	fmt.Println(cache.Get("hobby"))
	fmt.Println(cache.Get("job"))
	fmt.Println(cache.Get("city"))
	fmt.Println(cache.Get("country"))
	fmt.Println(cache.Get("language"))
	fmt.Println(cache.Get("system"))
}
