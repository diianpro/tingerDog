package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path"
	"runtime"

	"github.com/diianpro/tingerDog/internal/storage/postgres/models"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *pgxpool.Pool
}

type Repo interface {
	GetUsers(ctx context.Context) ([]models.User, error)

	Do(ctx context.Context, fn func(c context.Context) error) error
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

func (r *Repository) database(ctx context.Context) Tx {
	return DefaultTrOrDB(ctx, r.DB())
}

func (r *Repository) txFactory(ctx context.Context, options pgx.TxOptions) (pgx.Tx, error) {
	return r.Start(ctx, options)
}

func (r *Repository) Do(ctx context.Context, fn func(c context.Context) error) error {
	opts := pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	}
	tx, err := r.txFactory(ctx, opts)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = errors.Join(tx.Rollback(ctx))
		} else {
			err = errors.Join(tx.Commit(ctx))
		}
	}()

	c := context.WithValue(ctx, defaultCtxKey, tx)
	err = fn(c)
	if err != nil {
		return err
	}
	return nil
}

type ctxKey struct{}

var defaultCtxKey = ctxKey{}

func DefaultTrOrDB(ctx context.Context, db *pgxpool.Pool) Tx {
	if tr, ok := ctx.Value(defaultCtxKey).(pgx.Tx); ok {
		return tr
	}
	return db
}

type Tx interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
}
