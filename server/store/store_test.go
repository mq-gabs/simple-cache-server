package store

import (
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	s := New()

	if err := s.Set("name", []byte("John Doe")); err != nil {
		t.Fatal(err)
	}

	res, err := s.Get("name")

	if err != nil {
		t.Fatal(err)
	}

	if string(*res) != "John Doe" {
		t.Fatalf("invalid response: %v", res)
	}
}

func TestTTL(t *testing.T) {
	key := "name"
	value := "John Doe"
	keyTTL := 1
	keyDestroyTimeDelay := time.Second

	s := New()
	s.configBasic.keyDestroyTimeDelay = keyDestroyTimeDelay

	if err := s.SetWithTTL(key, []byte(value), uint32(keyTTL)); err != nil {
		t.Fatal(err)
	}

	res, err := s.Get(key)
	if err != nil {
		t.Fatal(err)
	}

	if string(*res) != value {
		t.Fatalf("invalid response: %v", res)
	}

	time.Sleep(keyDestroyTimeDelay + time.Second)

	res, err = s.Get(key)
	if err == nil {
		t.Fatalf("error should not be nil | response: %v", res)
	}
}
