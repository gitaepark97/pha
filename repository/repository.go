package repository

import (
	"database/sql"
)

type Repository interface {
	Querier
}

type repository struct {
	db *sql.DB
	*Queries
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db:      db,
		Queries: New(db),
	}
}
