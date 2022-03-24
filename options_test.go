package async

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestCreateWorkQueueWithOptions(t *testing.T) {
	workQueue := NewJobQueue(WithMaxWaitQueueLength(10), WithMaxWorkQueueLength(10))
	workQueue.Start()
	stop := make(chan bool, 1)

	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			select {
			case _, ok := <-stop:
				if !ok {
					log.Info("channel is closed, just quit")
				}
				log.Info("received signal and quit")
				return
			default:
				id := fmt.Sprintf("%d", rand.Intn(1000000))
				workQueue.AddJob(NewJob(id, fakeJobV2, url))
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
	return "Hello World V2" + value
}
