package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/437d5/pr-review-manager/internal/application/http/handlers"
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
		return db.NewUnitOfWork(conn), nil
	}

	teamService := services.NewTeamService(uowFactory)
	teamHandler := handlers.NewTeamHandler(teamService)

	userService := services.NewUserService(uowFactory)
	userHandler := handlers.NewUserHandler(userService)

	prService := services.NewPRService(uowFactory)
	prHandler := handlers.NewPRHandler(prService)

	http.HandleFunc("POST /team/add", teamHandler.CreateTeam)
	http.HandleFunc("GET /team/get", teamHandler.GetTeam)

	http.HandleFunc("POST /users/setIsActive", userHandler.SetIsActive)
	http.HandleFunc("GET /users/getReview", userHandler.GetPRs)

	http.HandleFunc("POST /pullRequest/create", prHandler.CreatePR)
	http.HandleFunc("POST /pullRequest/merge", prHandler.Merge)
	http.HandleFunc("POST /pullRequest/reassign", prHandler.Reassign)

	server := &http.Server{
		Addr: cfg.Address,
	}

	log.Fatal(server.ListenAndServe())
}
