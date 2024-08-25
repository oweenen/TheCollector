package tasks

import (
	"TheCollectorDG/db"
	"TheCollectorDG/riot"
	"context"
	"errors"
	"fmt"
	"log"
)

type AccountByPuuidTask struct {
	Puuid   string
	Cluster string
	Queries *db.Queries
}

func (task AccountByPuuidTask) Id() string {
	return fmt.Sprintf("AccountByPuuidTask-%v", task.Puuid)
}

func (task AccountByPuuidTask) Exec(ctx context.Context) error {
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

type AccountByNameTagTask struct {
	Name    string
	Tag     string
	Cluster string
	Queries *db.Queries
}

func (task AccountByNameTagTask) Id() string {
	return fmt.Sprintf("AccountByNameTagTask-%v#%v", task.Name, task.Tag)
}

func (task AccountByNameTagTask) Exec(ctx context.Context) error {
	res, err := riot.GetAccountByName(task.Cluster, task.Name, task.Tag)
	if err != nil {
		return err
	}

	err = task.Queries.InsertAccount(ctx, db.InsertAccountParams{
		Puuid: res.Puuid,
		Name:  res.Name,
		Tag:   res.Tag,
	})

	log.Printf("Account collected with name %v#%v\n", task.Name, task.Tag)

	return err
}
