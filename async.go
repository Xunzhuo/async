package async

// NewJobQueue() create an async events queue
var Engine JobWorkQueue

func init() {
	Engine = NewJobQueue()
}
