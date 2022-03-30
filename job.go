package async

import (
	"crypto/md5"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
)

// Job Item in Async Queue
type Job struct {
	*key
	startTime int64

	handler reflect.Value
	params  []reflect.Value
}

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
	j.subID = fmt.Sprintf("%x", md5.Sum([]byte(j.jobID)))
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

func NewJob(handler interface{}, params ...interface{}) *Job {
	uuid := uuid.New()
	jobID := uuid.String()
	return newJob(jobID, time.Now().Unix(), handler, params...)
}

func newJob(jobID string, startTime int64,
	handler interface{}, params ...interface{}) *Job {

	newJob := Job{
		key:       newKey(),
		startTime: startTime,
		handler:   reflect.ValueOf(handler),
		params:    make([]reflect.Value, 0),
	}
	newJob.jobID = jobID

	if len(params) != 0 {
		for _, param := range params {
			newJob.params = append(newJob.params, reflect.ValueOf(param))
		}
	}

	return &newJob
}

func newKey() *key {
	return &key{
		subjobIDs: make([]string, 0),
	}
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
