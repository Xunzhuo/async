package async

import (
	"sync"
)

var queueLockers = newLocker()

type asyncJobLocker struct {
	locker   *sync.RWMutex
	lockList map[string]bool
}

func (l *asyncJobLocker) hasLock(job *Job) bool {
	if status, ok := l.lockList[job.jobID]; ok && status {
		return true
	}
	return false
}

func (l *asyncJobLocker) getLockJobs() []string {
	lockJobs := make([]string, 0)
	for lockJob := range l.lockList {
		lockJobs = append(lockJobs, lockJob)
	}
	return lockJobs
}

func newLocker() *asyncJobLocker {
	return &asyncJobLocker{
		locker:   new(sync.RWMutex),
		lockList: make(map[string]bool),
	}
}

func (a *Queue) LockJob(job *Job) {
	queueLockers.locker.Lock()
	defer queueLockers.locker.Unlock()
	queueLockers.lockList[job.jobID] = true
}

func (a *Queue) IsLock(job *Job) bool {
	return queueLockers.hasLock(job)
}

func (a *Queue) UnLockJob(job *Job) {
	queueLockers.locker.Lock()
	defer queueLockers.locker.Unlock()
	queueLockers.lockList[job.jobID] = false
}

func (a *Queue) GetLockJobs() []string {
	return queueLockers.getLockJobs()
}
