package async

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestCreateWorkQueueWithOptions(t *testing.T) {
	workQueue := NewJobQueue(
		WithMaxWaitQueueLength(10),
		WithMaxWorkQueueLength(10),
	)

	workQueue.Start()
	stop := make(chan bool)

	go func() {
		for {
			select {
			case _, ok := <-stop:
				if !ok {
					return
				}
				return
			default:
				id := fmt.Sprintf("%d", rand.Intn(1000000))
				workQueue.AddJob(NewJob(id, fakeJobV2, "xunzhuo"))
			}
		}
	}()

	time.Sleep(3 * time.Second)
	stop <- true
	log.Info("sent signal and quit")
	close(stop)
	time.Sleep(5 * time.Second)
}

func fakeJobV2(value string) string {
	return "Hello World from " + value
}
