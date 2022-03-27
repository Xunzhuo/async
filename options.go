package async

var defaultOptions = newDefaultOptions()

// JobWorkQueueOptions some options when creating job workqueue
type options struct {
	maxWaitQueueLength int
	maxWorkQueueLength int
}

func newDefaultOptions() *options {
	return &options{
		maxWaitQueueLength: DefaultWaitQueueCapacity,
		maxWorkQueueLength: DefaultWaitQueueCapacity,
	}
}

func (o *options) setMaxWaitQueueLength(length int) {
	o.maxWaitQueueLength = length
}

func (o *options) setMaxWorkQueueLength(length int) {
	o.maxWorkQueueLength = length
}

func (a *Queue) SetMaxWaitQueueLength(length int) *Queue {
	defaultOptions.setMaxWaitQueueLength(length)
	return a
}

func (a *Queue) SetMaxWorkQueueLength(length int) *Queue {
	defaultOptions.setMaxWorkQueueLength(length)
	return a
}
