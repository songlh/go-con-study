package txfun

import (
	"github.com/boltdb/bolt"

	"encoding/binary"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
)

type DB struct {
	state    *bolt.DB
	lock     sync.Mutex
	epoch    uint64
	inflight map[uint32]*Tx
}

func NewDB() (*DB, error) {
	boltdb, err := bolt.Open("/tmp/bolt.db", 0600)
	if err != nil {
		return nil, err
	}

	err = boltdb.Update(func(tx *bolt.Tx) error {
		meta, err := tx.CreateBucketIfNotExists([]byte("meta"))
		if err != nil {
			return err
		}

		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, 1)

		err = meta.Put([]byte("epoch"), buf)
		return err
	})

	return &DB{
		state:    boltdb,
		lock:     sync.Mutex{},
		epoch:    1,
		inflight: make(map[uint32]*Tx),
	}, nil
}

func (db *DB) Begin() *Tx {
	db.lock.Lock()
	defer db.lock.Unlock()

	id := rand.Uint32()

	for _, present := db.inflight[id]; present; id = rand.Uint32() {
	}

	boltTx, err := db.state.Begin(false)
	if err != nil {
		log.Print(err)
	}

	tx := &Tx{
		id:          id,
		db:          db,
		view:        boltTx,
		state:       newList(),
		keysWritten: make(map[string]struct{}),
	}

	db.inflight[id] = tx
	return tx
}

func (db *DB) commitTx(tx *Tx) error {
	if tx.conflicted {
		tx.conflicted = false
		tx.commits = tx.commits[:0]
		tx.view.Commit()
		boltTx, _ := db.state.Begin(false)
		tx.view = boltTx
		return ErrConflict
	}

	db.lock.Lock()

	err := db.state.Update(func(boltTx *bolt.Tx) error {
		b, err := boltTx.CreateBucketIfNotExists([]byte("data"))
		if err != nil {
			return err
		}

		for n := tx.state.root; n != nil; n = n.next {
			err = b.Put(n.key, n.value)
			if err != nil {
				return err
			}
		}

		return nil
	})

	db.lock.Unlock()

	if err != nil {
		return err
	}

	for id, inflightTx := range db.inflight {
		if id == tx.id {
			continue
		}

		inflightTx.addCommits(tx.keysWritten)
	}

	atomic.AddUint64(&db.epoch, 1)
	delete(db.inflight, tx.id)

	return nil
}
