package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/qerdcv/qerdcv/internal/config"
	"github.com/qerdcv/qerdcv/internal/repositories"
	"github.com/qerdcv/qerdcv/internal/repositories/migrations"
	"github.com/qerdcv/qerdcv/internal/server"
	"github.com/qerdcv/qerdcv/internal/server/handlers"
	"github.com/qerdcv/qerdcv/internal/server/middlewares"
	"github.com/qerdcv/qerdcv/internal/services"
	"github.com/qerdcv/qerdcv/pkg/migrator"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg, err := config.New()
	if err != nil {
		logger.Error("failed to build config", slog.Any("err", err))
		return
	}

	v, d, err := migrator.Migrate(migrations.Migrations, cfg.DB.DSN())
	if err != nil {
		logger.Error("apply migration",
			slog.Any("err", err),
			slog.Int("version", int(v)),
			slog.Bool("dirty", d))
		return
	}

	logger.Info("apply migration",
		slog.Any("err", err),
		slog.Int("version", int(v)),
		slog.Bool("dirty", d))

	appCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
	defer cancel()

	db, err := sql.Open("postgres", cfg.DB.DSN())
	if err != nil {
		logger.Error("sql open", slog.Any("err", err))
		return
	}

	userService := services.NewUserService(
		repositories.NewUserRepo(db),
	)

	s := server.New(
		logger,
		cfg.Server,
		middlewares.Auth(userService),
		handlers.NewUserHandler(
			logger,
			userService,
		),
		handlers.NewBudgetHandler(
			logger,
			services.NewBudgetService(
				repositories.NewBudgetRepo(db),
			),
		),
	)

	if err = s.Run(appCtx); err != nil {
		logger.Error("run server", slog.Any("err", err))
		return
	}
}
