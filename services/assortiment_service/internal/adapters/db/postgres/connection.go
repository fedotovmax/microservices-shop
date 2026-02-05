package postgres

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

var syncOnce sync.Once

var initErr error

var ErrBadConnection = errors.New("cannot create new db connection")

var ErrBadPing = errors.New("error when ping connection")

var ErrCloseTimeout = errors.New("the time to safely terminate the connection to the postgres pool has expired")

type postgresPool struct {
	*pgxpool.Pool
}

type PostgresPool interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	Stop(ctx context.Context) error
}

func New(ctx context.Context, newCfg *Config) (PostgresPool, error) {
	const op = "adapters.db.postgres.New"
	syncOnce.Do(func() {

		poolcfg, err := pgxpool.ParseConfig(newCfg.DSN)

		if err != nil {
			initErr = fmt.Errorf("%s: %w: %v", op, ErrBadConnection, err)
			return
		}

		//TODO: maybe add some logic to poolcfg?

		dbPool, err := pgxpool.NewWithConfig(ctx, poolcfg)

		if err != nil {
			initErr = fmt.Errorf("%s: %w: %v", op, ErrBadConnection, err)
			return
		}

		err = dbPool.Ping(ctx)

		if err != nil {
			initErr = fmt.Errorf("%s: %w: %v", op, ErrBadPing, err)
			return
		}

		pool = dbPool
	})

	if initErr != nil {
		return nil, initErr
	}

	if pool == nil {
		return nil, fmt.Errorf("%s: pool is empty: %w", op, ErrBadConnection)
	}

	return &postgresPool{
		Pool: pool,
	}, nil
}

func (p *postgresPool) Stop(ctx context.Context) error {
	op := "adapters.db.postgres.Stop"

	done := make(chan struct{})

	go func() {
		defer close(done)
		p.Pool.Close()
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("%s: %w: %v", op, ErrCloseTimeout, ctx.Err())
	}
}
