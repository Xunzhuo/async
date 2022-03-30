package async

import (
	"github.com/go-logr/logr"
	"k8s.io/klog"
	"k8s.io/klog/v2/klogr"
)

var defaultOptions = newDefaultOptions()
var defaultLogger = klogr.New()

func init() {
	klog.InitFlags(nil)
}

// JobWorkQueueOptions some options when creating job workqueue
type options struct {
	maxWaitQueueLength int
	maxWorkQueueLength int
	logger             logr.Logger
}

func newDefaultOptions() *options {
	logger := defaultLogger.WithName("Async")
	return &options{
		maxWaitQueueLength: DefaultWaitQueueCapacity,
		maxWorkQueueLength: DefaultWaitQueueCapacity,
		logger:             logger,
	}
}

func (o *options) setLogger(logger logr.Logger) {
	o.logger = logger
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

func (a *Queue) SetLogger(logger logr.Logger) *Queue {
	defaultOptions.setLogger(logger)
	return a
}
