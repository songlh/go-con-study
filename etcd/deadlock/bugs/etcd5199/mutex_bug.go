// Copyright 2016 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package concurrency

import (
	"context"
	"sync"

	v3 "github.com/coreos/etcd/clientv3"
)

// MutexBug implements the sync Locker interface with etcd
// MutexBug implements the sync Locker interface with etcd
type MutexBug struct {
	client *v3.Client

	pfx   string
	myKey string
	myRev int64
}

func NewMutexBug(client *v3.Client, pfx string) *MutexBug {
	return &MutexBug{client, pfx, "", -1}
}

// Lock locks the MutexBug with a cancellable context. If the context is cancelled
// while trying to acquire the lock, the MutexBug tries to clean its stale lock entry.
func (m *MutexBug) LockBug(ctx context.Context) error {
	s, err := NewSession(m.client)
	if err != nil {
		return err
	}
	m.myKey, m.myRev, err = NewUniqueKey(ctx, m.client, m.pfx, v3.WithLease(s.Lease()))
	// wait for deletion revisions prior to myKey
	err = waitDeletesBug(ctx, m.client, m.pfx, v3.WithPrefix(), v3.WithRev(m.myRev-1))
	// release lock key if cancelled
	select {
	case <-ctx.Done():
		m.Unlock()
	default:
	}
	return err
}

func (m *MutexBug) Unlock() error {
	if _, err := m.client.Delete(m.client.Ctx(), m.myKey); err != nil {
		return err
	}
	m.myKey = "\x00"
	m.myRev = -1
	return nil
}

func (m *MutexBug) IsOwner() v3.Cmp {
	return v3.Compare(v3.CreateRevision(m.myKey), "=", m.myRev)
}

func (m *MutexBug) Key() string { return m.myKey }

type lockerMutexBug struct{ *MutexBug }

func (lm *lockerMutexBug) Lock() {
	if err := lm.MutexBug.LockBug(lm.client.Ctx()); err != nil {
		panic(err)
	}
}
func (lm *lockerMutexBug) Unlock() {
	if err := lm.MutexBug.Unlock(); err != nil {
		panic(err)
	}
}

// NewLocker creates a sync.Locker backed by an etcd MutexBug.
func NewLocker_Bug(client *v3.Client, pfx string) sync.Locker {
	return &lockerMutexBug{NewMutexBug(client, pfx)}
}
