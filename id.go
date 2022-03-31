package async

import (
	"crypto/md5"
	"fmt"
)

type key struct {
	jobID        string
	subID        string
	enableSubjob bool
	subjobIDs    map[string]string
}

func (j *Job) EncodejobIDWithMD5() {
	j.jobID = fmt.Sprintf("%x", md5.Sum([]byte(j.jobID)))
}

func (j *Job) EncodesubIDWithMD5() {
	j.subID = fmt.Sprintf("%x", md5.Sum([]byte(j.subID)))
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
	return j.jobID + "/" + keyOfnoSubID
}

func (j *Job) SetSubID(id string) *Job {
	j.enableSubjob = true
	j.subID = id
	queueLockers.locker.Lock()
	j.subjobIDs[id] = StatusPending
	queueLockers.locker.Unlock()
	return j
}

func (j *Job) GetSubIDs() []string {
	if j.enableSubjob {
		subIDs := make([]string, 0)
		queueLockers.locker.RLock()
		defer queueLockers.locker.RUnlock()
		for subID := range j.subjobIDs {
			subIDs = append(subIDs, subID)
		}
		return subIDs
	}
	return []string{j.GetSubID()}
}

func (a *Queue) GetJobsID() map[string][]string {
	history := make(map[string][]string)
	jobIDs := a.GetJobStatuses()
	for job, subJobs := range jobIDs {
		for subJob := range subJobs.subjobIDs {
			history[job] = append(history[job], subJob)
		}
	}
	return history
}

func (a *Queue) GetJobByID(id string) *Job {
	jobIDs := a.GetJobStatuses()
	if job, ok := jobIDs[id]; ok {
		return job
	}
	return nil
}
