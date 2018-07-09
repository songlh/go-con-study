package txfun

import (
	"bytes"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/boltdb/bolt"
)

func printState(db *bolt.DB) {
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("data")).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Println(string(k), "=>", string(v))
		}
		return nil
	})
}

func Test1(t *testing.T) {
	db, _ := NewDB()
	tx := db.Begin()

	tx.Set([]byte("foo"), []byte("bar"))

	err := tx.Commit()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("We should just see `foo => bar`")
	printState(db.state)

	tx = db.Begin()
	val, err := tx.Get([]byte("foo"))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(val, []byte("bar")) != 0 {
		t.Fatalf("Expected %v, got %v", string([]byte("bar")), val)
	}

	// Why would this cause a panic?
	//tx.Set([]byte("foo"), append(val, []byte("baz")...))

	fmt.Println("We should just see `foo => bar`")
	printState(db.state)
	tx.Commit()
	fmt.Println("We should just see `foo => bar`")
	printState(db.state)

	txA := db.Begin()
	txB := db.Begin()

	fmt.Println("We should just see `foo => bar`")
	printState(db.state)
	txA.Set([]byte("some_key"), []byte("a"))
	txA.Commit()

	fmt.Println("We should now also see `some_key => a`")
	printState(db.state)

	txB.Set([]byte("some_key"), []byte("b"))

	fmt.Println("Same thing again")
	printState(db.state)

	t.Log(txB.Commit())

	fmt.Println("We should see the same thing because there was a conflict")
	printState(db.state)

	fmt.Println("We should now see `some_key => b` because we're retrying")
	t.Log(txB.Commit())
	printState(db.state)

	db.state.Close()
}

func TestConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	const N = 55 // # goroutines
	var cases []map[string]bool
	keys := map[string]bool{}

	for i := 0; i < N; i++ {
		m := map[string]bool{}
		for j := 0; j < 50; j++ {
			str := fmt.Sprint(rand.Int())
			m[str] = false
			keys[str] = false
		}

		cases = append(cases, m)
	}

	db, _ := NewDB()

	wg := sync.WaitGroup{}

	for i := 0; i < N; i++ {
		wg.Add(1)

		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		go func(n int) {
			tx := db.Begin()
			t.Logf("[tx %d] starting tx", tx.id)

			for key := range cases[n] {
				tx.Set([]byte(key), []byte(fmt.Sprint(n)))
			}

			txStart := time.Now()
			for err := tx.Commit(); err != nil; err = tx.Commit() {
				t.Logf("[tx %d] retrying", tx.id)
				txStart = time.Now()
			}
			t.Logf("[tx %d] took %v to commit", tx.id, time.Now().Sub(txStart))
			wg.Done()
		}(i)
	}

	wg.Wait()

	err := db.state.View(func(boltTx *bolt.Tx) error {
		c := boltTx.Bucket([]byte("data")).Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			keys[string(k)] = true
		}

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	for key, val := range keys {
		if !val {
			t.Errorf("key %v was not found", key)
		}
	}
}
