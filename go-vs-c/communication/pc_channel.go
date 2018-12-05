package main

import (
	"fmt"
	"time"
	"sync"
)

const np = 32  // 1-64
const WorkCount = 1000000

var wg sync.WaitGroup

func consumer(cname string, ch chan int) {
	for i := 0; i < WorkCount; i++ {
		<-ch
	}
	wg.Done()
}

func producer(pname string, ch chan int) {
	for i := 0; i < WorkCount; i++ {
		ch <- i
	}
	wg.Done()
}

func main() {
	t1 := time.Now()
	data := make(chan int, 64)
	wg.Add(2 * np)
	for i := 0; i < np; i++ {
		go producer("p", data)
	}
	for i := 0; i < np; i++ {
		go consumer("c", data)
	}
	wg.Wait()
	close(data)
	elapsed := time.Since(t1)
	fmt.Printf("The number of goroutine is %v, communication time is %v ns\n", np * 2, elapsed.Nanoseconds())
}
