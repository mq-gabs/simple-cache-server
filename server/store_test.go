package main

import (
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	s := NewStore()

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
	s := NewStore()

	if err := s.SetWithTTL("name", []byte("John Doe"), 1000); err != nil {
		t.Fatal(err)
	}

	res, err := s.Get("name")

	if err != nil {
		t.Fatal(err)
	}

	if string(*res) != "John Doe" {
		t.Fatalf("invalid response: %v", res)
	}

	time.Sleep(time.Second * 2)

	res, err = s.Get("name")

	if err == nil {
		t.Fatalf("error should not be nil | response: %v", res)
	}
}
