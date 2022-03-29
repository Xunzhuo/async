package main

import (
	"time"

	"github.com/Xunzhuo/async"
	log "github.com/sirupsen/logrus"
)

func main() {
	queue := async.Q().
		SetMaxWaitQueueLength(100).
		SetMaxWorkQueueLength(100).
		Start()
	job, err := async.NewJob(longTimeJob)
	if err != nil {
		return
	}
	queue.AddJobAndRun(job)
	time.Sleep(5 * time.Second)
}

func longTimeJob() {
	log.Info("Running long time job")
	time.Sleep(500 * time.Millisecond)
}
