package main

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

	"github.com/MikebangSfilya/tree-tracker/internal/completion"
	"github.com/MikebangSfilya/tree-tracker/internal/config"
	"github.com/MikebangSfilya/tree-tracker/internal/plant"
	"github.com/MikebangSfilya/tree-tracker/internal/progress"
	"github.com/MikebangSfilya/tree-tracker/internal/routine"
	"github.com/MikebangSfilya/tree-tracker/internal/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}
}

func run() error {
	logger := initLogger()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	configPath := os.Getenv("CONFIG_PATH")

	config, err := config.LoadConfig(configPath)
	if err != nil {
		logger.Error("failed to load config",
			"config path", configPath,
			"error", err,
		)
		return err
	}

	server, err := reg(ctx, logger, config)
	if err != nil {
		logger.Error("failed to register server components",
			"error", err,
		)
		return err
	}

	errCh := make(chan error, 1)

	go func() {
		logger.Info("Listening", "address", server.Addr)
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		logger.Info("received shutdown signal")
		return server.Shutdown(shutdownCtx)
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("error while listening", "error", err)
			return err
		}
		return nil
	}
}

func reg(ctx context.Context, logger *slog.Logger, config *config.Config) (*http.Server, error) {
	pool, err := connectDatabase(ctx, config.DataBaseURL)
	if err != nil {
		logger.Error("failed to connect to database",
			"database url", config.DataBaseURL,
			"error", err,
		)
		return nil, err
	}

	completionRepo := completion.NewCompletionRepository1(pool, logger)
	plantRepo := plant.NewRepository(pool)
	progressRepo := progress.NewRepository(pool)
	routineRepo := routine.NewPostgresRepository(pool)
	userRepo := user.NewRepoDB(pool)

	progressServ := progress.NewService(progressRepo, completionRepo)
	plantServ := plant.NewService(plantRepo)
	routineServ := routine.NewService(routineRepo)
	userServ := user.NewService(logger, userRepo)

	completionHandler := completion.NewCompletionHandler(progressServ)
	plantHandler := plant.NewHandler(plantServ, logger)
	routineHandler := routine.NewHandler(routineServ)
	userHandler := user.NewHandler(logger, userServ)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	completionHandler.Routes(r)
	plantHandler.Routes(r)
	routineHandler.RegisterRoutes(r)
	userHandler.Routes(r)

	addr := fmt.Sprintf("%s:%d", config.APIHost, config.APIPort)
	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return server, nil
}

func connectDatabase(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

func initLogger() *slog.Logger {
	var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	slog.SetDefault(logger)
	return logger
}
