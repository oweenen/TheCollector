package tasks

import (
	"TheCollectorDG/db"
	"TheCollectorDG/riot"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgtype"
)

type AccountDetailsTask struct {
	Puuid   string
	Cluster string
	Queries *db.Queries
}

func (task AccountDetailsTask) Id() string {
	return fmt.Sprintf("AccountDetailsTask-%v", task.Puuid)
}

func (task AccountDetailsTask) Exec(ctx context.Context) error {
	res, err := riot.GetAccountByPuuid(task.Cluster, task.Puuid)
	if err != nil {
		return err
	}

	err = task.Queries.UpdateAccount(ctx, db.UpdateAccountParams{
		Puuid: task.Puuid,
		Name: pgtype.Text{
			String: res.Name,
			Valid:  true,
		},
		Tag: pgtype.Text{
			String: res.Tag,
			Valid:  true,
		},
	})

	log.Printf("Account details collected for %v\n", task.Puuid)

	return err
}
