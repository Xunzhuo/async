package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Xunzhuo/async"
	log "github.com/sirupsen/logrus"
)

func main() {
	workQueue := async.NewJobQueue(
		async.WithMaxWaitQueueLength(100),
		async.WithMaxWorkQueueLength(100),
	)

	workQueue.Start()

	stop := make(chan bool)
	stopData := make(chan bool)
	jobID := make(chan string, 1000)

	go func() {
		for {
			select {
			case _, ok := <-stop:
				if !ok {
					return
				}
				return
			default:
				masterID := fmt.Sprintf("%d", rand.Intn(1000000))
				counter := 3
				if !workQueue.IsFull() {
					for {
						slaveID := fmt.Sprintf("%d", rand.Intn(1000000))
						job := async.NewJob(masterID, longTimeJob, "xunzhuo")
						job.AddSubJob(slaveID)
						if ok := workQueue.AddJobAndRun(job); ok {
							jobID <- masterID
							log.Warning("Send Job ID: ", masterID)
						}
						if counter < 1 {
							break
						}
						counter--
					}
				}
			}
		}
	}()
	time.Sleep(60 * time.Second)
	stop <- true

	log.Info("Send signal to close workQueue")
	close(stop)

	go func() {
		for {
			log.Info("Start to receive Job Data")
			select {
			case _, ok := <-stopData:
				if !ok {
					return
				}
				return
			case job := <-jobID:
				log.Info("Received Job ID: ", job)
				if datas, ok := workQueue.GetJobsData(job); ok {
					for _, data := range datas {
						log.Warningf(fmt.Sprintf("Get data from workQueue %s with ID: %s", data[0].(string), job))
					}
				}
			}
		}
	}()

	time.Sleep(10 * time.Second)
	stopData <- true
	close(stopData)
}

func longTimeJob(value string) string {
	time.Sleep(1000 * time.Millisecond)
	return "Hello World from " + value
}
