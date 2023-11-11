package pkg

import (
	"sync"
	"sync/atomic"

	"github.com/xmopen/golib/pkg/xlogging"
)

// poolStatus pool status
type poolStatus int32

// pool status enum
const (
	// poolStatusClosed the pool is closed
	poolStatusClosed poolStatus = iota
	// poolStatusRunning hte pool is running
	poolStatusRunning
)

// default variables
var (
	defaultLocker = &sync.RWMutex{}
	defaultCond   = sync.NewCond(&sync.RWMutex{})
)

// Pool goroutine pool
// pay attention to memory alignment
type Pool struct {
	capacity   int32
	state      int32
	running    int64
	waiting    int64
	trace      bool
	options    *Option // pointer is 8 byte
	workerPool *sync.Pool
	workerCTL  *IWorkerContainer

	cond   *sync.Cond      // cond Lock and UnLock
	xlog   *xlogging.Entry // xlog pointer
	locker sync.Locker     // locker is interface, interface is 16 byte
}

// New a goroutine pool instance
func New(size int, ops ...OptionFun) *Pool {
	pool := &Pool{
		locker: defaultLocker,
		cond:   defaultCond,
		xlog:   xlogging.Tag("eagle.pool"),
	}
	initEaglePoolOption(pool, ops...)
	return pool
}

func (p *Pool) atomicAddRunningSize(delta int64) {
	atomic.AddInt64(&p.running, delta)
}

func (p *Pool) atomicAddWaitingSize(delta int64) {
	atomic.AddInt64(&p.waiting, delta)
}
