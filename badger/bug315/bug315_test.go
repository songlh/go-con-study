package badger

import (
	"fmt"
	"github.com/dgraph-io/badger/y"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

import (
	"math/rand"
)

func TestWriteDeadlock(t *testing.T) {
	dir, err := ioutil.TempDir("", "badger")
	fmt.Println(dir)
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	opt := DefaultOptions
	opt.Dir = dir
	opt.ValueDir = dir
	opt.ValueLogFileSize = 10 << 20
	db, err := Open(opt)
	require.NoError(t, err)

	print2 := func(count *int) {
		*count++
		if *count%100 == 0 {
			fmt.Printf("%05d\r", *count)
		}
	}

	var count int
	val := make([]byte, 10000)
	require.NoError(t, db.Update(func(txn *Txn) error {
		for i := 0; i < 1500; i++ {
			key := fmt.Sprintf("%d", i)
			rand.Read(val)
			require.NoError(t, txn.Set([]byte(key), val))
			print2(&count)
		}
		return nil
	}))

	count = 0
	fmt.Println("\nWrites done. Iteration and updates starting...")
	err = db.Update(func(txn *Txn) error {
		opt := DefaultIteratorOptions
		opt.PrefetchValues = false
		it := txn.NewIterator(opt)
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()

			// Using Value() would cause deadlock.
			// item.Value()
			out, err := item.Value()
			require.NoError(t, err)
			require.Equal(t, len(val), len(out))

			key := y.Copy(item.Key())
			rand.Read(val)
			require.NoError(t, txn.Set(key, val))
			print2(&count)
		}
		return nil
	})
	require.NoError(t, err)
}
