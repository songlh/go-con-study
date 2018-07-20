package grpc

import (
	"testing"
	"time"
)

func TestClientLBWatcherWithClosedBalancer(t *testing.T) {
	b := newBlockingBalancer()
	cc := &ClientConn{dopts: dialOptions{balancer: b}}

	doneChan := make(chan struct{})
	go cc.lbWatcher(doneChan)
	// Balancer closes before any successful connections.
	b.Close()

	select {
	case <-doneChan:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("lbWatcher with closed balancer didn't close doneChan after 100ms")
	}
}
