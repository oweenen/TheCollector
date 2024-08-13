package jobqueue

type Job interface {
	id() string
	execute()
}

type JobQueue struct {
	queue  chan Job
	jobMap map[string]*Job
}

func (jq *JobQueue) Push(j Job) {
	if _, has := jq.jobMap[j.id()]; !has {
		jq.jobMap[j.id()] = &j
		jq.queue <- j
	}
}

func (jq *JobQueue) Pop() Job {
	j := <-jq.queue
	delete(jq.jobMap, j.id())
	return <-jq.queue
}
