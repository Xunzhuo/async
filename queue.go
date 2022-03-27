package async

import (
	"fmt"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"
)

// Queue The Async Job WorkQueue
type Queue struct {
	*options

	workJobsQueue chan Job
	waitJobsQueue chan Job

	workQueueLength int

	SharedJobData        map[string]map[string][]interface{}
	SharedJobDataChannel chan map[string]map[string][]interface{}
}

func Q() *Queue {
	return &Queue{
		options:              defaultOptions,
		workJobsQueue:        make(chan Job, DefaultWorkQueueCapacity),
		SharedJobData:        make(map[string]map[string][]interface{}),
		SharedJobDataChannel: make(chan map[string]map[string][]interface{}),
		waitJobsQueue:        make(chan Job, DefaultWaitQueueCapacity),
	}
}

// CheckJobs check if Job is valid to add to queue
func (a *Queue) CheckJob(job *Job) (bool, error) {
	if a.IsFull() {
		return false, fmt.Errorf("work Queue is full")
	}

	log.Debug("checking Job to see if it can be added to queue")
	if a.GetJobStatus(job) == StatusFailure {
		return false, fmt.Errorf("job created failed")
	}

	if !job.handler.IsValid() {
		return false, fmt.Errorf("job handler is nil")
	}

	if a.IsLock(job) {
		return false, fmt.Errorf("job has been locked")
	}

	if job.handler.Kind() != reflect.Func {
		return false, fmt.Errorf("job has been given wrong handler, make sure it is a func")
	}

	if !job.EnableSubjob {
		if e := a.GetJobStatus(job); e == StatusRunning {
			return false, fmt.Errorf("jobs %s has been added", job.JobID)
		}
	} else {
		if e := a.GetJobStatus(job); e == StatusRunning {
			return false, fmt.Errorf("sub jobs has been added")
		}
	}
	return true, nil
}

// AddJobAndRun Add Job to workQueue and run it
func (a *Queue) AddJobAndRun(job *Job) bool {
	if _, err := a.CheckJob(job); err != nil {
		log.Warningf(fmt.Sprintf("Add Job Denied: %v", err.Error()))
		return false
	}

	a.SetJobStatus(job, StatusPending)
	a.workJobsQueue <- *job
	a.workQueueLength++
	log.Info("Add Job to WorkQueue and Run With JobID: ", job.JobID, " Number of Jobs in Queue :", a.workQueueLength)
	return true
}

// AddJob Add Job to waitQueue
func (a *Queue) AddJob(job *Job) bool {
	if _, err := a.CheckJob(job); err != nil {
		log.Info(fmt.Sprintf("Add Job Denied: %v", err.Error()))
		return false
	}

	if len(a.waitJobsQueue) == a.maxWaitQueueLength {
		log.Info(fmt.Sprintf("Async Job %s Has moved into workQueue for space pressure", job.JobID))
		log.Info("Job has moved into work queue with JobID: ", job.JobID)
		headJob := <-a.waitJobsQueue
		a.workJobsQueue <- headJob
	}

	a.waitJobsQueue <- *job

	log.Info(fmt.Sprintf("Job has added into wait work queue with JobID: %s with length %d", job.JobID, len(a.waitJobsQueue)))
	return true
}

func (a *Queue) Run() bool {
	if len(a.waitJobsQueue) == 0 {
		log.Info("wait work queue is empty")
		return false
	}

	go func(waitJobs chan Job) {
		for {
			if len(waitJobs) < 1 {
				break
			}
			headJob := <-waitJobs
			a.workJobsQueue <- headJob
			log.Info(fmt.Sprintf("Job %s has added into work queue from waited queue", headJob.JobID))
		}
	}(a.waitJobsQueue)

	log.Info("All Job has added into work queue from waited queue")
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

			log.Info(fmt.Sprintf("Async Job %s Finished with SubID: %s", jobID, subID))
		}
	}(a.SharedJobData, dataChans)

	go func(jobs chan Job, dataChans chan map[string]interface{}) {
		for {
			job := <-a.workJobsQueue
			jobID := job.JobID
			subID := job.GetSubID()

			jobData := make([]interface{}, 0)
			if a.GetJobStatus(&job) == StatusRunning {
				log.Info(fmt.Sprintf("Async Job %s Has Started", jobID))
				continue
			}

			values := job.handler.Call(job.params)
			a.SetJobStatus(&job, StatusRunning)
			log.Info(fmt.Sprintf("Async Job %s Has Done", jobID))

			if valuesNum := len(values); valuesNum > 0 {
				resultItems := make([]interface{}, valuesNum)
				for k, v := range values {
					resultItems[k] = v.Interface()
				}
				jobData = resultItems
			}

			dataChans <- map[string]interface{}{keyOfJobID: jobID, keyOfSubID: subID, keyOfjobData: jobData}
			log.Info(fmt.Sprintf("Async Job %s Has sent Job Data", jobID))
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
