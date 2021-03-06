package banlist

import (
	"sync"
	"time"
)

type BanList struct {
	keys map[uint64]time.Time

	mu sync.RWMutex
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

func (bl *BanList) All() map[uint64]time.Time {
	bl.mu.RLock()
	defer bl.mu.RUnlock()
	return bl.keys
}

func (bl *BanList) IsBanned(k uint64) (expires *time.Time, ok bool) {
	bl.mu.RLock()
	defer bl.mu.RUnlock()
	if e, ok := bl.keys[k]; ok {
		return &e, true
	}

	return nil, false
}

func (bl *BanList) Clear() {
	now := time.Now()
	bl.mu.Lock()
	for k, e := range bl.keys {
		if now.After(e) {
			delete(bl.keys, k)
		}
	}
	bl.mu.Unlock()
}
