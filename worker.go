package ticker

import (
	"fmt"
	"sync"
	"time"
)

type Worker interface {
	GetName() string
	Start(*sync.WaitGroup)
	Stop(*sync.WaitGroup)
	Enqueue(*Job) error
	delete(string)
}

var _ Worker = &genericWorker{}

type genericWorker struct {
	Name   string
	jobs   jobMap
	ticker *time.Ticker
	lock   sync.Mutex
}

func (w *genericWorker) delete(jobId string) {
	w.jobs.Delete(jobId)
}

func (w *genericWorker) GetName() string {
	return w.Name
}

func (w *genericWorker) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		for next := range w.ticker.C {
			w.run(next)
		}
		wg.Done()
	}(wg)
}

func (w *genericWorker) Stop(wg *sync.WaitGroup) {
	w.ticker.Stop()
	wg.Done()
}

func (w *genericWorker) Enqueue(job *Job) error {
	if _, found := w.jobs.Load(job.ID); found {
		return fmt.Errorf("job of id %s is already queued", job.ID)
	}
	w.lock.Lock()
	w.jobs.Store(job.ID, job)
	w.lock.Unlock()
	return nil
}

func (w *genericWorker) run(now time.Time) {
	if w.jobs.Count() == 0 {
		return
	}

	w.lock.Lock()
	var keys []string
	w.jobs.Range(func(key string, value *Job) bool {
		go func(j *Job) {
			j.Object.Tick(j, now)
		}(value)
		keys = append(keys, key)
		return true
	})
	for _, key := range keys {
		w.jobs.Delete(key)
	}
	w.lock.Unlock()

}

func NewGenericWorker(name string, d time.Duration) *genericWorker {
	return &genericWorker{
		Name:   name,
		jobs:   jobMap{},
		ticker: time.NewTicker(d),
	}
}
