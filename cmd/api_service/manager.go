package main

import (
	"sync"
	"time"
)

type StorageManager struct {
	mutex     sync.RWMutex
	addresses map[string]time.Time
}

func (sm *StorageManager) AddAddress(addr string) {
	sm.mutex.Lock()
	sm.addresses[addr] = time.Now()
	sm.mutex.Unlock()
}

func (sm *StorageManager) DelRange(f func(string, time.Time) bool) {
	sm.mutex.Lock()
	for k, v := range sm.addresses {
		if f(k, v) {
			delete(sm.addresses, k)
		}
	}
	sm.mutex.Unlock()
}
