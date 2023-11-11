package pkg

import "runtime"

var (
	workerChannelSize = func() int {
		// runtime.GOMAXPROCS(0) return the current number of p in GMP,the number of cpu cores
		if runtime.GOMAXPROCS(0) == 1 {
			return 0
		}
		return 1
	}
)
