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
	closeChannel chan struct{}
}

type PostgresPool interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	GracefulStop(ctx context.Context) error
}

func (p *postgresPool) GracefulStop(ctx context.Context) error {
	op := "postgresPool.Close"

	go func() {
		defer close(p.closeChannel)
		p.Pool.Close()
	}()

	select {
	case <-p.closeChannel:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("%s: %w: %v", op, ErrCloseTimeout, ctx.Err())
	}
}

func New(ctx context.Context, connection string) (PostgresPool, error) {
	const op = "postgresPool.New"
	syncOnce.Do(func() {
		dbPool, err := pgxpool.New(ctx, connection)

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
		return nil, fmt.Errorf("%s: pool is empty after connection: %w", op, ErrBadConnection)
	}

	return &postgresPool{
		Pool:         pool,
		closeChannel: make(chan struct{}),
	}, nil
}
