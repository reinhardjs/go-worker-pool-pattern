package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Job Payload
type Job struct {
	ID        int
	Payload   string
}

var (
	jobQueue   = make(chan Job, 100) // Buffer 100
	maxWorkers = 3
)

func worker(id int, jobs <-chan Job) {
	for job := range jobs {
		fmt.Printf("[Worker %d] Processing Job %d: %s\n", id, job.ID, job.Payload)
		time.Sleep(2 * time.Second) // Simulasi proses berat
		fmt.Printf("[Worker %d] Finished Job %d\n", id, job.ID)
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Non-blocking send dengan select
	select {
	case jobQueue <- Job{ID: int(time.Now().Unix()), Payload: "User Upload"}:
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Request diterima. Sedang diproses di background."))
	default:
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Server sibuk. Coba lagi nanti."))
	}
}

func main() {
	// Start Workers (Background Daemon)
	for w := 1; w <= maxWorkers; w++ {
		go worker(w, jobQueue)
	}

	http.HandleFunc("/process", requestHandler)
	
	fmt.Println("Server running on :8080 (Daemon Pool Mode)")
	http.ListenAndServe(":8080", nil)
}