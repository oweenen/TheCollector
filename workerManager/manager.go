package workerManager

import (
	"fmt"
)

const QUEUE_BUFFER_LENGTH = 1000

type Manager struct {
	queueMap map[string]chan Task
}

func New() *Manager {
	return &Manager{
		make(map[string]chan Task),
	}
}

func (m *Manager) AddWorker(workerId string, workerFn func(chan Task)) error {
	if _, ok := m.queueMap[workerId]; ok {
		return fmt.Errorf("failed to add worker: worker with id '%v' already exists", workerId)
	}

	q := make(chan Task, QUEUE_BUFFER_LENGTH)
	m.queueMap[workerId] = q

	go workerFn(q)

	return nil
}

func (m *Manager) AssignTask(workerId string, task Task) error {
	q, ok := m.queueMap[workerId]
	if !ok {
		return fmt.Errorf("failed to assign task: no worker with id '%v'", workerId)
	}

	q <- task

	return nil
}
