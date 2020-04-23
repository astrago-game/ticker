package ticker

import (
	"sort"
	"sync"
)

type workerMap struct {
	data sync.Map
}

// Delete the key from the map
func (m *workerMap) Delete(key string) {
	m.data.Delete(key)
}

// Load the key from the map.
// Returns Listener or bool.
// A false return indicates either the key was not found
// or the value is not of type Listener
func (m *workerMap) Load(key string) (Worker, bool) {
	i, ok := m.data.Load(key)
	if !ok {
		return nil, false
	}
	s, ok := i.(Worker)
	return s, ok
}

// LoadOrStore will return an existing key or
// store the value if not already in the map
func (m *workerMap) LoadOrStore(key string, value Worker) (Worker, bool) {
	i, _ := m.data.LoadOrStore(key, value)
	s, ok := i.(Worker)
	return s, ok
}

// Range over the Listener values in the map
func (m *workerMap) Range(f func(key string, value Worker) bool) {
	m.data.Range(func(k, v interface{}) bool {
		key, ok := k.(string)
		if !ok {
			return false
		}
		value, ok := v.(Worker)
		if !ok {
			return false
		}
		return f(key, value)
	})
}

// Store a Listener in the map
func (m *workerMap) Store(key string, value Worker) {
	m.data.Store(key, value)
}

// Keys returns a list of keys in the map
func (m *workerMap) Keys() []string {
	var keys []string
	m.Range(func(key string, value Worker) bool {
		keys = append(keys, key)
		return true
	})
	sort.Strings(keys)
	return keys
}
