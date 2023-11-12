package pkg

import (
	"fmt"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	pool := New(10, func(opt *Option) {
		opt.EnableTicket = true
		opt.TicketDuration = 5 * time.Second
		opt.EnableTrace = true
	})
	for i := 0; i < 20; i++ {
		fmt.Println("submit")
		temp := i
		err := pool.Submit(func() {
			fmt.Printf("submit fn res:%d\n", temp)
		})
		if err != nil {
			panic(err)
		}
	}
	time.Sleep(10 * time.Minute)
	pool.Close()
	time.Sleep(3 * time.Second)
	// 0 1 2 3 4 5  6 7 9 10 12 14 15 16 17 18 19
	// 8 11  13 为什么呢？
}
