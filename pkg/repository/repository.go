package repository

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB      *pgxpool.Pool
	Queries *db.Queries
}

func NewRepository(dbConnection *pgxpool.Pool) *Repository {
	return &Repository{
		DB:      dbConnection,
		Queries: db.New(dbConnection),
	}
}
