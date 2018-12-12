package badger

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetSetDeadlock(t *testing.T) {
	dir, err := ioutil.TempDir("", "badger")
	fmt.Println(dir)
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	opt := DefaultOptions
	opt.Dir = dir
	opt.ValueDir = dir
	opt.ValueLogFileSize = 1 << 20
	db, err := Open(opt)
	require.NoError(t, err)
	defer db.Close()

	val := make([]byte, 1<<19)
	key := []byte("key1")
	require.NoError(t, db.Update(func(txn *Txn) error {
		rand.Read(val)
		require.NoError(t, txn.Set(key, val))
		return nil
	}))

	timeout, done := time.After(10*time.Second), make(chan bool)

	go func() {
		db.Update(func(txn *Txn) error {
			item, err := txn.Get(key)
			require.NoError(t, err)
			_, err = item.Value() // This take a RLock on file
			require.NoError(t, err)

			rand.Read(val)
			require.NoError(t, txn.Set(key, val))
			require.NoError(t, txn.Set([]byte("key2"), val))
			return nil
		})
		fmt.Println("finish updating.....")
		done <- true
	}()

	select {
	case <-timeout:
		fmt.Println("timeout....")
		t.Fatal("db.Update did not finish within 10s, assuming deadlock.")
	case <-done:
		t.Log("db.Update finished.")
	}
}
