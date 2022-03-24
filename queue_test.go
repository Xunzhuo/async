package async

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

const url = "https://baidu.com"

func TestAsyncWithAddTaskAndRun(t *testing.T) {
	count := 10
	jobs := make([]string, 0)

	Engine.Start()

	for {
		count--
		jobID := fmt.Sprintf("%d", rand.Intn(100000))
		jobs = append(jobs, jobID)
		Engine.AddTaskAndRun(NewJob(jobID, sendRequest, url))
		if count < 1 {
			break
		}
	}
	time.Sleep(5 * time.Second)
	for _, job := range jobs {
		_, ok := Engine.GetJobData(job)
		log.Printf("GetJobData %t With JobID %s", ok, job)
	}
}

func TestAsyncWithAddTaskAfterRun(t *testing.T) {
	count := 60
	jobs := make([]string, 0)

	Engine.Start()
	Engine.SetMaxWaitQueueLength(50)

	for {
		count--
		jobID := fmt.Sprintf("%d", rand.Intn(100000))
		jobs = append(jobs, jobID)
		Engine.AddTask(NewJob(jobID, sendRequest, url))
		if count < 1 {
			break
		}
	}

	Engine.Run()

	time.Sleep(5 * time.Second)
	for _, job := range jobs {
		_, ok := Engine.GetJobData(job)
		log.Printf("GetJobData %t With JobID %s", ok, job)
	}
}

func TestAsyncWithAddBlindTaskAndRun(t *testing.T) {
	count := 10

	Engine.Start()

	for {
		count--
		Engine.AddTaskAndRun(NewBlindJob(sendRequest, url))
		if count < 1 {
			break
		}
	}

	time.Sleep(5 * time.Second)

	for jobD := range Engine.GetAllJobID() {
		_, ok := Engine.GetJobData(jobD)
		log.Printf("GetJobData %t With JobID %s", ok, jobD)
	}
}

func TestAsyncWithAddBlindTaskAfterRun(t *testing.T) {
	count := 60

	Engine.Start()
	Engine.SetMaxWaitQueueLength(50)

	for {
		count--
		Engine.AddTask(NewBlindJob(sendRequest, url))
		if count < 1 {
			break
		}
	}

	Engine.Run()

	time.Sleep(5 * time.Second)

	for jobD := range Engine.GetAllJobID() {
		_, ok := Engine.GetJobData(jobD)
		log.Printf("GetJobData %t With JobID %s", ok, jobD)
	}
}

func TestAsyncWithSameID(t *testing.T) {
	count := 10
	jobs := make([]string, 0)

	Engine.Start()

	for {
		count--
		jobID := "id"
		jobs = append(jobs, jobID)
		Engine.AddTaskAndRun(NewJob(jobID, sendRequest, url))
		if count < 1 {
			break
		}
	}
	time.Sleep(5 * time.Second)
	for _, job := range jobs {
		_, ok := Engine.GetJobData(job)
		log.Printf("GetJobData %t With JobID %s", ok, job)
	}
}

func TestAsyncWithSubJobs(t *testing.T) {
	count := 10
	jobs := make([]string, 0)

	Engine.Start()

	for {
		count--
		jobID := "master"
		subID := fmt.Sprintf("%d", rand.Intn(100000))
		jobs = append(jobs, jobID)
		masterJob := NewJob(jobID, sendRequest, url)
		masterJob.AddSubJob(subID)
		Engine.AddTaskAndRun(masterJob)

		if count < 1 {
			break
		}
	}
	time.Sleep(5 * time.Second)
	for _, job := range jobs {
		_, ok := Engine.GetJobData(job)
		log.Printf("GetJobData %t With JobID %s", ok, job)
	}
}

func sendRequest(url string) string {
	if rs, err := http.Get(url); err == nil {
		defer rs.Body.Close()
		if res, err := ioutil.ReadAll(rs.Body); err == nil {
			return string(res)
		}
	}
	return ""
}
