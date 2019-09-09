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
	err := scheduler.Add(job, time.Second)
	assert.Nil(t, err)

	go scheduler.Start()

	time.Sleep(2500 * time.Millisecond)
	scheduler.Stop()

	assert.Equal(t, 2, job.Calls())
}

func TestJobFunc(t *testing.T) {
	var counter int32 = 0
	fn := func() {
		atomic.AddInt32(&counter, 1)
	}

	scheduler := scheduler.Scheduler{}
	err := scheduler.AddFunc(fn, 100*time.Millisecond)
	assert.Nil(t, err)

	go scheduler.Start()

	time.Sleep(350 * time.Millisecond)
	scheduler.Stop()

	assert.Equal(t, int32(3), atomic.LoadInt32(&counter))
}

func TestAlreadyRunning(t *testing.T) {
	var counter1 int32 = 0
	var counter2 int32 = 0
	fn1 := func() {
		atomic.AddInt32(&counter1, 1)
	}
	fn2 := func() {
		atomic.AddInt32(&counter2, 1)
	}

	scheduler := scheduler.Scheduler{}
	err := scheduler.AddFunc(fn1, 100*time.Millisecond)
	assert.Nil(t, err)

	// start scheduler multiple times, only one instance of each job should be launched
	go scheduler.Start()
	go scheduler.Start()
	go scheduler.Start()

	time.Sleep(100 * time.Millisecond)

	err = scheduler.AddFunc(fn2, 100*time.Millisecond)
	assert.EqualError(t, err, "unable to add item to scheduler: scheduler is already running")

	time.Sleep(250 * time.Millisecond)

	// stop scheduler multiple times, only first Stop() will stop jobs
	scheduler.Stop()
	scheduler.Stop()
	scheduler.Stop()

	assert.Equal(t, int32(3), atomic.LoadInt32(&counter1))
	assert.Equal(t, int32(0), atomic.LoadInt32(&counter2))
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
		err := scheduler.Add(job.job, job.duration)
		assert.Nil(t, err)
	}

	go scheduler.Start()

	time.Sleep(totalDuration)
	scheduler.Stop()

	for _, job := range jobs {
		assert.Equal(t, job.expectedCalls, job.job.Calls())
	}
}
