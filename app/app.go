package app

import (
	"context"
	"log"
	"log/slog"

	"github.com/diianpro/tingerDog/internal/config"
	"github.com/diianpro/tingerDog/internal/service"
	"github.com/diianpro/tingerDog/internal/storage/postgres"
	"github.com/diianpro/tingerDog/internal/transport"
	"github.com/diianpro/tingerDog/internal/transport/handler"
)

type LoadConfigFn func() (*config.Config, error)

type App struct {
	cfg      *config.Config
	srv      *transport.Server
	ctx      context.Context
	cancelFn context.CancelFunc
}

func New(loadConfigFn LoadConfigFn) *App {
	ctx, cancelFn := context.WithCancel(context.Background())
	cfg, err := loadConfigFn()
	if err != nil {
		log.Fatal("failed to load config")
	}

	return &App{
		cfg:      cfg,
		ctx:      ctx,
		cancelFn: cancelFn,
	}
}

func (a *App) Start() {
	defer a.cancelFn()
	ctx := context.Background()
	db, err := postgres.New(a.ctx, &a.cfg.Postgres)
	if err != nil {
		slog.Error("Could not setup storage", err.Error())
	}

	if err = postgres.ApplyMigrate(postgres.ConnectionString(&a.cfg.Postgres), "../../../migration"); err != nil {
		slog.Error("Could not apply migrations.")
	}
	defer db.Close()

	repo, err := postgres.New(ctx, &a.cfg.Postgres)
	if err != nil {
		slog.Error("Could not set up repository.")
	}
	src := service.New(repo)
	hndl := handler.New(src)

	a.srv = transport.New(hndl)
}

func (a *App) Stop() {
	a.cancelFn()
}
