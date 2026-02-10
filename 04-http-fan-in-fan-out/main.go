package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// SubTask adalah pecahan pekerjaan dari satu request besar
type SubTask struct {
	ID     int
	Source string
}

// SubResult adalah hasil dari sub-task
type SubResult struct {
	ID     int
	Status string
}

// Kita menggunakan worker pool global untuk memproses sub-task
var (
	taskQueue  = make(chan func(), 50) // Channel of Functions (Pattern menarik!)
	maxWorkers = 5
)

// Worker generik yang mengeksekusi fungsi apapun yang dikirim ke queue
func worker(id int, tasks <-chan func()) {
	for task := range tasks {
		task() // Eksekusi fungsi closure
	}
}

func init() {
	// Start Global Workers
	for w := 1; w <= maxWorkers; w++ {
		go worker(w, taskQueue)
	}
}

func aggregatorHandler(w http.ResponseWriter, r *http.Request) {
	// Skenario: User minta data dashboard yang butuh fetch dari 5 sumber berbeda
	sources := []string{"Database A", "API B", "Cache C", "File D", "Service E"}
	
	// Channel lokal untuk Fan-In hasil HANYA dari request ini
	resultsCh := make(chan SubResult, len(sources))
	var wg sync.WaitGroup

	// FAN-OUT: Pecah request menjadi 5 sub-tasks
	for i, src := range sources {
		wg.Add(1)
		
		// Capture variable untuk closure
		idx, sourceName := i, src
		
		// Kirim tugas ke Global Worker Pool
		taskQueue <- func() {
			defer wg.Done()
			
			// Simulasi proses (misal: fetch API)
			time.Sleep(500 * time.Millisecond) 
			
			// Kirim hasil ke channel lokal (FAN-IN)
			resultsCh <- SubResult{
				ID:     idx,
				Status: fmt.Sprintf("Data from %s loaded", sourceName),
			}
		}
	}

	// Goroutine untuk menutup channel results setelah semua selesai
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	// Kumpulkan hasil (Aggregasi)
	var finalReport []SubResult
	for res := range resultsCh {
		finalReport = append(finalReport, res)
	}

	// Response JSON ke Client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Dashboard Data Loaded",
		"data":    finalReport,
	})
}

func main() {
	http.HandleFunc("/dashboard", aggregatorHandler)
	fmt.Println("Server running on :8080 (Fan-In/Fan-Out Mode)")
	http.ListenAndServe(":8080", nil)
}