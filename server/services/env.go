package services

import (
	"TheCollectorDG/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceEnv struct {
	Pool    *pgxpool.Pool
	Queries *db.Queries
}
