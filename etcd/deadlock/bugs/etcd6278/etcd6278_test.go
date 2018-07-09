package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"golang.org/x/net/context"
)


// TestElectionOnPrefixOfExistingKey checks that a single
// candidate can be elected on a new key that is a prefix
// of an existing key. To wit, check for regression
// of bug #6278. https://github.com/coreos/etcd/issues/6278
//
func TestElectionOnPrefixOfExistingKey(t *testing.T) {
	clus := NewClusterV3(t, &ClusterConfig{Size: 1})
	defer clus.Terminate(t)

	cli := clus.RandClient()
	if _, err := cli.Put(context.TODO(), "testa", "value"); err != nil {
		t.Fatal(err)
	}
	s, serr := concurrency.NewSession(cli)
	if serr != nil {
		t.Fatal(serr)
	}
	e := concurrency.NewElection(s, "test")
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	err := e.Campaign(ctx, "abc")
	cancel()
	if err != nil {
		// after 5 seconds, deadlock results in
		// 'context deadline exceeded' here.
		t.Fatal(err)
	}
}
