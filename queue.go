package async

import (
	"fmt"
	"reflect"
	"time"

	"github.com/go-logr/logr"
)

// Queue The Async Job WorkQueue
type Queue struct {
	*options
	logger logr.Logger

	workJobsQueue chan *Job
	waitJobsQueue chan *Job

	workQueueLength int

	sharedJobData        map[string]map[string][]interface{}
	sharedJobDataChannel chan map[string]map[string][]interface{}
}

func Q() *Queue {
	return &Queue{
		options:              defaultOptions,
		logger:               defaultOptions.logger,
		workJobsQueue:        make(chan *Job, defaultOptions.maxWorkQueueLength),
		sharedJobData:        make(map[string]map[string][]interface{}),
		sharedJobDataChannel: make(chan map[string]map[string][]interface{}),
		waitJobsQueue:        make(chan *Job, defaultOptions.maxWorkQueueLength),
	}
}

// CheckJobs check if Job is valid to add to queue
func (a *Queue) CheckJob(job *Job) (bool, error) {
	if a.IsFull() {
		return false, fmt.Errorf("work Queue is full")
	}

	if a.GetJobStatus(job) == StatusFailure {
		return false, fmt.Errorf("job created failed")
	}

	// if a.GetJobStatus(job) == StatusUnknown {
	// 	return false, fmt.Errorf("unknown job")
	// }

	if !job.handler.IsValid() {
		return false, fmt.Errorf("job handler is nil")
	}

	if a.IsLock(job) {
		return false, fmt.Errorf("job has been locked")
	}

	if job.handler.Kind() != reflect.Func {
		return false, fmt.Errorf("job has been given a wrong handler, make sure it is a func")
	}

	if e := a.GetJobStatus(job); e == StatusRunning {
		return false, fmt.Errorf("jobs %s has been added", job.jobID)
	}

	return true, nil
}

// AddJobAndRun Add Job to workQueue and run it
func (a *Queue) AddJobAndRun(job *Job) bool {
	if _, err := a.CheckJob(job); err != nil {
		a.logger.Error(err, "Add Job Denied")
		return false
	}

	a.SetJobStatus(job, StatusPending)
	a.workJobsQueue <- job
	a.workQueueLength++
	a.logger.Info("Add Job to WorkQueue and Run", "jobID", job.jobID, "Queue Length", a.workQueueLength)
	return true
}

// AddJob Add Job to waitQueue
func (a *Queue) AddJob(job *Job) bool {
	if _, err := a.CheckJob(job); err != nil {
		a.logger.Error(err, "Add Job Denied")
		return false
	}

	if len(a.waitJobsQueue) == a.maxWaitQueueLength {
		a.logger.Info("Job moved into work queue", "jobID", job.jobID)
		headJob := <-a.waitJobsQueue
		a.workJobsQueue <- headJob
	}

	a.waitJobsQueue <- job

	a.logger.Info("Add Job to WaitQueue", "jobID", job.jobID, "Queue Length", len(a.waitJobsQueue))
	return true
}

func (a *Queue) Run() bool {
	if len(a.waitJobsQueue) == 0 {
		a.logger.Info("wait work queue is empty")
		return false
	}

	go func(waitJobs chan *Job) {
		for {
			if len(waitJobs) < 1 {
				break
			}
			headJob := <-waitJobs
			a.workJobsQueue <- headJob
			a.logger.Info("Job added into work queue from wait queue", "jobID", headJob.jobID)
		}
	}(a.waitJobsQueue)

	a.logger.Info("All Job added to work queue from wait queue")
	return true
}

// Start Start the Queue
func (a *Queue) Start() *Queue {
	dataChans := make(chan map[string]interface{}, a.maxWorkQueueLength)

	go func(result map[string]map[string][]interface{}, dataChans chan map[string]interface{}) {
		subData := make(map[string][]interface{})
		for {
			res := <-dataChans
			a.workQueueLength--
			jobID := res[keyOfJobID].(string)
			subID := res[keyOfSubID].(string)
			jobdata := res[keyOfjobData].([]interface{})

			subData[subID] = jobdata
			result[jobID] = subData

			a.logger.Info("Job Wrote Data", "jobID", jobID, "SubID", subID)
		}
	}(a.sharedJobData, dataChans)

	go func(jobs chan *Job, dataChans chan map[string]interface{}) {
		for {
			job := <-a.workJobsQueue
			jobID := job.jobID
			subID := job.GetSubID()

			jobData := make([]interface{}, 0)
			if a.GetJobStatus(job) == StatusRunning {
				a.logger.Info("Job Started", "jobID", jobID, "SubID", subID)
				continue
			}

			values := job.handler.Call(job.params)
			a.SetJobStatus(job, StatusRunning)
			a.logger.Info("Job Done", "jobID", jobID, "SubID", subID)

			if valuesNum := len(values); valuesNum > 0 {
				resultItems := make([]interface{}, valuesNum)
				for k, v := range values {
					resultItems[k] = v.Interface()
				}
				jobData = resultItems
			}

			dataChans <- map[string]interface{}{keyOfJobID: jobID, keyOfSubID: subID, keyOfjobData: jobData}
			a.logger.Info("Job Sent Data", "jobID", jobID, "SubID", subID)
		}
	}(a.workJobsQueue, dataChans)
	return a
}

func (a *Queue) WaitForTime(waitTime time.Duration) *Queue {
	time.Sleep(waitTime)
	return a
}

// IsFull check workqueue if it is full
func (a *Queue) IsFull() bool {
	return a.maxWorkQueueLength <= a.workQueueLength
}
