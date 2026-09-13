package cache_test

import (
	"bytes"
	"fmt"
	"scas/cache"
	"sync"
	"testing"
	"time"
)

func TestCacheSetGet(t *testing.T) {
	c := cache.New()

	want := []byte("hello")
	c.Set("foo", want)

	got, ok := c.Get("foo")
	if !ok {
		t.Fatal("expected key to exist")
	}

	if !bytes.Equal(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestCacheGetMissing(t *testing.T) {
	c := cache.New()

	got, ok := c.Get("missing")
	if ok {
		t.Fatal("expected key to be missing")
	}
	if got != nil {
		t.Fatalf("expected nil value, got %v", got)
	}
}

func TestCacheDelete(t *testing.T) {
	c := cache.New()

	c.Set("foo", []byte("bar"))
	c.Delete("foo")

	if _, ok := c.Get("foo"); ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestCacheHas(t *testing.T) {
	c := cache.New()

	if c.Has("foo") {
		t.Fatal("expected key to not exist")
	}

	c.Set("foo", []byte("bar"))

	if !c.Has("foo") {
		t.Fatal("expected key to exist")
	}
}

func TestCacheLen(t *testing.T) {
	c := cache.New()

	if got := c.Len(); got != 0 {
		t.Fatalf("got %d, want 0", got)
	}

	c.Set("a", []byte("1"))
	c.Set("b", []byte("2"))

	if got := c.Len(); got != 2 {
		t.Fatalf("got %d, want 2", got)
	}
}

func TestCacheFlush(t *testing.T) {
	c := cache.New()

	c.Set("a", []byte("1"))
	c.Set("b", []byte("2"))

	c.Flush()

	if got := c.Len(); got != 0 {
		t.Fatalf("got %d, want 0", got)
	}

	if c.Has("a") || c.Has("b") {
		t.Fatal("expected cache to be empty")
	}
}

func TestCacheSetCopiesValue(t *testing.T) {
	c := cache.New()

	value := []byte("hello")
	c.Set("foo", value)

	value[0] = 'H'

	got, _ := c.Get("foo")

	if bytes.Equal(got, value) {
		t.Fatal("cache stored original slice instead of a copy")
	}

	if string(got) != "hello" {
		t.Fatalf("got %q, want %q", got, "hello")
	}
}

func TestCacheGetReturnsCopy(t *testing.T) {
	c := cache.New()

	c.Set("foo", []byte("hello"))

	got, _ := c.Get("foo")
	got[0] = 'H'

	again, _ := c.Get("foo")

	if string(again) != "hello" {
		t.Fatalf("cache value was modified: got %q", again)
	}
}

func TestCacheConcurrentAccess(t *testing.T) {
	c := cache.New()

	const (
		writers = 10
		readers = 10
		ops     = 1000
	)

	var wg sync.WaitGroup

	for i := range writers {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			for j := range ops {
				key := fmt.Sprintf("key-%d", j%100)
				value := []byte(fmt.Sprintf("%d-%d", id, j))
				c.Set(key, value)
			}
		}(i)
	}

	for range readers {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for j := range ops {
				key := fmt.Sprintf("key-%d", j%100)
				c.Get(key)
				c.Has(key)
				c.Len()
			}
		}()
	}

	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestCacheConcurrentWritesAndVerify(t *testing.T) {
	c := cache.New()

	const entries = 10000

	var wg sync.WaitGroup

	for i := range entries {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := fmt.Sprintf("key-%d", i)
			value := []byte(fmt.Sprintf("value-%d", i))

			c.Set(key, value)
		}(i)
	}

	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout")
	}

	for i := range entries {
		key := fmt.Sprintf("key-%d", i)
		expected := []byte(fmt.Sprintf("value-%d", i))

		got, ok := c.Get(key)
		if !ok {
			t.Fatalf("missing key %q", key)
		}

		if !bytes.Equal(got, expected) {
			t.Fatalf("key %q: got %q, want %q", key, got, expected)
		}
	}

	if got := c.Len(); got != entries {
		t.Fatalf("got %d entries, want %d", got, entries)
	}
}

func TestCacheOptions(t *testing.T) {
	t.Run("block overwrite", func(t *testing.T) {
		c := cache.New(cache.WithBlockOverwrite())
		key := "name"
		value := []byte("john")
		otherValue := []byte("bob")

		if !c.Set(key, value) {
			t.Fatal("failed to set at first time")
		}

		got, ok := c.Get(key)
		if !ok {
			t.Fatal("value was not written")
		}
		if !bytes.Equal(got, value) {
			t.Fatal("written value is different")
		}

		if c.Set(key, otherValue) {
			t.Fatal("should not overwrite key")
		}

		got, ok = c.Get(key)
		if !ok {
			t.Fatal("value was removed somehow")
		}
		if !bytes.Equal(got, value) {
			t.Fatal("written value is different after try to overwrite")
		}
	})
}
