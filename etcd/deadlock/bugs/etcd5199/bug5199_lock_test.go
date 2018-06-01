package integration

import (
	"context"
	"testing"

	"github.com/coreos/etcd/clientv3/concurrency"
	"fmt"
)

func TestMutexSessionRelockBug5199(t *testing.T) {
	fmt.Println("starting to test.....")
	clus := NewClusterV3(t, &ClusterConfig{Size: 3})
	defer clus.Terminate(t)
	cli := clus.RandClient()
	m := concurrency.NewMutex(cli, "test-mutex")
	if err := m.Lock(context.TODO()); err != nil {
		t.Fatal(err)
	}
	m2 := concurrency.NewMutex(cli, "test-mutex")
	if err := m2.Lock(context.TODO()); err != nil {
		t.Fatal(err)
	}
}