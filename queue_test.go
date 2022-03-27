package async

// import (
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"net/http"
// 	"testing"
// 	"time"
// )

// const url = "https://baidu.com"

// func TestAsyncWithAddTaskAndRun(t *testing.T) {
// 	count := 10
// 	jobs := make([]string, 0)

// 	DefaultAsyncQueue().Start()

// 	for {
// 		count--
// 		jobID := fmt.Sprintf("%d", rand.Intn(100000))
// 		jobs = append(jobs, jobID)
// 		DefaultAsyncQueue().AddJobAndRun(NewJob(jobID, sendRequest, url))
// 		if count < 1 {
// 			break
// 		}
// 	}
// 	time.Sleep(5 * time.Second)
// 	for _, job := range jobs {
// 		data, ok := DefaultAsyncQueue().GetJobData(job)
// 		if ok {
// 			log.Printf("GetJobData %t With JobID %s of Data %s", ok, job, data[0].(string))
// 		}
// 	}
// }

// func TestAsyncWithAddTaskAfterRun(t *testing.T) {
// 	count := 10
// 	jobs := make([]string, 0)

// 	DefaultAsyncQueue().Start()
// 	DefaultAsyncQueue().SetMaxWaitQueueLength(50)

// 	for {
// 		count--
// 		jobID := fmt.Sprintf("%d", rand.Intn(100000))
// 		jobs = append(jobs, jobID)
// 		DefaultAsyncQueue().AddJob(NewJob(jobID, sendRequest, url))
// 		if count < 1 {
// 			break
// 		}
// 	}

// 	DefaultAsyncQueue().Run()
// 	time.Sleep(5 * time.Second)

// 	for _, job := range jobs {
// 		data, ok := DefaultAsyncQueue().GetJobData(job)
// 		if ok {
// 			log.Printf("GetJobData %t With JobID %s of Data %s", ok, job, data[0].(string))
// 		}
// 	}
// }

// func TestAsyncWithAddBlindTaskAndRun(t *testing.T) {
// 	count := 10

// 	DefaultAsyncQueue().Start()

// 	for {
// 		count--
// 		DefaultAsyncQueue().AddJobAndRun(NewBlindJob(sendRequest, url))
// 		if count < 1 {
// 			break
// 		}
// 	}

// 	time.Sleep(5 * time.Second)

// 	for jobD := range DefaultAsyncQueue().GetAllJobID() {
// 		data, ok := DefaultAsyncQueue().GetJobData(jobD)
// 		if ok {
// 			log.Printf("GetJobData %t With JobID %s of Data %s", ok, jobD, data[0].(string))
// 		}
// 	}
// }

// func TestAsyncWithAddBlindTaskAfterRun(t *testing.T) {
// 	count := 60

// 	DefaultAsyncQueue().Start()
// 	DefaultAsyncQueue().SetMaxWaitQueueLength(50)

// 	for {
// 		count--
// 		DefaultAsyncQueue().AddJob(NewBlindJob(sendRequest, url))
// 		if count < 1 {
// 			break
// 		}
// 	}

// 	DefaultAsyncQueue().Run()

// 	time.Sleep(5 * time.Second)

// 	for jobD := range DefaultAsyncQueue().GetAllJobID() {
// 		data, ok := DefaultAsyncQueue().GetJobData(jobD)
// 		if ok {
// 			log.Printf("GetJobData %t With JobID %s of Data %s", ok, jobD, data[0].(string))
// 		}
// 	}
// }

// func TestAsyncWithSameID(t *testing.T) {
// 	count := 10
// 	jobs := make([]string, 0)

// 	DefaultAsyncQueue().Start()

// 	for {
// 		count--
// 		jobID := "sameID"
// 		jobs = append(jobs, jobID)
// 		DefaultAsyncQueue().AddJobAndRun(NewJob(jobID, sendRequest, url))
// 		if count < 1 {
// 			break
// 		}
// 	}
// 	time.Sleep(5 * time.Second)
// 	for _, job := range jobs {
// 		data, ok := DefaultAsyncQueue().GetJobData(job)
// 		if ok {
// 			log.Printf("GetJobData %t With JobID %s of Data %s", ok, job, data[0].(string))
// 		}
// 	}
// }

