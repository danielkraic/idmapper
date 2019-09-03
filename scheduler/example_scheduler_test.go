package scheduler_test

import (
	"fmt"
	"time"

	"github.com/danielkraic/idmapper/scheduler"
)

func Example() {
	job1 := &Job{}
	job2 := &Job{}

	scheduler := scheduler.Scheduler{}
	scheduler.Add(job1, 600*time.Millisecond)
	scheduler.Add(job2, 300*time.Millisecond)

	done := make(chan struct{})
	go scheduler.Run(done)

	time.Sleep(time.Second)
	done <- struct{}{}

	fmt.Printf("job1 calls: %d\n", job1.Calls())
	fmt.Printf("job2 calls: %d\n", job2.Calls())

	// Output:
	// job1 calls: 1
	// job2 calls: 3
}
