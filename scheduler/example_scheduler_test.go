package scheduler_test

import (
	"fmt"
	"time"

	"github.com/danielkraic/idmapper/scheduler"
)

type FirstJob struct {
}

func (firstJob FirstJob) Run() {
	fmt.Printf("FirstJob called\n")
}

func SecondJobFunc() {
	fmt.Printf("SecondJob called\n")
}

func Example() {
	scheduler := scheduler.Scheduler{}
	//scheduler.Add(&FirstJob{}, 100*time.Millisecond)
	err := scheduler.AddFunc(SecondJobFunc, 300*time.Millisecond)
	if err != nil {
		panic(err)
	}

	go scheduler.Start()

	time.Sleep(time.Second)
	scheduler.Stop()

	// Output:
	// SecondJob called
	// SecondJob called
	// SecondJob called
}
