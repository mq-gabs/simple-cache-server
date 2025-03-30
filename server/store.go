package main

import (
	"fmt"
	"sync"
	"time"
)

var mutex = sync.Mutex{}

type Store struct {
	data map[string][]byte
	lastChange time.Time
}

func NewStore() *Store {
	return &Store{
		data: make(map[string][]byte),
	}
}

func (s *Store) setLastChange() {
	s.lastChange = time.Now()
}

func (s *Store) Get(key string) ([]byte, error) {
	value, ok := s.data[key]

	if !ok {
		return nil, fmt.Errorf("key does not exists: %s", key)
	}

	return value, nil
}

func (s *Store) Set(key string, value []byte) error {
	mutex.Lock()
	
	s.data[key] = value

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