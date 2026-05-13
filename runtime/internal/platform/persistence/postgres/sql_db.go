package postgres

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func OpenSQLDBFromPool(pool *pgxpool.Pool) (*sql.DB, error) {
	if pool == nil {
		return nil, errors.New("postgres config: pool is required")
	}
	return stdlib.OpenDBFromPool(pool), nil
}
