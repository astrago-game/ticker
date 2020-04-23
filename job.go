package ticker

import "time"

type JobObject interface {
	Tick(job *Job, now time.Time)
}

type Job struct {
	ID         string
	ticker     Manager
	workerName string
	Object     JobObject
	done       bool
	df         DeleteFn
}

func (j *Job) Requeue() error {
	return j.ticker.RequeueJob(j)
}

func (j *Job) MoveTo(workerName string) (*Job, error) {
	j.workerName = workerName
	j.df()
	return j.ticker.RegisterJob(workerName, j.Object)
}

func (j *Job) Done() {
	j.done = true
	j.df()
}
