package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/Xunzhuo/async"
	"k8s.io/klog/v2/klogr"
)

var log = klogr.New()

func main() {
	workQueue := async.Q().
		SetMaxWaitQueueLength(100).
		SetMaxWorkQueueLength(100)

	workQueue.Start()
	jobs := make(chan async.Job, 1000)
	stop := make(chan bool)

	wg := sync.WaitGroup{}
	count := 10
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			masterID := fmt.Sprintf("%d", rand.Intn(100))
			counter := 3
			if !workQueue.IsFull() {
				for {
					slaveID := fmt.Sprintf("%d", rand.Intn(100))
					job := async.NewJob(longTimeJob, "xunzhuo").SetJobID(masterID)
					job.SetSubID(slaveID)
					if ok := workQueue.AddJobAndRun(job); ok {
						jobs <- *job
					}
					if counter < 1 {
						break
					}
					counter--
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	go func() {
		for {
			select {
			case _, ok := <-stop:
				if !ok {
					return
				}
				return
			case job := <-jobs:
				if datas, ok := workQueue.GetJobsData(job); ok {
					for _, data := range datas {
						log.Info("Get Data", "Data", data[0].(string))
					}
				}
			}
		}
	}()

	time.Sleep(10 * time.Second)
	stop <- true
	close(stop)
}

func longTimeJob(value string) string {
	time.Sleep(100 * time.Millisecond)
	return "Hello World from " + value
}
