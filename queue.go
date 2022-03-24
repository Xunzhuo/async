package async

import (
	"fmt"
	"reflect"
	"sync"

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
func NewJobQueue() JobWorkQueue {
	return JobWorkQueue{
		maxWaitQueueLength:   100,
		maxWorkQueueLength:   100,
		workJobsQueue:        make(chan Job, 100),
		workJobsStatus:       make(map[string]map[string]string),
		workJobIDHisory:      make(map[string][]string),
		SharedJobData:        make(map[string]map[string][]interface{}),
		SharedJobDataChannel: make(chan map[string]map[string][]interface{}),
		lockJobIDList:        make(map[string]bool),
		waitJobsQueue:        make(chan Job, 100),
	}
}

// CheckJobs check if Job is valid to add to queue
func (a *JobWorkQueue) CheckJob(job Job) (bool, error) {
	log.Debug("checking Job to see if it can be added to queue")

	if a.IsLock(job) {
		return false, fmt.Errorf("job has been locked")
	}

	if job.Handler.Kind() != reflect.Func {
		return false, fmt.Errorf("job has been given wrong handler, make sure it is a func")
	}

	if !job.EnableSubjob {
		if _, e := a.workJobsStatus[job.JobID][NaN]; e {
			return false, fmt.Errorf("jobs %s has been added", job.JobID)
		}
	} else {
		if _, e := a.workJobsStatus[job.JobID][job.NewSubJobID]; e {
			return false, fmt.Errorf("sub jobs has been added")
		}
	}
	return true, nil
}

// AddTaskAndRun Add Job to workQueue and run it
func (a *JobWorkQueue) AddTaskAndRun(job Job) bool {
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

// AddTask Add Job to waitQueue
func (a *JobWorkQueue) AddTask(job Job) bool {
	if _, err := a.CheckJob(job); err != nil {
		log.Info(fmt.Sprintf("Add Job Denied: %v", err.Error()))
		return false
	}

	if len(a.waitJobsQueue) == a.maxWaitQueueLength {
		log.Info("Job has added into wait work queue with JobID: ", job.JobID)
		headJob := <-a.waitJobsQueue
		a.workJobsQueue <- headJob
	}

	a.waitJobsQueue <- job
	log.Info("Job has added into wait work queue with JobID: ", job.JobID)
	return true
}

func (a *JobWorkQueue) Run() bool {
	if len(a.waitJobsQueue) == 0 {
		log.Info("Wait work queue is empty")
		return false
	}

	go func(waitJobs chan Job) {
		for {
			if len(waitJobs) < 1 {
				break
			}
			headJob := <-waitJobs
			a.workJobsQueue <- headJob
			log.Info("Job has added into work queue from waited queue with JobID: ", headJob.JobID)
		}
	}(a.waitJobsQueue)

	log.Info("All Job has added into work queue from waited queue")
	return true
}

func (a *JobWorkQueue) LockJob(job Job) {
	var lock sync.Mutex
	lock.Lock()
	log.Warning("Lock Job with JobID: ", job.JobID)
	a.lockJobIDList[job.JobID] = true
	lock.Unlock()
}

func (a *JobWorkQueue) IsLock(job Job) bool {
	log.Debug("Checking if job is locked with JobID: ", job.JobID)
	if _, ok := a.lockJobIDList[job.JobID]; ok {
		return true
	}
	return false
}

func (a *JobWorkQueue) UnLockJob(job *Job) {
	var lock sync.Mutex
	lock.Lock()
	log.Info("UnLock Job with JobID: ", job.JobID)
	job.Lock = false
	lock.Unlock()
}

// Start Start the JobWorkQueue
func (a *JobWorkQueue) Start() {
	dataChans := make(chan map[string]interface{}, a.maxWorkQueueLength)

	go func(result map[string]map[string][]interface{}, dataChans chan map[string]interface{}) {
		for {
			res := <-dataChans
			a.workQueueLength--
			subData := make(map[string][]interface{})
			subData[res[SUBID].(string)] = res[JOBDATA].([]interface{})
			result[res[JOBID].(string)] = subData
			log.Info("Async Job Finished with JobID: ", res[JOBID].(string), " subID: ", res[SUBID].(string))
		}
	}(a.SharedJobData, dataChans)

	go func(jobs chan Job, dataChans chan map[string]interface{}) {
		for {
			task := <-a.workJobsQueue
			jobID := task.JobID
			subID := task.GetSubID()

			jobData := make([]interface{}, 0)
			if task.Status == StatusRunning {
				log.Info("Async Job Has Started with JobID: ", jobID)
				continue
			}

			values := task.Handler.Call(task.Params)
			task.Status = StatusRunning
			log.Info("Start Async Job with JobID: ", jobID)

			if valuesNum := len(values); valuesNum > 0 {
				resultItems := make([]interface{}, valuesNum)
				for k, v := range values {
					resultItems[k] = v.Interface()
				}
				jobData = resultItems
				a.workJobIDHisory[jobID] = append(a.workJobIDHisory[jobID], subID)
			}

			dataChans <- map[string]interface{}{JOBID: jobID, SUBID: subID, JOBDATA: jobData}
			log.Info("Send Async Data with JobID: ", jobID)
		}
	}(a.workJobsQueue, dataChans)
}

// CleanJobData
func (a *JobWorkQueue) CleanJobData(jobID string) {
	delete(a.SharedJobData, jobID)
}

// CleanJobDatas
func (a *JobWorkQueue) CleanJobDatas(jobIDList ...string) {
	for _, jobID := range jobIDList {
		a.CleanJobData(jobID)
	}
}

// Stop
func (a *JobWorkQueue) Stop() {
	close(a.workJobsQueue)
}

// HasJobKey
func (a *JobWorkQueue) HasJobKey(key string) bool {
	return len(a.workJobsStatus[key]) != 0
}

func (a *JobWorkQueue) GetJobData(jobID string) ([][]interface{}, bool) {
	if len(a.SharedJobData) == 0 {
		return nil, false
	}

	var jobDataList [][]interface{}
	jobDataList = make([][]interface{}, 0)

	if _, ok := a.workJobIDHisory[jobID]; ok {
		for _, subID := range a.workJobIDHisory[jobID] {
			if _, ok := a.SharedJobData[jobID][subID]; ok {
				jobDataList = append(jobDataList, a.SharedJobData[jobID][subID])
			}
		}
	} else {
		return nil, false
	}

	return jobDataList, true
}

func (a *JobWorkQueue) SetMaxWorkQueueLength(len int) {
	a.maxWorkQueueLength = len
}

func (a *JobWorkQueue) SetMaxWaitQueueLength(len int) {
	a.maxWaitQueueLength = len
}
