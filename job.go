package async

import (
	"reflect"
	"time"

	"github.com/google/uuid"
)

// Job Item in Async Queue
type Job struct {
	JobKey

	TaskName  string
	StartTime int64
	Status    string

	Handler reflect.Value
	Params  []reflect.Value
}

type JobKey struct {
	Lock         bool
	JobID        string
	NewSubJobID  string
	SubJobIDs    []string
	EnableSubjob bool
}

func (j *Job) SetStatus(status string) {
	j.Status = status
}

func (j *Job) SetTaskName(taskName string) {
	j.TaskName = taskName
}

func (j *Job) GetSubID() string {
	if j.EnableSubjob {
		return j.NewSubJobID
	}
	return NaN
}

func (j *Job) AddSubJob(id string) {
	j.EnableSubjob = true
	j.NewSubJobID = id
	j.SubJobIDs = append(j.SubJobIDs, id)
}

func newJob(jobID string, taskName string, startTime int64,
	handler interface{}, params ...interface{}) Job {

	newJob := Job{
		TaskName:  taskName,
		StartTime: startTime,
		Handler:   reflect.ValueOf(handler),
		Params:    make([]reflect.Value, 0),
	}

	newJob.JobID = jobID

	if len(params) != 0 {
		for _, param := range params {
			newJob.Params = append(newJob.Params, reflect.ValueOf(param))
		}
	}

	return newJob
}

func NewJob(jobID string, handler interface{}, params ...interface{}) Job {
	return newJob(jobID, "", time.Now().Unix(), handler, params...)
}

func NewBlindJob(handler interface{}, params ...interface{}) Job {
	uuid := uuid.New()
	jobID := uuid.String()
	return newJob(jobID, "", time.Now().Unix(), handler, params...)
}
