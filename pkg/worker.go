package pkg

import (
	"sync/atomic"
	"time"

	"github.com/xmopen/golib/pkg/xlogging"
)

// Worker pool worker
type Worker struct {
	id           int64
	pool         *Pool
	taskChannel  chan func()
	xlog         *xlogging.Entry
	close        chan struct{}
	lastUsedTime time.Time
}

func newWorker(pool *Pool) IWorker {
	return &Worker{
		id:           atomic.AddInt64(&pool.workID, 1),
		pool:         pool,
		taskChannel:  make(chan func(), workerChannelSize()), // 这里是阻塞还是非阻塞呢？
		close:        make(chan struct{}),
		lastUsedTime: time.Now(),
	}
}

func (w Worker) run() {
	go func() {
		// release resource
		defer w.pool.cond.Signal()
		var handler PanicHandler
		if w.pool.options.PanicHandler != nil {
			handler = w.pool.options.PanicHandler
		}
		defer panicHandler(w.xlog, handler)
		w.workerLoop()
	}()
}

func (w Worker) workerLoop() {
	for {
		select {
		case <-w.close:
			return
		case task := <-w.taskChannel:
			task()
			// 将Worker回访到container中,重复利用
			w.pool.recycleWorker(w)
		}
	}
}

// finish worker end
// ants 中是通过向 taskChannel 传送一个nil func来实现close实际上这个不太合理
func (w Worker) finish() {
	close(w.close)
}

func (w Worker) laseUsedTime() time.Time {
	return w.lastUsedTime
}

func (w Worker) addTaskFunc(task func()) {
	w.taskChannel <- task
}

func (w Worker) updateLastUsedTime(lastTime time.Time) {
	w.lastUsedTime = lastTime
}

func (w Worker) workerID() int64 {
	return w.id
}
