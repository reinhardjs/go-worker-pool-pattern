package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID int
}

// Global Channels
var (
	// Channel Masuk (Input) - Unbuffered is fine because the dispatcher reads instantly
	inputChan = make(chan Job) 
	// Channel Keluar (Output ke Worker)
	workerChan = make(chan Job)
)

// 1. THE DISPATCHER (The Magic of "Infinite" Queue)
// Goroutine ini bertugas memindahkan data dari Input ke Slice (jika worker sibuk)
// atau langsung ke Worker (jika worker idle).
func dispatcher() {
	// Slice sebagai antrean tak terbatas (in-memory)
	var queue []Job

	for {
		// Tentukan channel mana yang aktif untuk pengiriman
		var activeWorkerChan chan Job
		var nextJob Job

		// Jika ada antrean di memori, kita coba kirim job terdepan ke worker
		if len(queue) > 0 {
			activeWorkerChan = workerChan
			nextJob = queue[0]
		}

		select {
		// KASUS A: Ada Job baru masuk
		case job := <-inputChan:
			// Terima job dan masukkan ke antrean memori (TIDAK PERNAH BLOK)
			queue = append(queue, job)
			fmt.Printf("[Queue] Job %d buffered. Queue Size: %d\n", job.ID, len(queue))

		// KASUS B: Worker siap menerima job (dan ada job di antrean)
		case activeWorkerChan <- nextJob:
			// Hapus job dari antrean memori karena sudah diambil worker
			queue = queue[1:]
			fmt.Printf("[Queue] Job %d sent to worker. Queue Size: %d\n", nextJob.ID, len(queue))
		}
	}
}

// 2. THE WORKER (Standard)
func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range workerChan {
		fmt.Printf("   [Worker %d] Processing Job %d...\n", id, job.ID)
		time.Sleep(2 * time.Second) // Simulasi kerja lambat
		fmt.Printf("   [Worker %d] Done Job %d\n", id, job.ID)
	}
}

func main() {
	// Start Dispatcher (Infinite Queue Manager)
	go dispatcher()

	// Start Workers (Limited Concurrency)
	var wg sync.WaitGroup
	numWorkers := 3
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, &wg)
	}

	// 3. SUBMITTER (Flooding Tasks)
	// Kita kirim 20 tugas secepat kilat.
	// Karena dispatcher pakai slice, pengirim TIDAK AKAN BLOKir meski worker lambat.
	go func() {
		for i := 1; i <= 20; i++ {
			fmt.Printf("-> Submitting Job %d\n", i)
			inputChan <- Job{ID: i} // Ini return instantly!
			time.Sleep(50 * time.Millisecond) // Kirim cepat
		}
		// Note: Di real app, kita butuh mekanisme shutdown yang lebih rapi untuk dispatcher
	}()

	// Biarkan main thread jalan cukup lama untuk melihat proses
	time.Sleep(15 * time.Second)
}
