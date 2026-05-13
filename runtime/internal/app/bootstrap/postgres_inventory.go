package bootstrap

import (
	"context"
	"errors"
	"fmt"

	"github.com/iceiko/vibit/runtime/internal/modules/inventory"
	"github.com/iceiko/vibit/runtime/internal/platform/tx"
)

type UnitOfWorkInventoryRepositoryFactory interface {
	NewInventoryRepository() (inventory.Repository, error)
}

type PostgresInventoryRepositoryProvider struct {
	QueryRepository inventory.Repository
}

func (p PostgresInventoryRepositoryProvider) ForCommand(_ context.Context, unit tx.UnitOfWork) (inventory.Repository, error) {
	factory, ok := unit.(UnitOfWorkInventoryRepositoryFactory)
	if !ok {
		return nil, fmt.Errorf("inventory bootstrap: unit of work %T cannot create inventory repository", unit)
	}

	repository, err := factory.NewInventoryRepository()
	if err != nil {
		return nil, err
	}
	if repository == nil {
		return nil, errors.New("inventory bootstrap: unit of work returned nil inventory repository")
	}
	return repository, nil
}

func (p PostgresInventoryRepositoryProvider) ForQuery(context.Context) (inventory.Repository, error) {
	if p.QueryRepository == nil {
		return nil, errors.New("inventory bootstrap: query repository is required")
	}
	return p.QueryRepository, nil
}
