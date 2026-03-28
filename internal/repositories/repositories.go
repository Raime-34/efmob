package repositories

import "github.com/jackc/pgx/v5/pgxpool"

type Repo struct{}

func NewRepo(conn *pgxpool.Pool) *Repo {
	return &Repo{}
}
