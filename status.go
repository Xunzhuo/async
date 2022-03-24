package async

// Stop
func (a *JobWorkQueue) Stop() {
	close(a.workJobsQueue)
}
