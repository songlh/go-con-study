package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const numThread = 2
const workCount = 1000000
var wg sync.WaitGroup

var startTime time.Time
var endTime time.Time
var minLatency = time.Duration(0)
var maxLatency = time.Duration(0)
var avgLatency = time.Duration(0)
var totalLatency = time.Duration(0)


type _counter struct {
	lock *sync.Mutex
	count int64
}

var l = &sync.Mutex{}

var counter = _counter{
	lock: l,
	count: 0,
}

func reportSwitchTime() {
	if counter.count > 0 {
		duration := endTime.Sub(startTime)
		totalLatency += duration

		if minLatency > duration {
			minLatency = duration
		}

		if maxLatency < duration {
			maxLatency = duration
		}

		avgLatency = (avgLatency + duration) / 2

		fmt.Printf("Total time is %v, counter is %v, current time is %v," +
			" avg time is %v, min time is %v, max time is %v \n",
			totalLatency, counter.count, duration, avgLatency, minLatency, maxLatency)
	}
}

func task() {
	for counter.count < workCount {
		counter.lock.Lock()
		endTime = time.Now()
		reportSwitchTime()
		counter.count += 1
		startTime = time.Now()
		counter.lock.Unlock()
		runtime.Gosched()
	}
	wg.Done()
}

func main() {
	wg.Add(numThread)
	for i := 0; i < numThread; i++ {
		go task()
	}
	wg.Wait()

}