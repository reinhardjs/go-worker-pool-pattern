# Go Concurrency Patterns: Worker Pool

This repository demonstrates the evolution of the **Worker Pool Pattern** in Go (Golang). It progresses from a simple, fixed-size goroutine pool to complex, production-ready HTTP implementations using Fan-In/Fan-Out strategies and Unbounded Queues.

These examples are designed to illustrate how to manage limited resources (CPU/RAM) efficiently while handling high-volume workloads.

## 📂 Repository Structure

The code is organized into five distinct levels of complexity:

| Folder | Pattern Name | Use Case |
| :--- | :--- | :--- |
| **`01-basic-pool`** | Fixed Worker Pool | Simple background tasks (logs, simple calc). |
| **`02-robust-pool`** | Pool with Results | ETL jobs, image processing where output is needed. |
| **`03-http-fire-and-forget`** | Daemon / Async Pool | HTTP requests that trigger long-running background tasks. |
| **`04-http-fan-in-fan-out`** | Aggregator Pattern | Dashboard APIs requiring data from multiple sources simultaneously. |
| **`05-unlimited-queue`** | **Unbounded Queue** | Critical ingestion where **dropping requests is impossible** and submission must never block. |

---

## 🚀 Getting Started

### Prerequisites
* Go 1.18+ (due to standard library usage).

### How to Run
Navigate to the root directory and run the specific folder's `main.go` file:

```bash
go run 01-basic-pool/main.go
```

## 📦 Production-Ready Libraries

For production systems, consider using these battle-tested solutions instead of building from scratch:

* [**panjf2000/ants**](https://github.com/panjf2000/ants) - High-performance, low-memory footprint with auto-scaling capabilities.
* [**gammazero/workerpool**](https://github.com/gammazero/workerpool) - Simple, idiomatic implementation with built-in backpressure and graceful shutdown.
