package scheduler

import (
	"time"
)

// Job interface with Run method
type Job interface {
	Run()
}

type item struct {
	job      Job
	duration time.Duration
	done     chan struct{}
}

// Scheduler to run jobs repeatedly
type Scheduler struct {
	items []item
}

// Add adds job to scheduler
func (scheduler *Scheduler) Add(job Job, duration time.Duration) {
	scheduler.items = append(scheduler.items, item{
		job:      job,
		duration: duration,
		done:     make(chan struct{}),
	})
}

// Run runs scheduler
func (scheduler *Scheduler) Run(done <-chan struct{}) {
	for i := range scheduler.items {
		scheduler.runItem(i)
	}

	<-done

	for _, item := range scheduler.items {
		item.done <- struct{}{}
	}
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
