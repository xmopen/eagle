package pkg

import "time"

// IWorker is interface of pool worker
type IWorker interface {
	new(pool *Pool) IWorker
	run()
	finish()
	laseUsedTime() time.Time
	addTaskFunc(func())
	workerLoop()
}

// IWorkerContainer go worker idleContainer
type IWorkerContainer interface {
	len() int
	isEmpty() bool
	addWorker(worker IWorker) error
	tryGetWorker() IWorker
	swapWorkerToRunning(int64, IWorker)
	refresh(duration time.Duration) []IWorker
	reset()
}
