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

	e := concurrency.NewElectionBug(cli, "test-elect")
	if err := e.CampaignBug(context.TODO(), "abc"); err != nil {
		t.Fatal(err)
	}
	e2 := concurrency.NewElectionBug(cli, "test-elect")
	if err := e2.CampaignBug(context.TODO(), "def"); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	if resp := <-e.ObserveBug(ctx); len(resp.Kvs) == 0 || string(resp.Kvs[0].Value) != "def" {
		t.Fatalf("expected value=%q, got response %v", "def", resp)
	}
}
