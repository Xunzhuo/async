package async

var statusMap *statuses

func init() {
	statusMap = newDefaultStatuses()
}

func newDefaultStatuses() *statuses {
	return &statuses{
		workJobsStatus: make(map[string]*Job),
	}
}

type statuses struct {
	workJobsStatus map[string]*Job
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
	statusMap.setStatus(job, status)
	return a
}

func (j *statuses) setStatus(job *Job, status string) {
	queueLocker.locker.Lock()
	defer queueLocker.locker.Unlock()
	job.subjobIDs[job.GetSubID()] = status
	j.workJobsStatus[job.jobID] = job
}

func (j *statuses) getStatus(job Job) string {
	return j.getJobStatuses()[job.jobID].subjobIDs[job.GetSubID()]
}

func (j *statuses) getJobStatuses() map[string]*Job {
	return j.workJobsStatus
}
