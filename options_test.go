package async

import (
	"testing"
	"time"

	"k8s.io/klog/v2/klogr"
)

var log = klogr.New()

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
				workQueue.AddJobAndRun(NewJob(fakeJobV2, "xunzhuo"))
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
