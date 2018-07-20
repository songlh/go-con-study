func TestKVGetErrConnClosed(t *testing.T) {
	defer testutil.AfterTest(t)

	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer clus.Terminate(t)

	cli := clus.Client(0)
	kv := clientv3.NewKV(cli)

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		_, err := kv.Get(context.TODO(), "foo")
		if err != nil && err != rpctypes.ErrConnClosed {
			t.Fatalf("expected %v, got %v", rpctypes.ErrConnClosed, err)
		}
	}()

	if err := cli.Close(); err != nil {
		t.Fatal(err)
	}
	clus.TakeClient(0)

	select {
	case <-time.After(3 * time.Second):
		t.Fatal("kv.Get took too long")
	case <-donec:
	}
}
