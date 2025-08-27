package store

import (
	"scas/utils"
	"sync"
	"time"
)

const keyDestroyTimeDelayDefault = time.Minute

type StoreItem struct {
	DieIn *time.Time
	Value *[]byte
}

type StoreConfig struct {
	keyDestroyTimeDelay time.Duration
}

type Store struct {
	mu          sync.Mutex
	data        map[string]*StoreItem
	lastChange  time.Time
	configBasic StoreConfig
}

func New() *Store {
	s := &Store{
		mu:   sync.Mutex{},
		data: make(map[string]*StoreItem),
		configBasic: StoreConfig{
			keyDestroyTimeDelay: keyDestroyTimeDelayDefault,
		},
	}

	s.keyDestroyPeriodicStart()

	return s
}

func (s *Store) timeLastChangeUpdate() {
	s.lastChange = time.Now()
}

func (s *Store) keyGet(key string) (*StoreItem, error) {
	item, ok := s.data[key]
	if !ok {
		return nil, utils.FmtErr(errKeyDoesNotExists, key)
	}

	return item, nil
}

func (s *Store) keyAdd(key string, item *StoreItem) error {
	s.data[key] = item

	return nil
}

func (s *Store) keyDestroy(key string) error {
	delete(s.data, key)
	return nil
}

func (s *Store) keyDestroyPeriodicStart() {
	go func() {
		for {
			s.keyDestroyPeriodic()
			time.Sleep(s.configBasic.keyDestroyTimeDelay)
		}
	}()
}

func (s *Store) keyDestroyPeriodic() {
	s.mu.Lock()
	defer s.mu.Unlock()

	timeNow := time.Now()

	for key, item := range s.data {
		if item.DieIn == nil {
			continue
		}

		if timeNow.After(*item.DieIn) {
			s.keyDestroy(key)
		}
	}
}

func (s *Store) Get(key string) (*[]byte, error) {
	item, err := s.keyGet(key)
	if err != nil {
		return nil, err
	}

	return item.Value, nil
}

func (s *Store) Set(key string, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	item := &StoreItem{
		DieIn: nil,
		Value: &value,
	}

	if err := s.keyAdd(key, item); err != nil {
		return nil
	}

	s.timeLastChangeUpdate()

	return nil
}

func (s *Store) Erase(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.keyDestroy(key); err != nil {
		return err
	}

	s.timeLastChangeUpdate()

	return nil
}

func (s *Store) SetWithTTL(key string, value []byte, ttl uint32) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dieIn := time.Now().Add(time.Second * time.Duration(ttl))

	item := &StoreItem{
		DieIn: &dieIn,
		Value: &value,
	}

	if err := s.keyAdd(key, item); err != nil {
		return err
	}

	s.timeLastChangeUpdate()

	return nil
}
