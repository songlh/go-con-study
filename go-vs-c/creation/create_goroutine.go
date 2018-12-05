package main

import (
	"fmt"
	"runtime"
	"time"
)

const GoroutineNum = 10000

func consumer() {
	runtime.Gosched()
}

func main() {
	//runtime.GOMAXPROCS(4)
	//j := 1
	t1 := time.Now()
	for i := 0; i < GoroutineNum; i++ {
		go consumer()
		//j += 1
	}
	elapsed := time.Since(t1)

	fmt.Printf("The number of C thread is %v, creation time is %v ns", GoroutineNum, elapsed.Nanoseconds())
}
