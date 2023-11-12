package pkg

import "github.com/xmopen/golib/pkg/xlogging"

type PanicHandler = func(err any)

// panicHandler TODO: process panic info
func panicHandler(xlog *xlogging.Entry, handler PanicHandler) {
	if err := recover(); err != nil {
		// 日志用什么日志呢/
		if handler != nil {
			handler(err)
		}
	}
}
