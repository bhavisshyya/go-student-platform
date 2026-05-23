package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bhavisshyya/students-api/internal/config"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// database setup

	// setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to API"))
	})

	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("Server Started: ", slog.String("Addr", cfg.HTTPServer.Addr))
	// fmt.Printf("Server Started %s", cfg.HTTPServer.Addr)

	done := make(chan os.Signal, 1)

	// interupt signals
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")
		}
	}()

	<-done

	// gracefull Shutdown
	slog.Info("shutting down server")
	ctx, cancl := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancl()

	err := server.Shutdown(ctx)

	if err != nil {
		slog.Error("Fail to shurdown server", slog.String("error:", err.Error()))
	}
}
