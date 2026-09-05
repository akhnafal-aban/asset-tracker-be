package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akhnafal-aban/asset_tracker_be/internal/config"
	"github.com/akhnafal-aban/asset_tracker_be/internal/handler"
	"github.com/akhnafal-aban/asset_tracker_be/internal/middleware"
)

type Server struct {
	config     *config.Config
	httpServer *http.Server
}

func New(cfg *config.Config, assetHandler *handler.AssetHandler) *Server {
	mux := http.NewServeMux()

	if assetHandler != nil {
		mux.HandleFunc("GET /api/v1/assets", assetHandler.ListAssets)
		mux.HandleFunc("POST /api/v1/assets", assetHandler.CreateAsset)
		mux.HandleFunc("GET /api/v1/assets/{id}", assetHandler.GetAsset)
		mux.HandleFunc("PUT /api/v1/assets/{id}", assetHandler.UpdateAsset)
		mux.HandleFunc("DELETE /api/v1/assets/{id}", assetHandler.DeleteAsset)
		mux.HandleFunc("GET /api/v1/categories", assetHandler.ListCategories)
	}

	h := middleware.Logger(mux)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: h,
		// Good practice timeouts
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return &Server{
		config:     cfg,
		httpServer: srv,
	}
}

// Start runs the HTTP server and handles graceful shutdown.
func (s *Server) Start(ctx context.Context) error {
	// Channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		slog.Info("Starting server", slog.String("port", s.config.Port), slog.String("env", s.config.Env))
		serverErrors <- s.httpServer.ListenAndServe()
	}()

	// Channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received or an error occurs.
	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server error: %w", err)
		}
		return nil

	case sig := <-shutdown:
		slog.Info("Shutdown signal received", slog.String("signal", sig.String()))

		// Create context with timeout for graceful shutdown
		shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()

		// Ask server to shutdown gracefully
		err := s.httpServer.Shutdown(shutdownCtx)
		if err != nil {
			// graceful shutdown failed, forceful close
			slog.Error("Graceful shutdown failed, forcing close", slog.Any("error", err))
			if err = s.httpServer.Close(); err != nil {
				return fmt.Errorf("force close error: %w", err)
			}
			return err
		}

		slog.Info("Server stopped gracefully")
		return nil
	}
}