// func TestAsyncWithSubJobs(t *testing.T) {
// 	count := 10
// 	var jobID string

// 	DefaultAsyncQueue().Start()

// 	for {
// 		count--
// 		jobID = "master"
// 		subID := fmt.Sprintf("%d", rand.Intn(100000))
// 		masterJob := NewJob(jobID, sendRequest, url)
// 		masterJob.AddSubJob(subID)
// 		DefaultAsyncQueue().AddJobAndRun(masterJob)

// 		if count < 1 {
// 			break
// 		}
// 	}

// 	time.Sleep(5 * time.Second)

// 	datas, ok := DefaultAsyncQueue().GetJobsData(jobID)
// 	for subID, data := range datas {
// 		log.Printf("GetJobsData %t With subID %s with Data %s", ok, subID, data[0].(string))
// 	}
// }

// func TestAsyncWithSubJobData(t *testing.T) {
// 	count := 10
// 	var jobID string

// 	DefaultAsyncQueue().Start()

// 	for {
// 		count--
// 		jobID = "master"
// 		subID := fmt.Sprintf("%d", rand.Intn(100000))
// 		masterJob := NewJob(jobID, sendRequest, url)
// 		masterJob.AddSubJob(subID)
// 		DefaultAsyncQueue().AddJobAndRun(masterJob)

// 		if count < 1 {
// 			break
// 		}
// 	}

// 	time.Sleep(5 * time.Second)

// 	subIDs := DefaultAsyncQueue().GetJobSubID(jobID)
// 	for _, subID := range subIDs {
// 		if data, ok := DefaultAsyncQueue().GetSubJobData(jobID, subID); ok {
// 			log.Printf("GetJobsData %t With subID %s with Data %s", ok, subID, data[0].(string))
// 		}
// 	}
// }

// func TestAsyncWithSubJobsID(t *testing.T) {
// 	count := 10
// 	jobs := make([]string, 0)

// 	DefaultAsyncQueue().Start()

// 	for {
// 		count--
// 		jobID := "master"
// 		subID := fmt.Sprintf("%d", rand.Intn(100000))
// 		jobs = append(jobs, jobID)
// 		masterJob := NewJob(jobID, sendRequest, url)
// 		masterJob.AddSubJob(subID)
// 		DefaultAsyncQueue().AddJobAndRun(masterJob)

// 		if count < 1 {
// 			break
// 		}
// 	}

// 	time.Sleep(5 * time.Second)
// 	for _, job := range jobs {
// 		subIDs := DefaultAsyncQueue().GetJobSubID(job)
// 		for _, subID := range subIDs {
// 			log.Printf("GetJobsData subID %s with jobID %s", subID, job)
// 		}
// 	}
// }

// func sendRequest(url string) string {
// 	var msg string
// 	randNum := rand.Intn(100000)
// 	if rs, err := http.Get(url); err == nil {
// 		defer rs.Body.Close()
// 		msg = fmt.Sprintf("send request %d to %s with resp code %d", randNum, url, rs.StatusCode)
// 	}
// 	return msg
// }

// func TestAsyncWithParams(t *testing.T) {
// 	count := 10
// 	jobs := make([]string, 0)

// 	DefaultAsyncQueue().Start()

// 	for {
// 		count--
// 		jobID := fmt.Sprintf("%d", rand.Intn(100000))
// 		jobs = append(jobs, jobID)
// 		DefaultAsyncQueue().AddJobAndRun(NewJob(jobID, fakeJob, url, jobID))
// 		if count < 1 {
// 			break
// 		}
// 	}

// 	time.Sleep(5 * time.Second)
// 	for _, job := range jobs {
// 		data, ok := DefaultAsyncQueue().GetJobData(job)
// 		if ok {
// 			log.Printf("GetJobData %t With JobID %s of Data %s",
// 				ok, job, data[0].(string))
// 		}
// 	}
// }

// func fakeJob(value string) string {
// 	return "Hello World" + value
// }
