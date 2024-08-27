package kvstore

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestKVStore_SetAndGet(t *testing.T) {
	store := NewKVStore()

	tests := []struct {
		name     string
		key      string
		value    string
		expected string
		found    bool
	}{
		{"Set and Get key1", "key1", "value1", "value1", true},
		{"Set and Get key2", "key2", "value2", "value2", true},
		{"Set and Get empty key", "", "empty", "empty", true},
		{"Get non-existent key", "noKey", "", "", false},
	}

	for _, tt := range tests {
		// If we are setting a value, do that first.
		if tt.found {
			store.Set(tt.key, tt.value)
		}

		t.Run(tt.name, func(t *testing.T) {
			result, found := store.Get(tt.key)
			if found != tt.found {
				t.Fatalf("expected found=%v, got %v", tt.found, found)
			}
			if result != tt.expected {
				t.Fatalf("expected value=%s, got %s", tt.expected, result)
			}
		})
	}
}

func TestKVStore_SaveAndLoad(t *testing.T) {
	store := NewKVStore()

	// Set some key-value pairs
	store.Set("key1", "value1")
	store.Set("key2", "value2")

	// Save the store to a file
	filename := "test_kvstore.json"
	defer os.Remove(filename) // Clean up the test file after the test

	if err := store.SaveToFile(filename); err != nil {
		t.Fatalf("failed to save store to file: %v", err)
	}

	// Create a new KVStore and load the data from the file
	newStore := NewKVStore()
	if err := newStore.LoadFromFile(filename); err != nil {
		t.Fatalf("failed to load store from file: %v", err)
	}

	// Verify that the loaded data matches the original data
	value, found := newStore.Get("key1")
	if !found || value != "value1" {
		t.Fatalf("expected 'key1' to have value 'value1', got '%s'", value)
	}

	value, found = newStore.Get("key2")
	if !found || value != "value2" {
		t.Fatalf("expected 'key2' to have value 'value2', got '%s'", value)
	}
}

func TestKVStore_ConcurrentAccess(t *testing.T) {
	store := NewKVStore()

	var wg sync.WaitGroup
	numGoroutines := 15

	// Function to set key-value pairs concurrently
	setFunc := func(id int) {
		defer wg.Done()
		store.Set("key"+strconv.Itoa(id), "value"+strconv.Itoa(id))
	}

	// Function to get key-value pairs concurrently
	getFunc := func(id int) {
		defer wg.Done()
		value, found := store.Get("key" + strconv.Itoa(id))
		if found && value != "value"+strconv.Itoa(id) {
			t.Errorf("expected value%d, got %s", id, value)
		}
	}

	// Start concurrent Set operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go setFunc(i)
	}

	// Wait for all Set operations to finish
	wg.Wait()

	// Start concurrent Get operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go getFunc(i)
	}

	// Wait for all Get operations to finish
	wg.Wait()
}

func TestKVStore_DeadlockDetection(t *testing.T) {
	store := NewKVStore()
	done := make(chan bool)

	go func() {
		defer close(done)
		store.Set("key", "value")
		store.Get("key")
	}()

	select {
	case <-done:
		// Test passed, no deadlock
	case <-time.After(2 * time.Second):
		t.Fatal("deadlock detected")
	}
}

func TestKVStore_RandomizedStress(t *testing.T) {
	store := NewKVStore()
	var wg sync.WaitGroup
	numGoroutines := 100

	randomOps := func(id int) {
		defer wg.Done()
		key := fmt.Sprintf("key%d", id)
		store.Set(key, fmt.Sprintf("value%d", id))
		store.Get(key)
		// Optionally, if a Delete method is implemented:
		// store.Delete(key)
	}

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go randomOps(i)
	}

	wg.Wait()
}
