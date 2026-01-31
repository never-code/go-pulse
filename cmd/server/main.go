package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"go-pulse/internal/config"
	"go-pulse/internal/worker"
)

func main() {
	cfg := config.Load()

	// Initialize worker pool
	pool := worker.NewPool(cfg.WorkerCount)
	pool.Start()

	// Submit test jobs
	pool.Submit(worker.Job{URL: "https://google.com"})
	pool.Submit(worker.Job{URL: "https://github.com"})
	pool.Submit(worker.Job{URL: "https://httpstat.us/500"})
	pool.Submit(worker.Job{URL: "https://httpstat.us/200?sleep=2000"})

	// HTTP server setup
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	go func() {
		log.Println("🚀 Go-Pulse running on port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("🛑 Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(ctx)
	close(pool.Jobs)

	log.Println("✅ Shutdown complete")
}
