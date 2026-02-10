package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Job struct {
	ID       int
	ImageURL string
}

type Result struct {
	JobID  int
	Status string
	Err    error
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		// Simulasi latency acak
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		
		// Simulasi hasil
		results <- Result{
			JobID:  job.ID,
			Status: fmt.Sprintf("Worker %d processed %s", id, job.ImageURL),
			Err:    nil,
		}
	}
}

func main() {
	const numJobs = 10
	const numWorkers = 3

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)
	var wg sync.WaitGroup

	// 1. Start Workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// 2. Dispatch Jobs (Producer)
	go func() {
		for j := 1; j <= numJobs; j++ {
			jobs <- Job{ID: j, ImageURL: fmt.Sprintf("/img/%d.jpg", j)}
		}
		close(jobs)
	}()

	// 3. Wait & Close Results (Coordinator)
	go func() {
		wg.Wait()
		close(results)
	}()

	// 4. Collect Results (Consumer)
	for res := range results {
		fmt.Printf("Job %d: %s\n", res.JobID, res.Status)
	}
}