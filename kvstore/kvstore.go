package kvstore

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

// KVStore is a simple in-memory key-value store.
type KVStore struct {
	mu    sync.RWMutex
	store map[string]string
}

// NewKVStore initializes a new KVStore instance.
func NewKVStore() *KVStore {
	return &KVStore{
		store: make(map[string]string),
	}
}

// Get retrieves the value associated with the given key.
// It returns the value and a boolean indicating whether the key was found.
func (kv *KVStore) Get(key string) (string, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	value, found := kv.store[key]
	return value, found
}

// Set adds or updates the value associated with the given key.
func (kv *KVStore) Set(key, value string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.store[key] = value
	log.Printf("KVStore: Set key=%s, value=%s\n", key, value)
}

func (kv *KVStore) SaveToFile(filename string) error {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	data, err := json.MarshalIndent(kv.store, "", "  ")
	if err != nil {
		return err
	}

	log.Printf("Saving data to file %s: %s\n", filename, string(data))

	return os.WriteFile(filename, data, 0644)
}

func (kv *KVStore) LoadFromFile(filename string) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	log.Printf("Loaded data from file %s\n", filename)
	return json.Unmarshal(data, &kv.store)
}
