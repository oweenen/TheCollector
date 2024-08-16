package workers

import (
	"TheCollectorDG/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkerEnv struct {
	Pool    *pgxpool.Pool
	Queries *db.Queries
}
