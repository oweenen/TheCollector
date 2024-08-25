package workerManager

import "context"

type Task interface {
	Id() string
	Exec(ctx context.Context) error
}

type TaskWithDone struct {
	Task Task
	Done chan error
}

func (twd TaskWithDone) Id() string {
	return twd.Id()
}

func (twd TaskWithDone) Exec(ctx context.Context) error {
	err := twd.Task.Exec(ctx)
	twd.Done <- err
	close(twd.Done)
	return err
}
