## Kategori 1: Tuning & Performa

**Q1: Berapa jumlah worker yang ideal (numWorkers)?**

- **Jawaban Singkat**: Tergantung jenis tugasnya (CPU-Bound vs I/O-Bound).
-  **Penjelasan**:
  -  **CPU-Bound (Enkripsi, Image Processing)**: Gunakan jumlah Core CPU (`runtime.NumCPU()`). Menambah lebih dari itu justru memperlambat karena context _switching_.
  -  **I/O-Bound (Database, API Call, Upload)**: Bisa jauh lebih banyak (misal: 50, 100, atau lebih), karena worker lebih banyak "menunggu" respon daripada menggunakan CPU. Lakukan _benchmark_ untuk angka pastinya.

**Q2: Apa bedanya Worker Pool dengan semaphore (buffered channel kosong)?**

- **Jawaban Singkat**: Worker Pool lebih terstruktur untuk reuse goroutine, Semaphore hanya membatasi konkurensi.
- **Penjelasan**:
  - Kita bisa membatasi konkurensi hanya dengan `make(chan struct{}, 5)`. Tapi Worker Pool (dengan worker tetap) lebih hemat memori karena tidak terus-menerus membuat dan menghancurkan goroutine (mengurangi overhead _Garbage Collection_).

**Q3: Apakah channel buffer size berpengaruh pada performa?**
- **Jawaban Singkat**: Ya, sebagai _backpressure_.
- **Penjelasan**:
  - **Buffer Kecil**: Hemat memori, tapi producer (pengirim tugas) akan lebih sering terblokir jika worker lambat.
  - **Buffer Besar**: Producer jarang terblokir, tapi jika crash, banyak data di antrean yang hilang (potensi data _loss_ tinggi).


## Kategori 2: Error Handling & Reliability

**Q4: Apa yang terjadi jika satu worker `panic`? Apakah aplikasi crash?**

- **Jawaban Singkat**: Ya, seluruh aplikasi akan crash kecuali di-recover.
- **Penjelasan**: Sangat disarankan membungkus kode di dalam worker dengan `defer` dan `recover()`. Jika panic terjadi, catat log error-nya, lalu biarkan worker itu lanjut mengambil tugas berikutnya (atau restart worker tersebut).

**Q5: Bagaimana cara melakukan Graceful Shutdown yang benar?**
- **Jawaban Singkat**: Stop terima tugas -> Tutup channel jobs -> Tunggu worker selesai (`wg.Wait()`).
- **Penjelasan**: Jangan langsung `os.Exit()`.
  1. Stop HTTP Server / Producer.
  2. `close(jobChannel).`
  3. Tunggu semua worker selesai memproses sisa antrean menggunakan `sync.WaitGroup`.

**Q6: Bagaimana kalau antrean penuh dan kita tidak mau blocking (misal di HTTP Handler)?**
- **Jawaban Singkat**: Gunakan `select` dengan `default` case.
- **Penjelasan**:
  ```Go
  select {
  
  case jobQueue <- job:
  
      // Masuk antrean
  
  default:
  
      // Antrean penuh, return error 503 (Service Unavailable)
  
  }
  ```
  Ini mencegah request HTTP hanging selamanya.


## Kategori 3: Advanced Concepts

**Q7: Kapan saya harus pakai library (seperti `ants`) dibanding bikin sendiri?**
- **Jawaban Singkat**: Untuk production skala besar, pakai library.
- **Penjelasan**: Bikin sendiri bagus untuk belajar atau kasus simpel. Tapi library seperti `ants` punya fitur auto-scaling (menambah worker saat sibuk, mengurangi saat sepi) dan manajemen memori yang jauh lebih efisien (goroutine recycling) yang sulit dibuat sendiri dari nol.

**Q8: Bisakah worker pool memproses berbagai jenis tugas (bukan cuma 1 struct)?**
- **Jawaban Singkat**: Bisa, gunakan `interface{}` atau closure `func()`.
- **Penjelasan**: Channel bisa bertipe `chan func()`. Worker tinggal mengeksekusi fungsi tersebut (`job()`). Ini membuat worker pool jadi generik dan bisa mengerjakan apa saja.

**Q9: Apakah Worker Pool menjamin urutan (Ordering)?**
- **Jawaban Singkat**: Tidak.
- **Penjelasan**: Jika Job A masuk duluan daripada Job B, bisa jadi Job B selesai duluan (misal Job A lebih berat atau worker yang menghandle Job B lebih cepat). Jika butuh urutan ketat, Worker Pool **bukan** solusinya; gunakan pemrosesan serial atau shard berdasarkan ID.
