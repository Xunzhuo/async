package async

import log "github.com/sirupsen/logrus"

func (a *Queue) GetJobsData(job Job) (map[string][]interface{}, bool) {
	if len(a.SharedJobData) == 0 {
		return nil, false
	}

	jobDataList := make(map[string][]interface{})
	if ok := a.HasJob(job); ok {
		for _, subID := range a.GetJobSubID(job.JobID) {
			log.Info("Get Job Data with job.JobID: ", job.JobID, " subID: ", subID)
			if _, ok := a.SharedJobData[job.JobID][subID]; ok {
				jobDataList[subID] = append(jobDataList[subID], a.SharedJobData[job.JobID][subID])
			} else {
				log.Info("Can not get Job Data with job.JobID: ", job.JobID, " subID: ", subID)
			}
		}
	} else {
		return nil, false
	}

	return jobDataList, true
}

func (a *Queue) GetSubJobData(job Job) ([]interface{}, bool) {
	if jobsData, ok := a.GetJobsData(job); ok {
		return jobsData[job.SubID], true
	} else {
		return nil, false
	}
}

func (a *Queue) GetJobData(job Job) ([]interface{}, bool) {
	if len(a.SharedJobData) == 0 {
		return nil, false
	}

	if _, ok := a.SharedJobData[job.JobID]; ok {
		if _, ok := a.SharedJobData[job.JobID][a.GetSingleJobSubID(job.JobID)]; ok {
			return a.SharedJobData[job.JobID][a.GetSingleJobSubID(job.JobID)], true
		}
		return nil, false
	}

	return nil, false
}

// CleanJobData
func (a *Queue) CleanJobData(job Job) {
	delete(a.SharedJobData, job.JobID)
}

// CleanJobDatas
func (a *Queue) CleanJobDatas(jobIDList ...Job) {
	for _, job := range jobIDList {
		a.CleanJobData(job)
	}
}
