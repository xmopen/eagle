package pkg

import "time"

const defaultOptionMaxBlock = 2 << 10

type OptionFun = func(opt *Option)

// Option eagle init option
type Option struct {
	EnableTrace    bool          // EnableTrace If ture,log info message of the eagle pool
	EnablePurge    bool          // EnablePurge is set to true,the worker are purged,otherwise,it will not be purged
	ExpiredTime    time.Duration // ExpiredTime the worker will be purged will now - lasttime > this expiredTime
	Nonblocking    bool          // Nonblocking is true will return err when pool size is full
	MaxBlockTask   int64         // MaxBlockTask the pool max block task
	PanicHandler   PanicHandler  // PanicHandler customer panic handler
	EnableTicket   bool          // EnableTicket If ture, start goroutine to process eagle pool with TicketDuration interval
	TicketDuration time.Duration // TicketDuration eagle ticket interval,This TicketDuration takes effect Only EnableTicket is true
}

// initEaglePoolOption init eagle pool option
func initEaglePoolOption(pool *Pool, ops ...OptionFun) {
	pool.options = &Option{}
	for _, op := range ops {
		op(pool.options)
	}
	if pool.options.MaxBlockTask < 0 && !pool.options.Nonblocking {
		pool.options.MaxBlockTask = defaultOptionMaxBlock
	}
}
