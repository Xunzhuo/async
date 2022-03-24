package async

// HasJobKey
func (a *JobWorkQueue) HasJobKey(key string) bool {
	return len(a.workJobsStatus[key]) != 0
}

func (a *JobWorkQueue) GetAllJobID() map[string][]string {
	return a.workJobIDHisory
}

func (a *JobWorkQueue) GetJobSubID(id string) []string {
	return a.workJobIDHisory[id]
}

func (a *JobWorkQueue) GetSingleJobSubID(id string) string {
	return a.workJobIDHisory[id][0]
}
