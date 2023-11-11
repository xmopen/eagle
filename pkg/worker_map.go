package pkg

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/xmopen/golib/pkg/xlogging"
)

// WorkerMap is a worker container implements with Map
type WorkerMap struct {
	workID           int64
	trace            bool
	idleSize         int
	runningSize      int
	locker           sync.Locker
	xlog             *xlogging.Entry
	idleContainer    sync.Map // idleContainer is idle worker container
	runningContainer sync.Map // runningContainer is  running worker container
}

func (w *WorkerMap) len() int {
	return w.idleSize
}

func (w *WorkerMap) isEmpty() bool {
	return w.len() == 0
}

func (w *WorkerMap) addWorker(worker IWorker) error {
	w.locker.Lock()
	defer w.locker.Unlock()
	workID := w.generateWorkerID()
	if w.trace {
		w.xlog.Infof("worker container[WorkerMap] add worker:[%d]", workID)
	}
	w.idleSize++
	w.idleContainer.Store(workID, worker)
	return nil
}

func (w *WorkerMap) generateWorkerID() int64 {
	return atomic.AddInt64(&w.workID, 1)
}

func (w *WorkerMap) tryGetWorker() IWorker {
	if w.isEmpty() {
		return nil
	}
	var (
		workerID int64
		worker   IWorker
	)
	w.locker.Lock()
	defer w.locker.Unlock()
	w.idleContainer.Range(func(key, value any) bool {
		if worker == nil {
			return true
		}
		workerInstance, ok := value.(IWorker)
		if !ok {
			return true
		}
		wid, _ := key.(int64)
		workerID = wid
		worker = workerInstance
		return false
	})
	w.swapWorkerToRunning(workerID, worker)
	return worker
}

func (w *WorkerMap) swapWorkerToRunning(workerID int64, worker IWorker) {
	w.idleContainer.Delete(workerID)
	w.idleSize--
	w.runningContainer.Store(workerID, worker)
	w.runningSize++
}

// refresh clear expired worker with timeout
func (w *WorkerMap) refresh(timeout time.Duration) []IWorker {
	panic("implement me")
}

func (w *WorkerMap) reset() {
	panic("implement me")
}
