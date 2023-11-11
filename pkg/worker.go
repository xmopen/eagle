package pkg

import (
	"time"

	"github.com/xmopen/golib/pkg/xlogging"
)

// Worker pool worker
type Worker struct {
	pool         *Pool
	taskChannel  chan func()
	xlog         *xlogging.Entry
	close        chan struct{}
	lastUsedTime time.Time
}

func (w Worker) new(pool *Pool) IWorker {
	return &Worker{
		pool:         pool,
		taskChannel:  make(chan func(), workerChannelSize()), // 这里是阻塞还是非阻塞呢？
		close:        make(chan struct{}),
		lastUsedTime: time.Now(),
	}
}

func (w Worker) run() {
	w.pool.atomicAddRunningSize(1)
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
			// TODO: 执行完之后还需要将当前Worker放回到idleContainer中
			task()
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
