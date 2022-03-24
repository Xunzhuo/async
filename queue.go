package async

import (
	"fmt"
	"reflect"

	log "github.com/sirupsen/logrus"
)

// JobWorkQueue The Async Job WorkQueue
type JobWorkQueue struct {
	workJobsQueue   chan Job
	waitJobsQueue   chan Job
	workJobsStatus  map[string]map[string]string
	workJobIDHisory map[string][]string
	lockJobIDList   map[string]bool

	workQueueLength    int
	maxWorkQueueLength int
	maxWaitQueueLength int

	SharedJobData        map[string]map[string][]interface{}
	SharedJobDataChannel chan map[string]map[string][]interface{}
}

// NewJobQueue create a Empty JobWorkQueue
func NewJobQueue(opts ...Option) JobWorkQueue {
	opt := new(JobWorkQueueOptions)
	for _, o := range opts {
		o(opt)
	}
	if opt.MaxWaitQueueLength == 0 {
		opt.MaxWaitQueueLength = 100
	}
	if opt.MaxWorkQueueLength == 0 {
		opt.MaxWorkQueueLength = 100
	}
	return JobWorkQueue{
		maxWaitQueueLength:   opt.MaxWaitQueueLength,
		maxWorkQueueLength:   opt.MaxWorkQueueLength,
		workJobsQueue:        make(chan Job, opt.MaxWorkQueueLength),
		workJobsStatus:       make(map[string]map[string]string),
		workJobIDHisory:      make(map[string][]string),
		SharedJobData:        make(map[string]map[string][]interface{}),
		SharedJobDataChannel: make(chan map[string]map[string][]interface{}),
		lockJobIDList:        make(map[string]bool),
		waitJobsQueue:        make(chan Job, opt.MaxWaitQueueLength),
	}
}

func (a *JobWorkQueue) SetMaxWorkQueueLength(len int) {
	a.maxWorkQueueLength = len
}

func (a *JobWorkQueue) SetMaxWaitQueueLength(len int) {
	a.maxWaitQueueLength = len
}

// CheckJobs check if Job is valid to add to queue
func (a *JobWorkQueue) CheckJob(job Job) (bool, error) {
	if a.workQueueLength == a.maxWorkQueueLength {
		return false, fmt.Errorf("work Queue is full")
	}

	log.Debug("checking Job to see if it can be added to queue")
	if job.Status == StatusFailure {
		return false, fmt.Errorf("job created failed")
	}

	if a.IsLock(job) {
		return false, fmt.Errorf("job has been locked")
	}

	if job.Handler.Kind() != reflect.Func {
		return false, fmt.Errorf("job has been given wrong handler, make sure it is a func")
	}

	if !job.EnableSubjob {
		if _, e := a.workJobsStatus[job.JobID][job.GetSubID()]; e {
			return false, fmt.Errorf("jobs %s has been added", job.JobID)
		}
	} else {
		if _, e := a.workJobsStatus[job.JobID][job.GetSubID()]; e {
			return false, fmt.Errorf("sub jobs has been added")
		}
	}
	return true, nil
}

// AddJobAndRun Add Job to workQueue and run it
func (a *JobWorkQueue) AddJobAndRun(job Job) bool {
	if _, err := a.CheckJob(job); err != nil {
		log.Warningf(fmt.Sprintf("Add Job Denied: %v", err.Error()))
		return false
	}

	subID := map[string]string{job.GetSubID(): StatusPending}
	a.workJobsStatus[job.JobID] = subID
	a.workJobsQueue <- job
	a.workQueueLength++
	log.Info("Add Job to WorkQueue and Run With JobID: ", job.JobID, " Number of Jobs in Queue :", a.workQueueLength)
	return true
}

// AddJob Add Job to waitQueue
func (a *JobWorkQueue) AddJob(job Job) bool {
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

	a.waitJobsQueue <- job

	log.Info(fmt.Sprintf("Job has added into wait work queue with JobID: %s with length %d", job.JobID, len(a.waitJobsQueue)))
	return true
}

func (a *JobWorkQueue) Run() bool {
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

// Start Start the JobWorkQueue
func (a *JobWorkQueue) Start() {
	dataChans := make(chan map[string]interface{}, a.maxWorkQueueLength)

	go func(result map[string]map[string][]interface{}, dataChans chan map[string]interface{}) {
		subData := make(map[string][]interface{})
		for {
			res := <-dataChans
			a.workQueueLength--
			jobID := res[JOBID].(string)
			subID := res[SUBID].(string)
			jobdata := res[JOBDATA].([]interface{})

			subData[subID] = jobdata
			result[jobID] = subData

			log.Info(fmt.Sprintf("Async Job %s Finished with SubID: %s", res[JOBID].(string), res[SUBID].(string)))
		}
	}(a.SharedJobData, dataChans)

	go func(jobs chan Job, dataChans chan map[string]interface{}) {
		for {
			job := <-a.workJobsQueue
			jobID := job.JobID
			subID := job.GetSubID()

			jobData := make([]interface{}, 0)
			if job.Status == StatusRunning {
				log.Info(fmt.Sprintf("Async Job %s Has Started", jobID))
				continue
			}

			values := job.Handler.Call(job.Params)
			job.Status = StatusRunning
			log.Info(fmt.Sprintf("Async Job %s Has Done", jobID))

			if valuesNum := len(values); valuesNum > 0 {
				resultItems := make([]interface{}, valuesNum)
				for k, v := range values {
					resultItems[k] = v.Interface()
				}
				jobData = resultItems
				a.workJobIDHisory[jobID] = append(a.workJobIDHisory[jobID], subID)
			}

			dataChans <- map[string]interface{}{JOBID: jobID, SUBID: subID, JOBDATA: jobData}
			log.Info(fmt.Sprintf("Async Job %s Has sent Job Data", jobID))
		}
	}(a.workJobsQueue, dataChans)
}
