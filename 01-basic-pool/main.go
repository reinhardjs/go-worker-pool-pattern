package main

import (
	"fmt"
	"sync"
	"time"
)

// Worker sederhana yang hanya mencetak data
func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, j)
		time.Sleep(500 * time.Millisecond) // Simulasi kerja
	}
}

func main() {
	const numJobs = 10
	const numWorkers = 3

	jobs := make(chan int, numJobs)
	var wg sync.WaitGroup

	// 1. Start Workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, &wg)
	}

	// 2. Send Jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // Tutup channel agar worker tahu tugas selesai

	// 3. Wait
	wg.Wait()
	fmt.Println("All jobs processed.")
}