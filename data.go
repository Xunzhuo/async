package async

import log "github.com/sirupsen/logrus"

func (a *JobWorkQueue) GetJobsData(jobID string) (map[string][]interface{}, bool) {
	if len(a.SharedJobData) == 0 {
		return nil, false
	}

	jobDataList := make(map[string][]interface{})
	if ok := a.HasJobKey(jobID); ok {
		for _, subID := range a.GetJobSubID(jobID) {
			log.Info("Get Job Data with JobID: ", jobID, " subID: ", subID)
			if _, ok := a.SharedJobData[jobID][subID]; ok {
				jobDataList[subID] = a.SharedJobData[jobID][subID]
			} else {
				log.Info("Can not get Job Data with JobID: ", jobID, " subID: ", subID)
			}
		}
	} else {
		return nil, false
	}

	return jobDataList, true
}

func (a *JobWorkQueue) GetSubJobData(jobID, subJobID string) ([]interface{}, bool) {
	if jobsData, ok := a.GetJobsData(jobID); ok {
		return jobsData[subJobID], true
	} else {
		return nil, false
	}
}

func (a *JobWorkQueue) GetJobData(jobID string) ([]interface{}, bool) {
	if len(a.SharedJobData) == 0 {
		return nil, false
	}

	if _, ok := a.workJobIDHisory[jobID]; ok {
		if _, ok := a.SharedJobData[jobID][a.GetSingleJobSubID(jobID)]; ok {
			return a.SharedJobData[jobID][a.GetSingleJobSubID(jobID)], true
		}
		return nil, false
	}

	return nil, false
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
