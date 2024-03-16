package pokecache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	interval := time.Millisecond * 10
	cache := NewCache(interval)
	if cache.cache == nil {
		t.Error("Cache is nil")
	}
}

func TestCache_Add(t *testing.T) {
	interval := time.Millisecond * 10
	cache := NewCache(interval)
	cases := []struct {
		key   string
		value []byte
	}{
		{"key1", []byte("value1")},
		{"key2", []byte("value2")},
	}

	for _, c := range cases {
		cache.Add(c.key, c.value)
		if _, ok := cache.cache[c.key]; !ok {
			t.Error("Key not added to cache")
		}
		actual, err := cache.Get(c.key)
		if err != true {
			t.Error("Key not found in cache")
		}
		if string(actual) != string(c.value) {
			t.Errorf("wanted value \"%s\" does not match got \"%s\"", c.value, actual)
		}
	}
}

func TestCache_reap(t *testing.T) {
	interval := time.Millisecond * 10
	cache := NewCache(interval)
	cases := []struct {
		key   string
		value []byte
	}{
		{"key1", []byte("value1")},
		{"key2", []byte("value2")},
	}

	for _, c := range cases {
		cache.Add(c.key, c.value)
		time.Sleep(interval * 2)
		if _, ok := cache.Get(c.key); ok {
			t.Errorf("cache entry %v not reaped from cache", c.key)
		}
	}
}

func TestCache_reapFail(t *testing.T) {
	interval := time.Millisecond * 10
	cache := NewCache(interval)
	cases := []struct {
		key   string
		value []byte
	}{
		{"key1", []byte("value1")},
		{"key2", []byte("value2")},
	}

	for _, c := range cases {
		cache.Add(c.key, c.value)
		time.Sleep(interval / 2)
		if _, ok := cache.Get(c.key); !ok {
			t.Errorf("cache entry %v reaped from cache when it should have not been", c.key)
		}
	}

}
