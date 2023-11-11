package pkg

import "time"

type OptionFun = func(opt *Option)

// Option eagle init option
type Option struct {
	EnableTrace  bool
	EnablePurge  bool          // EnablePurge is set to true,the worker are purged,otherwise,it will not be purged
	ExpiredTime  time.Duration // ExpiredTime the worker will be purged will now - lasttime > this expiredTime
	Nonblocking  bool          // Nonblocking is true will return err when pool size is full
	MaxBlockTask int64         // MaxBlockTask the pool max block task
	PanicHandler PanicHandler  // PanicHandler panic handler
}

// initEaglePoolOption init eagle pool option
func initEaglePoolOption(pool *Pool, ops ...OptionFun) {
	pool.options = &Option{}
	for _, op := range ops {
		op(pool.options)
	}
}
