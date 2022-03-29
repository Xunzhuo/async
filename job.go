package async

import (
	"errors"
	"reflect"
	"time"

	"github.com/google/uuid"
)

// Job Item in Async Queue
type Job struct {
	key
	taskName  string
	startTime int64

	handler reflect.Value
	params  []reflect.Value
}

func (j *Job) SetTaskName(taskName string) *Job {
	j.taskName = taskName
	return j
}

func (j *Job) GetTaskName() string {
	return j.taskName
}

func (j *Job) SetJobID(jobID string) *Job {
	j.JobID = jobID
	return j
}

func (j *Job) GetJobID() string {
	return j.JobID
}

func (j *Job) GetSubID() string {
	if j.EnableSubjob {
		return j.JobID + "-" + j.SubID
	}
	return j.JobID + "-" + keyOfnoSubID
}

func (j *Job) SetSubID(id string) *Job {
	j.EnableSubjob = true
	j.SubID = id
	j.SubJobIDs = append(j.SubJobIDs, id)
	return j
}

func NewJob(handler interface{}, params ...interface{}) (*Job, error) {
	uuid := uuid.New()
	jobID := uuid.String()
	if reflect.TypeOf(handler).Kind() != reflect.Func {
		return nil, errors.New("Job handler is not func")
	}
	return newJob(jobID, "", time.Now().Unix(), handler, params...), nil
}

func newJob(jobID string, taskName string, startTime int64,
	handler interface{}, params ...interface{}) *Job {

	newJob := Job{
		key:       *newKey(),
		taskName:  taskName,
		startTime: startTime,
		handler:   reflect.ValueOf(handler),
		params:    make([]reflect.Value, 0),
	}
	newJob.JobID = jobID

	if len(params) != 0 {
		for _, param := range params {
			newJob.params = append(newJob.params, reflect.ValueOf(param))
		}
	}

	return &newJob
}
