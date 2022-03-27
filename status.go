package async

var statusMap *statuses

func init() {
	statusMap = newDefaultStatuses()
}

func newDefaultStatuses() *statuses {
	return &statuses{
		workJobsStatus: make(map[string]map[string]string),
	}
}

type statuses struct {
	workJobsStatus map[string]map[string]string
}

func (a *Queue) GetJobStatuses() map[string]map[string]string {
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
	statusMap.setStatus(*job, status)
	return a
}

func (j *statuses) setStatus(job Job, status string) {
	subJobStatus := make(map[string]string)
	subJobStatus[job.GetSubID()] = status
	j.workJobsStatus[job.JobID] = subJobStatus
}

func (j *statuses) getStatus(job Job) string {
	return j.getJobStatuses()[job.JobID][job.GetSubID()]
}

func (j *statuses) getJobStatuses() map[string]map[string]string {
	return j.workJobsStatus
}
