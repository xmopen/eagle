package pkg

import (
	"sync"
	"time"

	"github.com/xmopen/golib/pkg/xlogging"
)

// WorkerMap is a worker container implements with Map
type WorkerMap struct {
	pool             *Pool
	workID           int64
	trace            bool
	idleSize         int
	runningSize      int
	locker           sync.Locker
	xlog             *xlogging.Entry
	idleContainer    *sync.Map // idleContainer is idle worker container
	runningContainer *sync.Map // runningContainer is  running worker container
}

func newWorkerContainerMap(pool *Pool) IWorkerContainer {
	return &WorkerMap{
		pool:             pool,
		trace:            pool.options.EnableTrace,
		locker:           &sync.Mutex{},
		xlog:             pool.xlog,
		idleContainer:    &sync.Map{},
		runningContainer: &sync.Map{},
	}
}

// len return the container all worker number
func (w *WorkerMap) len() int {
	return w.idleSize + w.runningSize
}

func (w *WorkerMap) running() int {
	return w.runningSize
}

func (w *WorkerMap) idle() int {
	return w.idleSize
}

func (w *WorkerMap) isEmpty() bool {
	return w.len() == 0
}

// addWorker First get Worker from worker pool, add to worker container
func (w *WorkerMap) addWorker(worker IWorker) error {
	w.locker.Lock()
	defer w.locker.Unlock()
	if w.trace {
		w.xlog.Infof("worker container[WorkerMap] add worker:[%d]", worker.workerID())
	}
	w.runningSize++
	w.runningContainer.Store(worker.workerID(), worker)
	return nil
}

// tryGetWorker try get a worker from thw container
// However it is possible to return nil when the worker does not exist in the container
func (w *WorkerMap) tryGetWorker() IWorker {
	if w.isEmpty() {
		return nil
	}
	var worker IWorker
	w.idleContainer.Range(func(key, value any) bool {
		if key == nil || value == nil {
			return true
		}
		workerInstance, ok := value.(IWorker)
		if !ok {
			return true
		}
		worker = workerInstance
		return false
	})
	// TODO: why to add check nil?
	if worker == nil {
		return nil
	}
	w.swapWorkerToRunning(worker)
	return worker
}

func (w *WorkerMap) tryGetIdleWorker() IWorker {
	if w.idle() == 0 {
		return nil
	}
	var worker IWorker
	w.idleContainer.Range(func(key, value any) bool {
		if key == nil || value == nil {
			return true
		}
		obj, ok := value.(IWorker)
		if !ok {
			return true
		}
		worker = obj
		return false
	})
	return worker
}

func (w *WorkerMap) swapWorkerToRunning(worker IWorker) {
	w.idleContainer.Delete(worker.workerID())
	w.runningContainer.Store(worker.workerID(), worker)
	w.locker.Lock()
	defer w.locker.Unlock()
	w.idleSize--
	w.runningSize++
}

// swapWorkerToIdle swap the running worker to idle container
// thinking pool is close or
func (w *WorkerMap) swapWorkerToIdle(worker IWorker) {
	w.runningContainer.Delete(worker.workerID())
	w.idleContainer.Store(worker.workerID(), worker)
	if w.pool.options.EnableTrace {
		w.xlog.Infof("worker container signal")
	}
	w.pool.cond.Signal()
	w.locker.Lock()
	defer w.locker.Unlock()
	w.idleSize++
	w.runningSize--
}

// refresh clear expired worker with timeout
func (w *WorkerMap) refresh(timeout time.Duration) []IWorker {
	panic("implement me")
}

func (w *WorkerMap) reset() {
	panic("implement me")
}
