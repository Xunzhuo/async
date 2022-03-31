package main

import (
	"time"

	"github.com/Xunzhuo/async"
	"k8s.io/klog/v2/klogr"
)

var log = klogr.New()

func main() {

	async.DefaultAsyncQueue().Start()

	stop := make(chan bool)
	stopData := make(chan bool)
	jobID := make(chan async.Job, 1000)

	go func() {
		for {
			select {
			case _, ok := <-stop:
				if !ok {
					return
				}
				return
			default:
				job := async.NewJob(longTimeJob, "xunzhuo")
				if ok := async.DefaultAsyncQueue().AddJobAndRun(job); ok {
					jobID <- *job
					log.Info("Send Job", "JobID", job.GetJobID())
				} else {
					log.Info("Reject Job", "JobID", job.GetJobID())
				}
			}
		}
	}()

	time.Sleep(5 * time.Second)
	stop <- true
	close(stop)

	go func() {
		for {
			time.Sleep(1000 * time.Millisecond)
			select {
			case _, ok := <-stopData:
				if !ok {
					return
				}
				return
			case job := <-jobID:
				log.Info("Received Job ID", "JobID", job.GetJobID())
				if data, ok := async.DefaultAsyncQueue().GetJobData(job); ok {
					log.Info("Get data from WorkQueue", "Data", data[0].(string), "JobID", job.GetJobID())
				}
			}
		}
	}()

	time.Sleep(100 * time.Second)
	stopData <- true
	close(stopData)
}

func longTimeJob(value string) string {
	return "Hello World from " + value
}
