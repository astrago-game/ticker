package ticker

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type DeleteFn func()

type Manager interface {
	Start()
	Stop()
	RegisterWorker(string, Worker) (DeleteFn, error)
	RegisterJob(string, JobObject) (*Job, error)
	RequeueJob(*Job) error
	ListWorkers() ([]string, error)
	ListJobs() ([]string, error)
	GetJob(string) (*Job, bool)
}

// DefaultManager implements a map backed Manager
func DefaultManager() Manager {
	return &manager{
		wg:      &sync.WaitGroup{},
		workers: workerMap{},
		jobs:    jobMap{},
	}
}

// SetManager allows you to replace the default
// ticker manager with a custom one
func SetManager(m Manager) {
	ticker = m
}

var ticker Manager = DefaultManager()

type manager struct {
	wg      *sync.WaitGroup
	workers workerMap
	jobs    jobMap
}

func (m *manager) Start() {
	m.workers.Range(func(key string, w Worker) bool {
		w.Start(m.wg)
		return true
	})
}

func (m *manager) Stop() {
	m.workers.Range(func(key string, w Worker) bool {
		w.Stop(m.wg)
		return true
	})
	m.wg.Wait()
}

func (m *manager) RegisterWorker(name string, worker Worker) (DeleteFn, error) {
	if _, found := m.workers.Load(name); found {
		return nil, fmt.Errorf("worker named %s is already loaded", name)
	}

	m.workers.Store(name, worker)

	df := func() {
		m.workers.Delete(name)
	}

	return df, nil
}

func (m *manager) RegisterJob(workerName string, object JobObject) (*Job, error) {
	worker, found := m.workers.Load(workerName)
	if !found {
		return nil, fmt.Errorf("worker named %s not registered", workerName)
	}

	job := Job{
		ID:         uuid.New().String(),
		ticker:     m,
		workerName: workerName,
		Object:     object,
	}
	job.df = func() {
		worker.delete(job.ID)
		m.jobs.Delete(job.ID)
	}

	m.jobs.LoadOrStore(job.ID, &job)
	if err := worker.Enqueue(&job); err != nil {
		return nil, err
	}

	return &job, nil
}

func (m *manager) RequeueJob(job *Job) error {
	if job.done {
		return fmt.Errorf("job id %s is already marked as done and can't be requeued", job.ID)
	}

	worker, found := m.workers.Load(job.workerName)
	if !found {
		return fmt.Errorf("worker named %s not registered", job.workerName)
	}

	return worker.Enqueue(job)
}

func (m *manager) ListWorkers() ([]string, error) {
	return m.workers.Keys(), nil
}

func (m *manager) ListJobs() ([]string, error) {
	return m.jobs.Keys(), nil
}

func (m *manager) GetJob(id string) (*Job, bool) {
	return m.jobs.Load(id)
}
