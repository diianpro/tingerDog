package repository

import (
	"context"
	"errors"

	"github.com/diianpro/tingerDog/repository/models"
	"github.com/diianpro/tingerDog/storage/postgres"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	driver *postgres.Repository
}

type Repo interface {
	GetUsers(ctx context.Context) ([]models.User, error)

	Do(ctx context.Context, fn func(c context.Context) error) error
}

func NewRepository(db *postgres.Repository) *Repository {
	return &Repository{
		driver: db,
	}
}

func (r *Repository) db(ctx context.Context) Tx {
	return DefaultTrOrDB(ctx, r.driver.DB())
}

func (r *Repository) txFactory(ctx context.Context, options pgx.TxOptions) (pgx.Tx, error) {
	return r.driver.Start(ctx, options)
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
