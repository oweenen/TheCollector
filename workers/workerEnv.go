package workers

import (
	"TheCollectorDG/db"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const BACKOFF_TIME = time.Second * 1

type WorkerEnv struct {
	Pool    *pgxpool.Pool
	Queries *db.Queries
}
