package pkg

import (
	"fmt"
	"runtime"
)

var (
	workerChannelSize = func() int {
		// runtime.GOMAXPROCS(0) return the current number of P in GMP,the number of cpu cores
		if runtime.GOMAXPROCS(0) == 1 {
			return 0
		}
		return 1
	}
)

var (
	// ErrorPoolIsClosed Eagle pool is closed error
	ErrorPoolIsClosed = fmt.Errorf("eagle pool is closed")
	// ErrorPoolWaitingTaskMax Eagle pool waiting task is max
	ErrorPoolWaitingTaskMax = fmt.Errorf("eagle pool waiting is max")
)
