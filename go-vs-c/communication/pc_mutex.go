package main

import (
	"sync"
	"time"
	"fmt"
	"runtime"
)

var wg sync.WaitGroup

const np = 32  // 1-64
const WorkCount = 1000000

type buffer_ struct {
	buf         [64]int
	lock        *sync.Mutex
	len         int
	can_produce *sync.Cond
	can_consume *sync.Cond
}

var l = &sync.Mutex{}
var buffer = buffer_{
	len: 0,
	lock:        l,
	can_produce: sync.NewCond(l),
	can_consume: sync.NewCond(l),
}

func producer() {
	for i := 0; i < WorkCount; i++ {
		buffer.lock.Lock()
		for buffer.len == 64 {
			buffer.can_produce.Wait()
		}
		buffer.buf[buffer.len] = i
		buffer.len += 1
		buffer.can_consume.Signal()
		buffer.lock.Unlock()
	}
	wg.Done()
}

func consumer() {
	for i := 0; i < WorkCount; i++ {
		buffer.lock.Lock()
		for buffer.len == 0 {
			buffer.can_consume.Wait()
		}
		buffer.len -= 1
		buffer.can_produce.Signal()
		buffer.lock.Unlock()
	}
	wg.Done()
}

func main() {
	runtime.GOMAXPROCS(24)
	t1 := time.Now()
	wg.Add(np * 2)
	for i := 0; i < np; i++ {
		go producer()
	}
	for i := 0; i < np; i++ {
		go consumer()
	}
	wg.Wait()
	elapsed := time.Since(t1)
	fmt.Printf("The number of goroutine is %v, communication time is %v ns \n", np * 2, elapsed.Nanoseconds())
}
