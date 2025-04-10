package domain

// TaskExecutor defines the interface for executing tasks against another microservice
type TaskExecutor interface {
	// ExecuteTask executes a task for the given interest
	ExecuteTask(interest *Interest) error
}
