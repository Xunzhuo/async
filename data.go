package async

func (a *Queue) GetJobsData(job Job) (map[string][]interface{}, bool) {
	if len(a.sharedJobData) == 0 {
		return nil, false
	}

	jobDataList := make(map[string][]interface{})
	if ok := a.HasJob(job); ok {
		for _, subID := range job.GetSubIDs() {
			if _, ok := a.sharedJobData[job.jobID][subID]; ok {
				a.logger.Info("Get Job Data", "jobID", job.jobID, "subID", subID)
				jobDataList[subID] = a.sharedJobData[job.jobID][subID]
			} else {
				a.logger.Info("SubJob is not finished", "jobID", job.jobID, "subID", subID)
				return nil, false
			}
		}
	} else {
		return nil, false
	}

	return jobDataList, true
}

func (a *Queue) GetSubJobData(job Job) ([]interface{}, bool) {
	if jobsData, ok := a.GetJobsData(job); ok {
		return jobsData[job.subID], true
	} else {
		return nil, false
	}
}

func (a *Queue) GetJobData(job Job) ([]interface{}, bool) {
	if len(a.sharedJobData) == 0 {
		return nil, false
	}

	if job.enableSubjob {
		return nil, false
	}

	if _, ok := a.sharedJobData[job.jobID]; ok {
		if _, ok := a.sharedJobData[job.jobID][job.GetSubID()]; ok {
			return a.sharedJobData[job.jobID][job.GetSubID()], true
		}
		return nil, false
	}

	return nil, false
}

// CleanJobData
func (a *Queue) CleanJobData(job Job) {
	delete(a.sharedJobData, job.jobID)
}

// CleanJobDatas
func (a *Queue) CleanJobDatas(jobIDList ...Job) {
	for _, job := range jobIDList {
		a.CleanJobData(job)
	}
}

// CleanHistory
func (a *Queue) CleanHistory() {
	a.sharedJobData = make(map[string]map[string][]interface{})
}
