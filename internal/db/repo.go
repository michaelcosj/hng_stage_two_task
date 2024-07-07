package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Db interface {
	DBTX
	Begin(ctx context.Context) (pgx.Tx, error)
}

type RepoQuerier interface {
	Querier
	WithTx(tx pgx.Tx) RepoQuerier
	GetDB() Db
}

type RepoQueries struct {
	*Queries
	db Db
}

func (r *RepoQueries) WithTx(tx pgx.Tx) RepoQuerier {
	return &RepoQueries{
		Queries: r.Queries.WithTx(tx),
		db:      r.db,
	}
}

func (r *RepoQueries) GetDB() Db {
	return r.db
}

func NewRepoQuerier(q *Queries, db Db) RepoQuerier {
	return &RepoQueries{
		Queries: q,
		db:      db,
	}
}
