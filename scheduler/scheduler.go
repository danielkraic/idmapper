package scheduler

import (
	"fmt"
	"sync/atomic"
	"time"
)

// Job interface with Run method
type Job interface {
	Run()
}

// JobFunc is adapter to allow use ordinary function as Job
type JobFunc func()

// Run runs job
func (jobFunc JobFunc) Run() {
	jobFunc()
}

type item struct {
	job Job
	// duration between last call and next call of job
	duration time.Duration
	done     chan struct{}
}

// Scheduler to run periodic jobs
type Scheduler struct {
	items     []item
	isRunning int32
	done      chan struct{}
}

// IsRunning return true it scheduler is already running
func (scheduler *Scheduler) IsRunning() bool {
	return atomic.LoadInt32(&scheduler.isRunning) == 1
}

func (scheduler *Scheduler) setIsRunning(isRunning bool) {
	if isRunning {
		atomic.StoreInt32(&scheduler.isRunning, 1)
	} else {
		atomic.StoreInt32(&scheduler.isRunning, 0)
	}
}

// Add adds job to scheduler
func (scheduler *Scheduler) Add(job Job, duration time.Duration) error {
	if scheduler.IsRunning() {
		return fmt.Errorf("unable to add item to scheduler: scheduler is already running")
	}

	scheduler.items = append(scheduler.items, item{
		job:      job,
		duration: duration,
		done:     make(chan struct{}),
	})

	return nil
}

// AddFunc adds function to scheduler
func (scheduler *Scheduler) AddFunc(fn func(), duration time.Duration) error {
	return scheduler.Add(JobFunc(fn), duration)
}

// Start starts scheduler by running all its jobs repeatedly
func (scheduler *Scheduler) Start() {
	if scheduler.IsRunning() {
		return
	}

	scheduler.setIsRunning(true)

	for i := range scheduler.items {
		scheduler.runItem(i)
	}
}

// Stop stops scheduler
func (scheduler *Scheduler) Stop() {
	if !scheduler.IsRunning() {
		return
	}

	for _, item := range scheduler.items {
		item.done <- struct{}{}
	}

	defer scheduler.setIsRunning(false)
}

func (scheduler *Scheduler) runItem(index int) {
	ticker := time.NewTicker(scheduler.items[index].duration)

	go func() {
		for {
			select {
			case <-ticker.C:
				scheduler.items[index].job.Run()
				continue
			case <-scheduler.items[index].done:
				ticker.Stop()
				return
			}
		}
	}()
}
