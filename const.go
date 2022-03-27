package async

// values of some keys
const (
	keyOfJobID   = "JobID"
	keyOfSubID   = "SubID"
	keyOfjobData = "JobData"
	keyOfnoSubID = "noSubID"
)

// values of status
const (
	// StatusPending ...
	StatusPending = "pending"
	// StatusRunning ...
	StatusRunning = "running"
	// StatusFailure ...
	StatusFailure = "failure"
	// StatusSuccess ...
	StatusSuccess = "success"
	// StatusCancelled ...
	StatusCancelled = "cancelled"
	// StatusSkipped ...
	StatusSkipped = "skipped"
	// StatusUnknown ...
	StatusUnknown = "Unknown"
)

// values of options
const (
	// DefaultWorkQueueCapacity ...
	DefaultWorkQueueCapacity = 100
	// DefaultWaitQueueCapacity ...
	DefaultWaitQueueCapacity = 100
)
