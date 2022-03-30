package main

import (
	"time"

	"github.com/Xunzhuo/async"
)

func main() {
	async.Q().
		SetMaxWaitQueueLength(100).
		SetMaxWorkQueueLength(100).
		Start().AddJobAndRun(async.NewJob(longTimeJob))
	time.Sleep(5 * time.Second)
}

func longTimeJob() {
	time.Sleep(500 * time.Millisecond)
}
