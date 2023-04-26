package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	MaxConcurrentRequests = 8
	MaxRetries            = 5
	RequestTimeout        = 5 * time.Second
	BatchSize             = 20
	OutputFileName        = "results.txt"
)

func main() {
	var wg sync.WaitGroup
	sem := make(chan struct{}, MaxConcurrentRequests)
	results := make(chan string, BatchSize)

	go writeResultsToFile(results, &wg)

	for _, site := range sites {
		wg.Add(1)
		go getSite(sem, results, &wg, site)
	}

	wg.Wait()

	wg.Add(1)
	close(results)

	wg.Wait()
}

func getSite(sem chan struct{}, results chan<- string, wg *sync.WaitGroup, site string) {
	defer wg.Done()

	sem <- struct{}{}
	defer func() { <-sem }()

	status, err := checkSiteStatus(site)
	if err != nil {
		fmt.Printf("error: %s", err)

		return
	}

	results <- fmt.Sprintf("%s %d\n", site, status)
}

func checkSiteStatus(site string) (int, error) {
	client := http.Client{
		Timeout: RequestTimeout,
	}

	for i := 0; i < MaxRetries; i++ {
		resp, err := client.Get(site)
		if err != nil {
			fmt.Printf("failed to get site: %s error: %s\n", site, err)

			continue
		}

		resp.Body.Close()
		return resp.StatusCode, nil
	}

	return 0, fmt.Errorf("failed to get site: %s\n", site)
}

func writeResultsToFile(results <-chan string, wg *sync.WaitGroup) {
	var batch []string
	file, err := os.Create(OutputFileName)
	if err != nil {
		fmt.Printf("error creating file: %v \n", err)
		return
	}

	defer file.Close()
	for result := range results {
		batch = append(batch, result)

		if len(batch) >= BatchSize {
			writeBatch(file, batch)
			batch = nil
		}
	}

	if len(batch) > 0 {
		writeBatch(file, batch)
	}

	wg.Done()
}

func writeBatch(file *os.File, batch []string) {
	data := strings.Join(batch, "")
	if _, err := file.WriteString(data); err != nil {
		fmt.Printf("error writing to file: %v\n", err)
	}
}
