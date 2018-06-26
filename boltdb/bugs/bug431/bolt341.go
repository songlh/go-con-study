package main

import (
	"fmt"
	"os"

	"github.com/boltdb/bolt"
)

func main() {
	// Open the database.
	db, _ := bolt.Open("test.db", 0666, nil)
	defer os.Remove(db.Path())

	// Start a write transaction.
	db.Update(func(tx *bolt.Tx) error {
		// Create a bucket.
		tx.CreateBucket([]byte("widgets"))
		b := tx.Bucket([]byte("widgets"))

		// Set the value "bar" for the key "foo".
		b.Put([]byte("foo"), []byte("bar"))

		// Retrieve the key back from the database and verify it.
		value := b.Get([]byte("foo"))
		fmt.Printf("The value of 'foo' was: %s\n", value)
		return nil
	})

	// Retrieve the key again.
	tx, _ := db.Begin(false)

	db.Close()
	value := tx.Bucket([]byte("widgets")).Get([]byte("foo"))
	if value == nil {
		fmt.Printf("The value of 'foo' is now: nil\n")
	}
}

