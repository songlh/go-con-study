package bolt

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"testing/quick"
	"time"

	"github.com/stretchr/testify/assert"
)

// Ensure that multiple threads can use the DB without race detector errors.
func TestParallelTransactions(t *testing.T) {
	var mutex sync.RWMutex

	err := quick.Check(func(numReaders, batchSize uint, items testdata) bool {
		// Limit the readers & writers to something reasonable.
		numReaders = (numReaders % 10) + 1
		batchSize = (batchSize % 50) + 1

		// Maintain the current dataset.
		var current testdata

		withOpenDB(func(db *DB, path string) {
			db.CreateBucket("widgets")

			// Maintain a set of concurrent readers.
			var wg sync.WaitGroup
			var c = make(chan bool, 0)
			go func() {
				var readers = make(chan int, numReaders)
				for {
					wg.Add(1)

					// Attempt to start a new reader unless we're stopped.
					select {
					case readers <- 0:
					case <-c:
						wg.Done()
						return
					}

					go func() {
						mutex.RLock()
						local := current
						txn, err := db.Transaction()
						mutex.RUnlock()
						if !assert.NoError(t, err) {
							t.FailNow()
						}

						// Verify all data is in for local data list.
						for _, item := range local {
							value, err := txn.Get("widgets", item.Key)
							if !assert.NoError(t, err) || !assert.Equal(t, value, item.Value) {
								txn.Close()
								wg.Done()
								t.FailNow()
							}
						}

						txn.Close()
						wg.Done()
						<-readers
					}()
				}
			}()

			// Batch insert items.
			pending := items
			for {
				// Determine next batch.
				currentBatchSize := int(batchSize)
				if currentBatchSize > len(pending) {
					currentBatchSize = len(pending)
				}
				batchItems := pending[0:currentBatchSize]
				pending = pending[currentBatchSize:]

				// Start write transaction.
				txn, err := db.RWTransaction()
				if !assert.NoError(t, err) {
					t.FailNow()
				}

				// Insert whole batch.
				for _, item := range batchItems {
					err := txn.Put("widgets", item.Key, item.Value)
					if !assert.NoError(t, err) {
						t.FailNow()
					}
				}

				// Commit and update the current list.
				mutex.Lock()
				err = txn.Commit()
				current = append(current, batchItems...)
				mutex.Unlock()
				if !assert.NoError(t, err) {
					t.FailNow()
				}

				// If there are no more left then exit.
				if len(pending) == 0 {
					break
				}

				time.Sleep(1 * time.Millisecond)
			}

			// Notify readers to stop.
			close(c)

			// Wait for readers to finish.
			wg.Wait()
		})
		fmt.Fprint(os.Stderr, ".")
		return true
	}, qconfig())
	assert.NoError(t, err)
	fmt.Fprint(os.Stderr, "\n")
}
