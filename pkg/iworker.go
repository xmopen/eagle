package pkg

import "time"

// TODO: rename ipool

// IWorker is interface of pool worker
type IWorker interface {
	run()
	finish()
	laseUsedTime() time.Time
	addTaskFunc(func())
	workerLoop()
	updateLastUsedTime(lastTime time.Time)
	workerID() int64
}

// IWorkerContainer go worker idleContainer
type IWorkerContainer interface {
	len() int
	running() int
	idle() int
	isEmpty() bool
	addWorker(worker IWorker) error
	tryGetWorker() IWorker
	tryGetIdleWorker() IWorker
	swapWorkerToRunning(IWorker)
	swapWorkerToIdle(worker IWorker)
	refresh(duration time.Duration) []IWorker
	reset()
}
