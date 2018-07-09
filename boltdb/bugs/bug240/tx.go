package txfun

import (
	"bytes"
	"errors"

	"github.com/boltdb/bolt"
)

var (
	ErrNotFound = errors.New("tx: key not found")
	ErrConflict = errors.New("tx: conflict")
	ErrUnknown  = errors.New("unknown error")
)

type Tx struct {
	id          uint32
	db          *DB
	view        *bolt.Tx
	state       *list
	conflicted  bool
	keysWritten map[string]struct{}
	commits     []map[string]struct{}
}

func (tx *Tx) Rollback() {
	tx.view.Commit()
}

func (tx *Tx) Commit() error {
CONFLICT_CHECK:
	for _, committed := range tx.commits {
		for key := range committed {
			if _, present := tx.keysWritten[key]; present {
				tx.conflicted = true
				break CONFLICT_CHECK
			}
		}
	}
	return tx.db.commitTx(tx)
}

func (tx *Tx) Set(key, value []byte) error {
	tx.keysWritten[string(key)] = struct{}{}
	tx.state.insert(key, value)
	return nil
}

func (tx *Tx) Get(key []byte) ([]byte, error) {
	for n := tx.state.root; n != nil; n = n.next {
		cmp := bytes.Compare(key, n.key)
		switch {
		case cmp == 0:
			return n.value, nil
		case cmp > 0:
			break
		}
	}

	bucket := tx.view.Bucket([]byte("data"))
	if bucket == nil {
		return nil, ErrNotFound
	}

	c := bucket.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		cmp := bytes.Compare(key, k)
		switch {
		case cmp == 0:
			return v, nil
		case cmp > 0:
			break
		}
	}

	return nil, ErrNotFound
}

func (tx *Tx) addCommits(commits ...map[string]struct{}) {
	tx.commits = append(tx.commits, commits...)
}