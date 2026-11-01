package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielxfeng/hsl-watch/internal/config"
	"github.com/danielxfeng/hsl-watch/internal/hsl"
	"github.com/danielxfeng/hsl-watch/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	if len(cfg.ArrivalIDs) == 0 {
		log.Print("warning: ARRIVAL_IDS is empty, nothing will be fetched")
	}

	client := hsl.NewClient(cfg.HSLAPIKey)
	arrivals := service.NewArrivalService(client, cfg.ArrivalIDs, cfg.Debug)

	service.Schedule(ctx, arrivals.UpdateArrivals)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("GET /api/v1/arrivals", func(w http.ResponseWriter, r *http.Request) {
		data := arrivals.GetArrivals()
		if data == nil {
			http.Error(w, "no data yet", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("write response: %v", err)
		}
	})

	srv := &http.Server{Addr: ":8080", Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	shutdownDone := make(chan struct{})
	go func() {
		defer close(shutdownDone)
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown: %v", err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("http server: %v", err)
		return
	}
	<-shutdownDone
}
