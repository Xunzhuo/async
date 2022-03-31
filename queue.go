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

	workJobsQueue    chan *Job
	waitJobsQueue    chan *Job
	shutdownDataChan chan bool
	shutdownWorkChan chan bool

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
		shutdownDataChan:     make(chan bool),
		shutdownWorkChan:     make(chan bool),
	}
}

// CheckJobs check if Job is valid to add to queue
func (a *Queue) CheckJob(job *Job) (bool, error) {
	if job.handler.Kind() != reflect.Func {
		return false, fmt.Errorf("job has been given a wrong handler, make sure it is a func")
	}

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

	if oldJob := a.GetJobByID(job.jobID); oldJob != nil {
		if subID := job.GetSubID(); subID != keyOfnoSubID {
			oldJob.SetSubID(subID)
		}
		job.subjobIDs = oldJob.subjobIDs
	}

	a.SetJobStatus(job, StatusPending)
	a.workJobsQueue <- job
	a.logger.Info("Add Job to WorkQueue", "jobID", job.jobID, "subID", job.subID, "WorkQueue Length", a.Length())
	time.Sleep(1 * time.Second)
	return true
}

// AddJob Add Job to waitQueue
func (a *Queue) AddJob(job *Job) bool {
	if _, err := a.CheckJob(job); err != nil {
		a.logger.Error(err, "Add Job Denied")
		return false
	}

	if oldJob := a.GetJobByID(job.jobID); oldJob != nil {
		if subID := job.GetSubID(); subID != keyOfnoSubID {
			oldJob.SetSubID(subID)
		}
		job.subjobIDs = oldJob.subjobIDs
	}

	if len(a.waitJobsQueue) == a.maxWaitQueueLength {
		headJob := <-a.waitJobsQueue
		a.workJobsQueue <- headJob
	}

	a.SetJobStatus(job, StatusPending)
	a.waitJobsQueue <- job

	a.logger.Info("Add Job to WaitQueue", "jobID", job.jobID, "WaitQueue Length", a.WaitQueueLength())
	return true
}

func (a *Queue) Run() bool {
	if len(a.waitJobsQueue) == 0 {
		return false
	}

	go func(waitJobs chan *Job) {
		for {
			if len(waitJobs) < 1 {
				return
			}
			headJob := <-waitJobs
			a.workJobsQueue <- headJob
			a.logger.Info("Job added into work queue from wait queue", "jobID", headJob.jobID)
		}
	}(a.waitJobsQueue)

	return true
}

// Start Start the Queue
func (a *Queue) Start() *Queue {
	dataChans := make(chan map[string]interface{}, a.options.maxWorkQueueLength)
	a.startDataPipline(dataChans)
	a.startworkPipline(dataChans)
	return a
}

func (a *Queue) startDataPipline(dataChans chan map[string]interface{}) {
	go func(result map[string]map[string][]interface{}, dataChans chan map[string]interface{}) {
		subData := make(map[string][]interface{})
		for {
			select {
			case <-a.shutdownDataChan:
				a.logger.Info("shutdown DataPipline")
				return
			case res := <-dataChans:

				jobID := res[keyOfJobID].(string)
				subID := res[keyOfSubID].(string)
				jobdata := res[keyOfjobData].([]interface{})

				subData[subID] = jobdata
				result[jobID] = subData

				a.logger.Info("Job Wrote Data", "jobID", jobID, "SubID", subID)
			}
		}
	}(a.sharedJobData, dataChans)
}

func (a *Queue) startworkPipline(dataChans chan map[string]interface{}) {
	go func(jobs chan *Job, dataChans chan map[string]interface{}) {
		for {
			select {
			case <-a.shutdownWorkChan:
				a.logger.Info("shutdown workPipline")
				return
			case job := <-a.workJobsQueue:
				jobID := job.GetJobID()
				subID := job.GetSubID()
				jobData := make([]interface{}, 0)

				values := job.handler.Call(job.params)
				a.SetJobStatus(job, StatusRunning)

				if valuesNum := len(values); valuesNum > 0 {
					resultItems := make([]interface{}, valuesNum)
					for k, v := range values {
						resultItems[k] = v.Interface()
					}
					jobData = resultItems
				}
				dataChans <- map[string]interface{}{keyOfJobID: jobID, keyOfSubID: subID, keyOfjobData: jobData}
				a.logger.Info("WorkQueue Job Sent Data", "jobID", jobID, "SubID", subID)
			}
		}
	}(a.workJobsQueue, dataChans)
}

// IsFull check workqueue if it is full
func (a *Queue) IsFull() bool {
	return a.maxWorkQueueLength <= a.workQueueLength
}

// IsFull check workqueue if it is full
func (a *Queue) Length() int {
	return len(a.workJobsQueue)
}

// IsFull check workqueue if it is full
func (a *Queue) WaitQueueLength() int {
	return len(a.waitJobsQueue)
}

// Stop stop all piplines
func (a *Queue) Stop() {
	a.stopWorkChan()
	a.stopDataChan()
}

// StopWorkChan
func (a *Queue) stopWorkChan() {
	a.shutdownWorkChan <- true
}

// stopDataChan
func (a *Queue) stopDataChan() {
	a.shutdownDataChan <- true
}
