package pkg

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

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

const (
	UnLimitPoolGoroutineSize = -1
)

// default variables
var (
	defaultLocker = &sync.RWMutex{}
	defaultCond   = sync.NewCond(&sync.RWMutex{})
)

// Pool goroutine pool
// pay attention to memory alignment
type Pool struct {
	workID     int64
	capacity   int
	state      poolStatus
	waiting    int64
	trace      bool
	options    *Option // pointer is 8 byte
	workerPool *sync.Pool
	workerCTL  IWorkerContainer // workerCTL container contain with running number and idle number

	cond   *sync.Cond      // cond Lock and UnLock
	xlog   *xlogging.Entry // xlog pointer
	locker sync.Locker     // locker is interface, interface is 16 byte
	cancel context.CancelFunc
}

// New a goroutine pool instance
// If size is -1 means thant the Pool does not limit the goroutines
func New(size int, ops ...OptionFun) *Pool {
	pool := &Pool{
		capacity: size,
		locker:   defaultLocker,
		cond:     defaultCond,
		state:    poolStatusRunning,
		xlog:     xlogging.Tag("eagle.pool"),
	}
	initEaglePoolOption(pool, ops...)
	pool.workerPool = &sync.Pool{
		New: func() any {
			return newWorker(pool)
		},
	}
	pool.workerCTL = newWorkerContainerMap(pool)
	ctx, cancel := context.WithCancel(context.Background())
	pool.cancel = cancel
	go pool.ticket(ctx)
	return pool
}

func (p *Pool) atomicAddWaitingSize(delta int64) {
	atomic.AddInt64(&p.waiting, delta)
}

// Submit a task func to pool,return err when the pool is full and is blocking
func (p *Pool) Submit(fn func()) error {
	if p.IsClose() {
		return ErrorPoolIsClosed
	}
	// blocking
	worker, err := p.retrieveWorker()
	if err != nil {
		return err
	}
	worker.run()
	worker.addTaskFunc(fn)
	return nil
}

// Close the running pool
func (p *Pool) Close() {
	state := int32(p.state)
	atomic.StoreInt32(&state, int32(poolStatusClosed))
	p.cancel()
}

// IsClose return the pool is closed
func (p *Pool) IsClose() bool {
	state := int32(p.state)
	return atomic.LoadInt32(&state) == int32(poolStatusClosed)
}

// retrieveWorker retrieve worker from container
// If no worker exists in the containers,a new Worker is retrieved from the Pool.workerPoll
func (p *Pool) retrieveWorker() (IWorker, error) {
	for {
		if p.options.EnableTrace {
			p.xlog.Infof("retreve worker")
		}
		// First try get work dones not lock
		if worker := p.workerCTL.tryGetWorker(); worker != nil {
			return worker, nil
		}
		p.locker.Lock()
		if p.workerCTL.idle() > 0 {
			if worker := p.workerCTL.tryGetIdleWorker(); worker != nil {
				p.workerCTL.swapWorkerToRunning(worker)
				p.locker.Unlock()
				return worker, nil
			}
		}
		// verify pool capacity,read varlibe should to lock
		if p.capacity > p.WorkerSize() || p.capacity == UnLimitPoolGoroutineSize {
			p.locker.Unlock()
			worker := p.workerPool.Get().(IWorker)
			// 这里需要加到container中
			if err := p.workerCTL.addWorker(worker); err != nil {
				return nil, err
			}
			return worker, nil
		}
		// verify pool is Nonblocking or MaxWaiting
		// 11的时候 这里会阻塞.
		if p.options.Nonblocking || (p.options.MaxBlockTask != 0 && p.waiting >= p.options.MaxBlockTask) {
			return nil, ErrorPoolWaitingTaskMax
		}
		p.locker.Unlock()

		// blocking,wait the pool.container release IWorker
		// 这里不释放.
		p.wait()
		if p.IsClose() {
			return nil, ErrorPoolIsClosed
		}
	}
}

// wait blocking wait until the Pool Container releases the IWorker
func (p *Pool) wait() {
	atomic.AddInt64(&p.waiting, 1)
	if p.options.EnableTrace {
		p.xlog.Warnf("and wait,now waiting:[%+v]", atomic.LoadInt64(&p.waiting))
	}
	p.cond.L.Lock()
	p.cond.Wait()
	p.cond.L.Unlock()
	atomic.AddInt64(&p.waiting, -1)
	if p.options.EnableTrace {
		p.xlog.Warnf("wait done,now waiting:[%+v]", atomic.LoadInt64(&p.waiting))
	}
}

// WorkerSize return the Pool Container all worker size
func (p *Pool) WorkerSize() int {
	return p.workerCTL.len()
}

// recycleWorker put worker to container,recycle goroutine
func (p *Pool) recycleWorker(worker IWorker) {
	if p.IsClose() {
		p.cond.Broadcast()
		return
	}
	worker.updateLastUsedTime(time.Now())
	p.workerCTL.swapWorkerToIdle(worker)
	worker.finish()
}

func (p *Pool) ticket(ctx context.Context) {
	defer panicHandler(p.xlog, p.options.PanicHandler)
	if !p.options.EnableTicket {
		return
	}
	// 1、running worker
	// 2、idle worker
	// 3、waiting goroutine
	ticket := time.NewTicker(p.options.TicketDuration)
	for {
		select {
		case <-ctx.Done():
			if p.options.EnableTrace {
				p.xlog.Infof("ticket closed success")
			}
			return
		case <-ticket.C:
			if p.options.EnableTrace {
				p.ticketInfo()
			}
			// clear worker
		}
	}
}

func (p *Pool) ticketInfo() {
	p.xlog.Infof("running work:[%d] idle work:[%+v] wait:[%+v]", p.workerCTL.running(), p.workerCTL.idle(),
		atomic.LoadInt64(&p.waiting))
}
