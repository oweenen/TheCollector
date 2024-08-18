package tasks

import (
	"TheCollectorDG/db"
	"TheCollectorDG/riot"
	"context"
	"errors"
	"fmt"
	"log"
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
	if errors.Is(err, riot.NotFoundError) {
		err = task.Queries.SetSkipAccountFlag(ctx, db.SetSkipAccountFlagParams{
			Puuid:       task.Puuid,
			SkipAccount: true,
		})
		return err
	}
	if err != nil {
		return err
	}

	err = task.Queries.UpdateAccount(ctx, db.UpdateAccountParams{
		Puuid: task.Puuid,
		Name:  res.Name,
		Tag:   res.Tag,
	})

	log.Printf("Account details collected for %v\n", task.Puuid)

	return err
}
