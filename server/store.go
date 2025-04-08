package main

import (
	"fmt"
	"sync"
	"time"
)

var mutex = sync.Mutex{}

type StoreItem struct {
	DieIn *time.Time
	Value *[]byte
}

type Store struct {
	data       map[string]StoreItem
	lastChange time.Time
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]StoreItem),
	}
}

func (s *Store) setLastChange() {
	s.lastChange = time.Now()
}

func (s *Store) Get(key string) (*[]byte, error) {
	item, ok := s.data[key]

	if !ok {
		return nil, fmt.Errorf("key does not exists: %s", key)
	}

	return item.Value, nil
}

func (s *Store) Set(key string, value []byte) error {
	mutex.Lock()

	item := StoreItem{
		DieIn: nil,
		Value: &value,
	}

	s.data[key] = item

	s.setLastChange()

	mutex.Unlock()

	return nil
}

func (s *Store) Erase(key string) error {
	mutex.Lock()

	delete(s.data, key)

	s.setLastChange()

	mutex.Unlock()

	return nil
}

func (s *Store) SetWithTTL(key string, value []byte, ttl uint32) error {
	dieIn := time.Now().Local().Add(time.Second * time.Duration(ttl))

	item := StoreItem{
		DieIn: &dieIn,
		Value: &value,
	}

	mutex.Lock()

	s.data[key] = item

	mutex.Unlock()

	s.setLastChange()

	return nil
}
