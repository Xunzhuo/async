package async

import (
	"fmt"
	"reflect"

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

	sharedJobData map[string]map[string][]interface{}
}

func Q() *Queue {
	return &Queue{
		options:          defaultOptions,
		logger:           defaultOptions.logger,
		workJobsQueue:    make(chan *Job, defaultOptions.maxWorkQueueLength),
		sharedJobData:    make(map[string]map[string][]interface{}),
		waitJobsQueue:    make(chan *Job, defaultOptions.maxWorkQueueLength),
		shutdownDataChan: make(chan bool),
		shutdownWorkChan: make(chan bool),
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

	if e := a.GetJobStatus(job); e == StatusRunning || e == StatusPending {
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
		job.subjobIDs = oldJob.subjobIDs
	}

	a.SetJobStatus(job, StatusRunning)
	a.logger.Info("Add Job to WorkQueue", "jobID", job.GetJobID(), "subID", job.GetSubID(), "WorkQueue Length", a.Length())
	a.workJobsQueue <- job
	return true
}

// AddJob Add Job to waitQueue
func (a *Queue) AddJob(job *Job) bool {
	if _, err := a.CheckJob(job); err != nil {
		a.logger.Error(err, "Add Job Denied")
		return false
	}

	if oldJob := a.GetJobByID(job.jobID); oldJob != nil {
		job.subjobIDs = oldJob.subjobIDs
	}

	if len(a.waitJobsQueue) == a.maxWaitQueueLength {
		headJob := <-a.waitJobsQueue
		a.workJobsQueue <- headJob
	}

	a.SetJobStatus(job, StatusPending)
	a.waitJobsQueue <- job
	a.logger.Info("Add Job to WaitQueue", "jobID", job.GetSubID(), "subID", job.GetSubID(), "WaitQueue Length", a.WaitQueueLength())
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
			a.SetJobStatus(headJob, StatusRunning)
			a.workJobsQueue <- headJob
			a.logger.Info("Job added into work queue from wait queue", "jobID", headJob.GetJobID(), "subID", headJob.GetSubID())
		}
	}(a.waitJobsQueue)

	return true
}

// Start Start the Queue
func (a *Queue) Start() *Queue {
	dataChans := make(chan map[string]interface{}, a.options.maxWorkQueueLength)
	go a.startWorkPipline(dataChans)
	go a.startDataPipline(dataChans)
	return a
}

func (a *Queue) startDataPipline(dataChans chan map[string]interface{}) {
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
			a.sharedJobData[jobID] = subData

			a.logger.Info("DataPipline Done", "jobID", jobID, "SubID", subID)
		}
	}
}

func (a *Queue) startWorkPipline(dataChans chan map[string]interface{}) {
	for {
		select {
		case <-a.shutdownWorkChan:
			a.logger.Info("shutdown WorkPipline")
			return
		case job := <-a.workJobsQueue:
			jobID := job.GetJobID()
			subID := job.GetSubID()
			jobData := make([]interface{}, 0)

			values := job.handler.Call(job.params)

			if valuesNum := len(values); valuesNum > 0 {
				resultItems := make([]interface{}, valuesNum)
				for k, v := range values {
					resultItems[k] = v.Interface()
				}
				jobData = resultItems
			}
			dataChans <- map[string]interface{}{keyOfJobID: jobID, keyOfSubID: subID, keyOfjobData: jobData}
			a.logger.Info("WorkPipline Done", "jobID", jobID, "SubID", subID)
		}
	}
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
