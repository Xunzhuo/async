package async

// JobWorkQueueOptions some options when creating job workqueue
type JobWorkQueueOptions struct {
	MaxWaitQueueLength int
	MaxWorkQueueLength int
}

type Option func(*JobWorkQueueOptions)

func WithMaxWaitQueueLength(maxWaitQueueLength int) Option {
	return func(opt *JobWorkQueueOptions) {
		opt.MaxWaitQueueLength = maxWaitQueueLength
	}
}

func WithMaxWorkQueueLength(maxWorkQueueLength int) Option {
	return func(opt *JobWorkQueueOptions) {
		opt.MaxWorkQueueLength = maxWorkQueueLength
	}
}
