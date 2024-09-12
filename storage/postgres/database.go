package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, cfg *Config) (*Repository, error) {
	connConfig, err := pgxpool.ParseConfig(ConnectionString(cfg))
	if err != nil {
		return nil, err
	}
	connConfig.MinConns = cfg.MaxIdleConns
	connConfig.MaxConns = cfg.MaxOpenConns

	conn, err := pgxpool.ConnectConfig(ctx, connConfig)
	if err != nil {
		return nil, err
	}

	return &Repository{db: conn}, nil
}

func (r *Repository) Close() {
	r.db.Close()
}

func (r *Repository) Start(ctx context.Context, options pgx.TxOptions) (pgx.Tx, error) {
	return r.db.BeginTx(ctx, options)
}

func (r *Repository) DB() *pgxpool.Pool {
	return r.db
}
func ApplyMigrate(databaseUrl, migrationsDir string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("could not find migration path")
	}
	dir := path.Join(path.Dir(filename), migrationsDir)

	mig, err := migrate.New(
		fmt.Sprintf("file://%s", dir),
		databaseUrl)
	if err != nil {
		slog.Error("failed to create migrations instance", slog.Any("err", err))
		return err
	}

	if err = mig.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("could not exec migration", slog.Any("err", err))
		return err
	}
	return nil
}
