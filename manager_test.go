package ticker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestManager_RegisterWorker(t *testing.T) {
	m := DefaultManager().(*manager)

	df, err := m.RegisterWorker("foo", NewGenericWorker("foo", time.Minute))
	assert.NoError(t, err)

	w, ok := m.workers.Load("foo")
	assert.True(t, ok)
	assert.NotNil(t, w)

	df()

	w, ok = m.workers.Load("foo")
	assert.False(t, ok)
	assert.Nil(t, w)
}

type mockJob struct {
}

func (m *mockJob) Tick(_ *Job, _ time.Time) {
}

func TestManager_RegisterJob(t *testing.T) {
	m := DefaultManager().(*manager)
	worker := NewGenericWorker("foo", time.Minute)
	assert.Equal(t, 0, worker.jobs.Count())

	_, werr := m.RegisterWorker("foo", worker)
	assert.NoError(t, werr)

	job, err := m.RegisterJob("foo", &mockJob{})
	assert.NoError(t, err)

	assert.Equal(t, 1, worker.jobs.Count())
	w, ok := m.jobs.Load(job.ID)
	assert.True(t, ok)
	assert.NotNil(t, w)

	job.Done()

	assert.Equal(t, 0, worker.jobs.Count())
	w, ok = m.jobs.Load(job.ID)
	assert.False(t, ok)
	assert.Nil(t, w)
}

func TestManager_RequeueJob(t *testing.T) {
	m := DefaultManager().(*manager)
	worker := NewGenericWorker("foo", time.Minute)

	_, werr := m.RegisterWorker("foo", worker)
	assert.NoError(t, werr)

	job, err := m.RegisterJob("foo", &mockJob{})
	assert.NoError(t, err)

	assert.Equal(t, 1, worker.jobs.Count())
	w, ok := m.jobs.Load(job.ID)
	assert.True(t, ok)
	assert.NotNil(t, w)

	worker.run(time.Now())

	assert.Equal(t, 0, worker.jobs.Count())
	w, ok = m.jobs.Load(job.ID)
	assert.True(t, ok)
	assert.NotNil(t, w)

	err = m.RequeueJob(job)
	assert.NoError(t, err)

	assert.Equal(t, 1, worker.jobs.Count())
	w, ok = m.jobs.Load(job.ID)
	assert.True(t, ok)
	assert.NotNil(t, w)
}

func TestManager_RequeueJob_AlreadyQueued(t *testing.T) {
	m := DefaultManager().(*manager)
	worker := NewGenericWorker("foo", time.Minute)

	_, werr := m.RegisterWorker("foo", worker)
	assert.NoError(t, werr)

	job, err := m.RegisterJob("foo", &mockJob{})
	assert.NoError(t, err)

	assert.Equal(t, 1, worker.jobs.Count())
	w, ok := m.jobs.Load(job.ID)
	assert.True(t, ok)
	assert.NotNil(t, w)

	err = job.Requeue()
	assert.Error(t, err)
}

func TestManager_RequeueJob_JobAlreadyMarkedAsDone(t *testing.T) {
	m := DefaultManager().(*manager)
	worker := NewGenericWorker("foo", time.Minute)

	_, werr := m.RegisterWorker("foo", worker)
	assert.NoError(t, werr)

	job, err := m.RegisterJob("foo", &mockJob{})
	assert.NoError(t, err)

	assert.Equal(t, 1, worker.jobs.Count())
	w, ok := m.jobs.Load(job.ID)
	assert.True(t, ok)
	assert.NotNil(t, w)

	job.Done()

	assert.Equal(t, 0, worker.jobs.Count())
	w, ok = m.jobs.Load(job.ID)
	assert.False(t, ok)
	assert.Nil(t, w)

	err = job.Requeue()
	assert.Error(t, err)
}

func TestManager_ListWorkers(t *testing.T) {
	m := DefaultManager().(*manager)
	worker := NewGenericWorker("foo", time.Minute)
	df, werr := m.RegisterWorker("foo", worker)
	assert.NoError(t, werr)
	assert.NotNil(t, df)

	names, _ := m.ListWorkers()
	assert.Contains(t, names, "foo")

	df()

	names, _ = m.ListWorkers()
	assert.NotContains(t, names, "foo")

}
