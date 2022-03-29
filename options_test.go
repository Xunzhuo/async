package async

import (
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestCreateWorkQueueWithOptions(t *testing.T) {
	workQueue := Q().
		SetMaxWaitQueueLength(100).
		SetMaxWorkQueueLength(100).
		Start()

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
				time.Sleep(100 * time.Millisecond)
				// id := fmt.Sprintf("%d", rand.Intn(1000000))
				if job, err := NewJob(fakeJobV2, "xunzhuo"); err != nil {
					t.Error(err)
				} else {
					workQueue.AddJobAndRun(job)
				}

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
