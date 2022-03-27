package async

type key struct {
	JobID        string
	SubID        string
	SubJobIDs    []string
	EnableSubjob bool
}

func newKey() *key {
	return &key{
		SubJobIDs: make([]string, 0),
	}
}

func (a *Queue) GetAllJobID() map[string][]string {
	history := make(map[string][]string)
	jobIDs := a.GetJobStatuses()
	for job, subJobs := range jobIDs {
		for _, subJob := range subJobs {
			history[job] = append(history[job], subJob)
		}
	}
	return history
}

func (a *Queue) GetJobSubID(id string) []string {
	return a.GetAllJobID()[id]
}

func (a *Queue) GetSingleJobSubID(id string) string {
	return a.GetJobSubID(id)[0]
}
