package async

import (
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

func NewJob(handler interface{}, params ...interface{}) *Job {
	uuid := uuid.New()
	jobID := uuid.String()
	return newJob(jobID, time.Now().Unix(), handler, params...)
}

func newJob(jobID string, startTime int64,
	handler interface{}, params ...interface{}) *Job {

	newJob := &Job{
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

	return newJob
}

func newKey() *key {
	return &key{
		subjobIDs: make(map[string]string),
	}
}
