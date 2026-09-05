package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/akhnafal-aban/asset_tracker_be/internal/config"
	"github.com/akhnafal-aban/asset_tracker_be/internal/database"
	"github.com/akhnafal-aban/asset_tracker_be/internal/handler"
	"github.com/akhnafal-aban/asset_tracker_be/internal/repository"
	"github.com/akhnafal-aban/asset_tracker_be/internal/server"
	"github.com/akhnafal-aban/asset_tracker_be/internal/service"
)

func main() {
	// 1. Setup structured logging
	// We use JSON handler in production, Text handler in dev.
	// For simplicity, let's use TextHandler for now.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// 2. Load Configuration
	cfg := config.Load()

	// 3. Initialize Database
	db, err := database.InitDB(cfg.DBPath)
	if err != nil {
		slog.Error("Failed to initialize database", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	// 4. Initialize Dependencies (Repositories, Services, Handlers)
	assetRepo := repository.NewAssetRepository(db)
	assetSvc := service.NewAssetService(assetRepo)
	assetHandler := handler.NewAssetHandler(assetSvc)

	// 5. Initialize Server (Constructor DI)
	srv := server.New(cfg, assetHandler)

	// 4. Start Server with context
	ctx := context.Background()
	if err := srv.Start(ctx); err != nil {
		slog.Error("Failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}
