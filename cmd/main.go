// main.go is the application entry point. It bootstraps all layers using
// manual dependency injection and starts the HTTP server with graceful shutdown.
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/yourusername/go-skeleton-code/internal/delivery/http/handler"
	"github.com/yourusername/go-skeleton-code/internal/delivery/http/router"
	"github.com/yourusername/go-skeleton-code/internal/infrastructure/config"
	"github.com/yourusername/go-skeleton-code/internal/infrastructure/database"
	"github.com/yourusername/go-skeleton-code/internal/infrastructure/httpclient"
	"github.com/yourusername/go-skeleton-code/internal/infrastructure/logger"
	contactrepo "github.com/yourusername/go-skeleton-code/internal/repository/contact"
	contactusecase "github.com/yourusername/go-skeleton-code/internal/usecase/contact"
)

func main() {
	// ── 1. Configuration ──────────────────────────────────────────────────────
	cfg, err := config.Load(".env")
	if err != nil {
		// Cannot use structured logger yet; fall back to stdlib.
		panic("failed to load config: " + err.Error())
	}

	// ── 2. Logger ─────────────────────────────────────────────────────────────
	log, err := logger.New(cfg.LogLevel, cfg.IsProd())
	if err != nil {
		panic("failed to init logger: " + err.Error())
	}
	defer log.Sync() //nolint:errcheck

	log.Info("starting application",
		zap.String("env", cfg.AppEnv),
		zap.String("port", cfg.ServerPort),
	)

	// ── 3. Database ───────────────────────────────────────────────────────────
	db, err := database.NewPostgresDB(cfg, log)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}

	// ── 4. External HTTP Client ───────────────────────────────────────────────
	extClient := httpclient.New(cfg.ExternalAPIBaseURL, cfg.ExternalAPITimeout, log)

	// Optional: exercise the external client during startup (dev only).
	if !cfg.IsProd() {
		go func() {
			post, err := httpclient.FetchExamplePost(context.Background(), extClient)
			if err != nil {
				log.Warn("external API example call failed", zap.Error(err))
				return
			}
			log.Info("external API example response",
				zap.Int("post_id", post.ID),
				zap.String("title", post.Title),
			)
		}()
	}

	// ── 5. Dependency Injection (manual wire) ─────────────────────────────────
	// Repository layer
	contactRepo := contactrepo.New(db, log)

	// Usecase layer
	contactUC := contactusecase.New(contactRepo, log)

	// Delivery layer
	contactHandler := handler.NewContactHandler(contactUC, log)

	// ── 6. Router ─────────────────────────────────────────────────────────────
	engine := router.Setup(log, contactHandler)

	// ── 7. HTTP Server ────────────────────────────────────────────────────────
	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine so it doesn't block the shutdown signal listener.
	go func() {
		log.Info("http server listening", zap.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("http server error", zap.Error(err))
		}
	}()

	// ── 8. Graceful Shutdown ──────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server…")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown", zap.Error(err))
	}

	// Close database connection pool.
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}

	log.Info("server stopped gracefully")
}
