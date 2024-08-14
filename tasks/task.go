package tasks

import "context"

type Task interface {
	Id() string
	Exec(ctx context.Context) error
}
