package async

import "sync"

var statusMap *statuses

func init() {
	statusMap = newDefaultStatuses()
}

// var mutex sync.RWMutex

func newDefaultStatuses() *statuses {
	return &statuses{
		workJobsStatus: make(map[string]*Job),
		safeStatus:     sync.Map{},
	}
}

type statuses struct {
	workJobsStatus map[string]*Job
	safeStatus     sync.Map
}

func (a *Queue) GetJobStatuses() map[string]*Job {
	return statusMap.getJobStatuses()
}

func (a *Queue) GetJobStatus(job *Job) string {
	return statusMap.getStatus(*job)
}

func (a *Queue) HasJob(job Job) bool {
	if key := a.GetJobStatus(&job); key == "" {
		return false
	}
	return true
}

func (a *Queue) SetJobStatus(job *Job, status string) *Queue {
	// a.logger.Info("SetJobStatus", "JobID", job.jobID, "SubID", job.subID, "Status", status)
	statusMap.setStatus(job, status)
	return a
}

func (j *statuses) setStatus(job *Job, status string) {
	queueLockers.locker.Lock()
	defer queueLockers.locker.Unlock()
	job.subjobIDs[job.GetSubID()] = status
	j.workJobsStatus[job.jobID] = job
}

func (j *statuses) getStatus(job Job) string {
	queueLockers.locker.RLock()
	defer queueLockers.locker.RUnlock()
	statuses := j.workJobsStatus
	if status, ok := statuses[job.jobID]; ok {
		if s, ok := status.subjobIDs[job.GetSubID()]; ok {
			return s
		}
		return StatusUnknown
	}
	return StatusUnknown
}

func (j *statuses) getJobStatuses() map[string]*Job {
	return j.workJobsStatus
}
