package implementations

import (
	"sync"
	"time"

	"github.com/smnzlnsk/routing-manager/internal/domain"
	"go.uber.org/zap"
)

// TaskSchedulerObserver schedules regular tasks for interests
type TaskSchedulerObserver struct {
	*BaseObserver
	taskExecutor domain.TaskExecutor
	schedulers   map[string]*time.Ticker
	done         map[string]chan bool
	interval     time.Duration
	mutex        sync.Mutex
}

// NewTaskSchedulerObserver creates a new TaskSchedulerObserver
func NewTaskSchedulerObserver(
	logger *zap.Logger,
	taskExecutor domain.TaskExecutor,
	interval time.Duration,
) *TaskSchedulerObserver {
	if interval <= 0 {
		interval = 1 * time.Minute // Default interval
	}

	return &TaskSchedulerObserver{
		BaseObserver: NewBaseObserver("TaskSchedulerObserver", logger),
		taskExecutor: taskExecutor,
		schedulers:   make(map[string]*time.Ticker),
		done:         make(map[string]chan bool),
		interval:     interval,
	}
}

// Update handles interest events by starting or stopping the scheduled tasks
func (o *TaskSchedulerObserver) Update(event domain.InterestEvent) {
	interest := event.Interest
	appName := interest.AppName

	switch event.Type {
	case domain.InterestCreated:
		o.startTaskScheduler(interest)

	case domain.InterestUpdated:
		// If we have a scheduler, stop it and start a new one with the updated interest
		if o.hasScheduler(appName) {
			o.stopTaskScheduler(appName)
			o.startTaskScheduler(interest)
		} else {
			// If no scheduler exists, start a new one
			o.startTaskScheduler(interest)
		}

	case domain.InterestDeleted:
		// Stop the scheduler for deleted interests
		o.stopTaskScheduler(appName)
	}
}

// startTaskScheduler starts a scheduler for the given interest
func (o *TaskSchedulerObserver) startTaskScheduler(interest *domain.Interest) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	appName := interest.AppName

	// Make a copy of the interest to prevent issues with concurrent access
	interestCopy := &domain.Interest{
		AppName:   interest.AppName,
		ServiceIp: interest.ServiceIp,
		CreatedAt: interest.CreatedAt,
		UpdatedAt: interest.UpdatedAt,
	}

	// Create a ticker for the scheduler
	ticker := time.NewTicker(o.interval)
	done := make(chan bool)

	o.schedulers[appName] = ticker
	o.done[appName] = done

	o.logger.Info("Started task scheduler",
		zap.String("appName", appName),
		zap.Duration("interval", o.interval))

	// Start the scheduler in a goroutine
	go func() {
		// Execute immediately on start
		/*if err := o.taskExecutor.ExecuteTask(interestCopy); err != nil {
			o.logger.Error("Failed to execute initial task",
				zap.String("appName", appName),
				zap.Error(err))
		}*/

		// Then continue with the ticker
		for {
			select {
			case <-ticker.C:
				if err := o.taskExecutor.ExecuteTask(interestCopy); err != nil {
					o.logger.Error("Failed to execute scheduled task",
						zap.String("appName", appName),
						zap.Error(err))
				}
			case <-done:
				return
			}
		}
	}()
}

// stopTaskScheduler stops the scheduler for the given app name
func (o *TaskSchedulerObserver) stopTaskScheduler(appName string) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	if ticker, ok := o.schedulers[appName]; ok {
		ticker.Stop()
		close(o.done[appName])

		delete(o.schedulers, appName)
		delete(o.done, appName)

		o.logger.Info("Stopped task scheduler", zap.String("appName", appName))
	}
}

// hasScheduler checks if a scheduler exists for the given app name
func (o *TaskSchedulerObserver) hasScheduler(appName string) bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	_, exists := o.schedulers[appName]
	return exists
}

// Shutdown stops all schedulers
func (o *TaskSchedulerObserver) Shutdown() {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	for appName, ticker := range o.schedulers {
		ticker.Stop()
		close(o.done[appName])
		o.logger.Info("Stopped task scheduler during shutdown", zap.String("appName", appName))
	}

	o.schedulers = make(map[string]*time.Ticker)
	o.done = make(map[string]chan bool)

	o.logger.Info("All task schedulers stopped")
}
