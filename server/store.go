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

type TTLKey struct {
	key  string
	next *TTLKey
}

type Store struct {
	data        map[string]StoreItem
	lastChange  time.Time
	ttlFirstKey *TTLKey
}

func NewStore() *Store {
	s := &Store{
		data: make(map[string]StoreItem),
	}

	go s.checkTTL()

	return s
}

func (s *Store) checkTTL() {
	defer func() {
		time.Sleep(time.Second)
		go s.checkTTL()
	}()

	if s.ttlFirstKey == nil {
		return
	}

	var prevTTLKey *TTLKey
	ttlKey := s.ttlFirstKey

	for ttlKey != nil {
		if s.data[ttlKey.key].DieIn.After(time.Now()) {
			delete(s.data, ttlKey.key)
			if prevTTLKey != nil {
				prevTTLKey.next = ttlKey.next
			} else {
				s.ttlFirstKey = ttlKey.next
			}
		} else {
			prevTTLKey = ttlKey
		}

		ttlKey = ttlKey.next
	}
}

func (s *Store) addKeyToTTL(key string) {
	ttlKey := TTLKey{
		key: key,
	}

	if s.ttlFirstKey == nil {
		s.ttlFirstKey = &ttlKey
		return
	}

	currTTLKey := s.ttlFirstKey

	for currTTLKey.next != nil {
		currTTLKey = currTTLKey.next
	}

	currTTLKey.next = &ttlKey
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

	s.addKeyToTTL(key)

	mutex.Unlock()

	s.setLastChange()

	return nil
}
