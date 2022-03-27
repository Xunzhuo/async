package async

// async.DefaultAsyncQueue() create an default job workQueue
var defaultAsyncQueue *Queue = Q()

func DefaultAsyncQueue() *Queue {
	return defaultAsyncQueue
}

func SetDfaultAsyncQueue(q *Queue) {
	defaultAsyncQueue = q
}
