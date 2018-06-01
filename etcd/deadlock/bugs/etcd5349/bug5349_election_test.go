package integration

import (
	"context"
	"testing"

	"github.com/coreos/etcd/clientv3/concurrency"
)

// TestElectionSessionRelock ensures that campaigning twice on the same election
// with the same lock will Proclaim instead of deadlocking.
func TestElectionSessionRecampaignBug5349(t *testing.T) {
	clus := NewClusterV3(t, &ClusterConfig{Size: 1})
	defer clus.Terminate(t)
	cli := clus.RandClient()

	e := concurrency.NewElection(cli, "test-elect")
	if err := e.Campaign(context.TODO(), "abc"); err != nil {
		t.Fatal(err)
	}
	e2 := concurrency.NewElection(cli, "test-elect")
	if err := e2.Campaign(context.TODO(), "def"); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	if resp := <-e.Observe(ctx); len(resp.Kvs) == 0 || string(resp.Kvs[0].Value) != "def" {
		t.Fatalf("expected value=%q, got response %v", "def", resp)
	}
}
