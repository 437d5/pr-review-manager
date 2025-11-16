package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/437d5/pr-review-manager/internal/application/http/handlers"
	"github.com/437d5/pr-review-manager/internal/application/routers"
	"github.com/437d5/pr-review-manager/internal/domain/repositories"
	"github.com/437d5/pr-review-manager/internal/domain/services"
	"github.com/437d5/pr-review-manager/internal/infrastructure/db"
	"github.com/437d5/pr-review-manager/pkg/config"
	"github.com/437d5/pr-review-manager/pkg/logger"
	"github.com/jmoiron/sqlx"
)

func main() {
	cfg := config.MustLoadConfig()
	logger.InitLogger(cfg.Mode)

	slog.Info("starting server")
	slog.Debug("debug messages are enabled")

	conn, err := sqlx.Connect("postgres", cfg.GetConnectionString())
	if err != nil {
		slog.Error("db connection failed", "error", err.Error())
		os.Exit(1)
	}

	migrator := db.NewMigrator(conn)
	if err := migrator.Migrate(); err != nil {
		slog.Error("migration failed",
			"error", err.Error(),
		)
		os.Exit(1)
	}

	// for each request we have instance of UnitOfWork
	uowFactory := func(ctx context.Context) (repositories.UnitOfWork, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return db.NewUnitOfWork(conn), nil
		}
	}

	teamService := services.NewTeamService(uowFactory)
	teamHandler := handlers.NewTeamHandler(teamService)

	userService := services.NewUserService(uowFactory)
	userHandler := handlers.NewUserHandler(userService)

	prService := services.NewPRService(uowFactory)
	prHandler := handlers.NewPRHandler(prService)

	router := routers.InitRouter(
		teamHandler,
		userHandler,
		prHandler,
	)

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		slog.Debug("shutting down server")
		if err := server.Shutdown(context.Background()); err != nil {
			slog.Error("server shutdown failed",
				"error", err,
			)
		}
	}()

	slog.Info("starting server", "address", cfg.Address)
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server closed unexpectedly", "error", err)
			return
		}
	}
}
