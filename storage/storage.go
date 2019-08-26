package storage

import (
	"sync"
	"time"
)

type sku struct {
	Obj interface{}
	Exp int64
}

// Storage represents a cache interface
type Storage interface {
	Get(key string) (item interface{}, found bool)
	Put(key string, item interface{}, dur time.Duration)
}

type storage struct {
	sync.RWMutex
	data map[string]sku
}

// NewStore creates a new cache
func NewStore() Storage {
	stor := storage{data: make(map[string]sku)}
	// expiration
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			stor.Expire()
		}
	}()
	return &stor
}

func (s *storage) Expire() {
	s.Lock()
	now := time.Now().Unix()
	for k, v := range s.data {
		if now > v.Exp {
			delete(s.data, k)
		}
	}
	s.Unlock()
}

// Get returns a cache record
func (s *storage) Get(key string) (interface{}, bool) {
	s.RLock()
	entry, found := s.data[key]
	s.RUnlock()
	if !found {
		return nil, found
	}
	return entry.Obj, found
}

// Store puts a record into cache
func (s *storage) Put(key string, item interface{}, dur time.Duration) {
	s.Lock()
	_, found := s.data[key]
	if !found {
		s.data[key] = sku{Obj: item, Exp: time.Now().Add(dur).Unix()}
	}
	s.Unlock()
}
