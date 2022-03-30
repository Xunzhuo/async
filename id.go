package async

import (
	"crypto/md5"
	"fmt"
)

type key struct {
	jobID        string
	subID        string
	enableSubjob bool
	subjobIDs    []string
}

func (j *Job) EncodejobIDWithMD5() *Job {
	j.jobID = fmt.Sprintf("%x", md5.Sum([]byte(j.jobID)))
	return j
}

func (j *Job) EncodesubIDWithMD5() *Job {
	j.subID = fmt.Sprintf("%x", md5.Sum([]byte(j.subID)))
	return j
}

func (j *Job) GetJobID() string {
	return j.jobID
}

func (j *Job) SetJobID(jobID string) *Job {
	j.jobID = jobID
	return j
}

func (j *Job) GetSubID() string {
	if j.enableSubjob {
		return j.subID
	}
	return keyOfnoSubID
}

func (j *Job) SetSubID(id string) *Job {
	j.enableSubjob = true
	j.subID = id
	j.subjobIDs = append(j.subjobIDs, id)
	return j
}

func (j *Job) GetSubIDs() []string {
	if j.enableSubjob {
		return j.subjobIDs
	}
	return []string{keyOfnoSubID}
}

func (a *Queue) GetJobsID() map[string][]string {
	history := make(map[string][]string)
	jobIDs := a.GetJobStatuses()
	for job, subJobs := range jobIDs {
		for _, subJob := range subJobs {
			history[job] = append(history[job], subJob)
		}
	}
	return history
}
