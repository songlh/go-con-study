package main

import (
	"fmt"
	"sync"
)

// if NUM is greater than 1, it would be a deadlock
// if start a new goroutine run moby25384(), it would be a partial deadlock
// if moby25384 and main are in same goroutine, it would be a global deadlock

var NUM = 2

func moby25384() {
	var group sync.WaitGroup

	group.Add(NUM) // Must wait for 10 calls to 'done' before moving on
	for i := 0; i < NUM; i++ {
		go func(i int) {
			fmt.Println("Start gorouting ", i)
			defer group.Done()
			fmt.Println("End gorouting ", i)
		}(i)
		group.Wait()
	}
}

func main() {
	fmt.Println("start")
	go moby25384()
	fmt.Println("end")
}
