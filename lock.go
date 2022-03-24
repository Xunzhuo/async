package async

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

func (a *JobWorkQueue) LockJob(job Job) {
	var lock sync.Mutex
	lock.Lock()
	log.Warning("Lock Job with JobID: ", job.JobID)
	a.lockJobIDList[job.JobID] = true
	lock.Unlock()
}

func (a *JobWorkQueue) IsLock(job Job) bool {
	log.Debug("Checking if job is locked with JobID: ", job.JobID)
	if _, ok := a.lockJobIDList[job.JobID]; ok {
		return true
	}
	return false
}

func (a *JobWorkQueue) UnLockJob(job *Job) {
	var lock sync.Mutex
	lock.Lock()
	log.Info("UnLock Job with JobID: ", job.JobID)
	job.Lock = false
	lock.Unlock()
}
