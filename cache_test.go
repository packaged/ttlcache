package ttlmap_test

import (
	"github.com/packaged/ttlmap"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	cache := ttlmap.New(ttlmap.WithDefaultTTL(300*time.Millisecond), ttlmap.WithCleanupDuration(time.Millisecond*100))

	data, exists := cache.Get("hello")
	if exists || data != nil {
		t.Errorf("Expected empty cache to return no data")
	}

	cache.Set("hello", "world", nil)
	data, exists = cache.Get("hello")
	if !exists {
		t.Errorf("Expected cache to return data for `hello`")
	}
	if data.(string) != "world" {
		t.Errorf("Expected cache to return `world` for `hello`")
	}

	// Check to see if cleanup is clearing unexpired items
	time.Sleep(time.Millisecond * 200)
	data, exists = cache.Get("hello")
	if !exists || data == nil {
		t.Errorf("Expected cache to return data")
	}

	// Check Cache is re-touching after a get
	time.Sleep(time.Millisecond * 200)
	data, exists = cache.Get("hello")
	if !exists || data == nil {
		t.Errorf("Expected cache to return data")
	}

	// Check Cache is optionally re-touching after a get
	time.Sleep(time.Millisecond * 200)
	data, exists = cache.TouchGet("hello", false)
	if !exists || data == nil {
		t.Errorf("Expected cache to return data")
	}

	// Make sure cache clears after expiry
	time.Sleep(time.Millisecond * 200)
	data, exists = cache.Get("hello")
	if exists || data != nil {
		t.Errorf("Expected empty cache to return no data")
	}
}

func TestMaxLifetime(t *testing.T) {
	cache := ttlmap.New(ttlmap.WithMaxLifetime(time.Millisecond * 100))

	data, exists := cache.Get("hello")
	if exists || data != nil {
		t.Errorf("Expected empty cache to return no data")
	}

	cache.Set("hello", "world", nil)
	data, exists = cache.Get("hello")
	if !exists {
		t.Errorf("Expected cache to return data for `hello`")
	}
	if data.(string) != "world" {
		t.Errorf("Expected cache to return `world` for `hello`")
	}

	// Check to see if max lifetime has killed the item
	time.Sleep(time.Millisecond * 200)
	data, exists = cache.Get("hello")
	if exists || data != nil {
		t.Errorf("Expected empty cache to return no data")
	}
}

func TestItems(t *testing.T) {
	dur := time.Millisecond * 100

	cache := ttlmap.New()
	cache.Set("item1", "one", nil)
	cache.Set("item2", "two", nil)
	cache.Set("item3", "three", &dur)
	cache.Set("item4", "four", nil)

	if len(cache.Items()) != 4 {
		t.Errorf("Expected cache to return 4 items")
	}

	time.Sleep(dur)

	if len(cache.Items()) != 3 {
		t.Errorf("Expected cache to return 3 items after cache expiry")
	}
}

func TestCleanup(t *testing.T) {
	dur := time.Millisecond * 100

	cleanup := 0

	cache := ttlmap.New(ttlmap.WithCleanupDuration(time.Millisecond * 5))

	cache.SetWithCleanup("item1", "one", nil, func(item *ttlmap.Item) { cleanup++ })

	if cleanup != 0 {
		t.Errorf("Cache item cleaned up too early")
	}

	cache.Remove("item1")

	if cleanup != 1 {
		t.Errorf("Cache item cleanup not called on Remove")
	}

	cache.SetWithCleanup("item2", "two", &dur, func(item *ttlmap.Item) { cleanup++ })
	time.Sleep(dur)
	// Wait a few more milliseconds for cleanup to run
	time.Sleep(time.Millisecond * 20)

	if cleanup != 2 {
		t.Errorf("Cache item cleanup not called on expiry")
	}
}
