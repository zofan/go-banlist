package banlist

import (
	"sync"
	"time"
)

type BanList struct {
	keys map[uint64]time.Time

	mu sync.Mutex
}

func New() *BanList {
	return &BanList{
		keys: make(map[uint64]time.Time),
	}
}

func (bl *BanList) Ban(k uint64, expires time.Time) {
	bl.mu.Lock()
	bl.keys[k] = expires
	bl.mu.Unlock()
}

func (bl *BanList) UnBan(k uint64) {
	bl.mu.Lock()
	delete(bl.keys, k)
	bl.mu.Unlock()
}

func (bl *BanList) IsBanned(k uint64) (expires *time.Time, ok bool) {
	bl.mu.Lock()
	if e, ok := bl.keys[k]; ok {
		return &e, true
	}
	bl.mu.Unlock()

	return nil, false
}

func (bl *BanList) Clean() {
	bl.mu.Lock()
	now := time.Now()
	for k, e := range bl.keys {
		if now.After(e) {
			delete(bl.keys, k)
		}
	}
	bl.mu.Unlock()
}
