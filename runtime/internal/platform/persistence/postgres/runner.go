package postgres

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/iceiko/vibit/runtime/internal/app/session"
	"github.com/iceiko/vibit/runtime/internal/modules/authentication"
	"github.com/iceiko/vibit/runtime/internal/modules/friends"
	"github.com/iceiko/vibit/runtime/internal/modules/inventory"
	"github.com/iceiko/vibit/runtime/internal/modules/player"
	"github.com/iceiko/vibit/runtime/internal/modules/storage"
	"github.com/iceiko/vibit/runtime/internal/platform/tx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Beginner interface {
	Begin(context.Context) (pgx.Tx, error)
}

type Runner struct {
	Beginner Beginner
}

func NewRunner(beginner Beginner) *Runner {
	return &Runner{Beginner: beginner}
}

func NewPoolRunner(pool *pgxpool.Pool) *Runner {
	if pool == nil {
		return NewRunner(nil)
	}
	return NewRunner(pool)
}

func (r *Runner) WithinUnitOfWork(ctx context.Context, fn func(context.Context, tx.UnitOfWork) error) error {
	if fn == nil {
		return errors.New("postgres tx: unit-of-work function is required")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if r == nil || isNilBeginner(r.Beginner) {
		return errors.New("postgres tx: beginner is required")
	}

	pgxTx, err := r.Beginner.Begin(ctx)
	if err != nil {
		return fmt.Errorf("postgres tx: begin: %w", err)
	}

	unit := UnitOfWork{
		ctx:      ctx,
		executor: pgxTx,
	}
	if err := fn(ctx, unit); err != nil {
		if rollbackErr := pgxTx.Rollback(ctx); rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			return fmt.Errorf("postgres tx: rollback after function error: %v: %w", rollbackErr, err)
		}
		return err
	}

	if err := pgxTx.Commit(ctx); err != nil {
		if rollbackErr := pgxTx.Rollback(ctx); rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			return fmt.Errorf("postgres tx: rollback after commit error: %v: %w", rollbackErr, err)
		}
		return fmt.Errorf("postgres tx: commit: %w", err)
	}
	return nil
}

func isNilBeginner(beginner Beginner) bool {
	if beginner == nil {
		return true
	}
	value := reflect.ValueOf(beginner)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

type UnitOfWork struct {
	ctx      context.Context
	executor Executor
}

func (u UnitOfWork) Context() context.Context {
	if u.ctx == nil {
		return context.Background()
	}
	return u.ctx
}

func (u UnitOfWork) Executor() (Executor, error) {
	if u.executor == nil {
		return nil, errors.New("postgres tx: executor is required")
	}
	return u.executor, nil
}

func (u UnitOfWork) NewInventoryRepository() (inventory.Repository, error) {
	executor, err := u.Executor()
	if err != nil {
		return nil, err
	}
	return NewInventoryRepositoryForUnitOfWork(executor), nil
}

func (u UnitOfWork) NewPlayerAccountRepository() (player.Repository, error) {
	executor, err := u.Executor()
	if err != nil {
		return nil, err
	}
	return NewPlayerAccountRepositoryForUnitOfWork(executor), nil
}

func (u UnitOfWork) NewAuthenticationRepository() (authentication.Repository, error) {
	executor, err := u.Executor()
	if err != nil {
		return nil, err
	}
	return NewAuthenticationRepositoryForUnitOfWork(executor), nil
}

func (u UnitOfWork) NewSessionRepository() (session.Repository, error) {
	executor, err := u.Executor()
	if err != nil {
		return nil, err
	}
	return NewSessionRepositoryForUnitOfWork(executor), nil
}

func (u UnitOfWork) NewStorageObjectRepository() (storage.Repository, error) {
	executor, err := u.Executor()
	if err != nil {
		return nil, err
	}
	return NewStorageObjectRepositoryForUnitOfWork(executor), nil
}

func (u UnitOfWork) NewFriendRelationshipRepository() (friends.Repository, error) {
	executor, err := u.Executor()
	if err != nil {
		return nil, err
	}
	return NewFriendRelationshipRepositoryForUnitOfWork(executor), nil
}
