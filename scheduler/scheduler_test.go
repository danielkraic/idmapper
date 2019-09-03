package scheduler_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/danielkraic/idmapper/scheduler"
	"github.com/stretchr/testify/assert"
)

type Job struct {
	calls int32
}

func (job *Job) Run() {
	atomic.AddInt32(&job.calls, 1)
}

func (job *Job) Calls() int {
	return int(atomic.LoadInt32(&job.calls))
}

func TestJob(t *testing.T) {
	job := &Job{}

	scheduler := scheduler.Scheduler{}
	scheduler.Add(job, time.Second)

	done := make(chan struct{})
	go scheduler.Run(done)

	time.Sleep(2500 * time.Millisecond)
	done <- struct{}{}

	assert.Equal(t, 2, job.Calls())
}

func TestMultipleJobs(t *testing.T) {
	jobs := []struct {
		job           *Job
		duration      time.Duration
		expectedCalls int
	}{
		{&Job{}, 100 * time.Millisecond, 19},
		{&Job{}, 200 * time.Millisecond, 9},
		{&Job{}, 500 * time.Millisecond, 3},
		{&Job{}, 1000 * time.Millisecond, 1},
		{&Job{}, 2000 * time.Millisecond, 0},
		{&Job{}, 5000 * time.Millisecond, 0},
	}

	totalDuration := 1950 * time.Millisecond

	scheduler := scheduler.Scheduler{}
	for _, job := range jobs {
		scheduler.Add(job.job, job.duration)
	}

	done := make(chan struct{})
	go scheduler.Run(done)

	time.Sleep(totalDuration)
	done <- struct{}{}

	for _, job := range jobs {
		assert.Equal(t, job.expectedCalls, job.job.Calls())
	}
}
