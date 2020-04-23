package ticker

import (
	"sort"
	"sync"
)

type jobMap struct {
	data sync.Map
}

// Delete the key from the map
func (m *jobMap) Delete(key string) {
	m.data.Delete(key)
}

// Load the key from the map.
// Returns Listener or bool.
// A false return indicates either the key was not found
// or the value is not of type Listener
func (m *jobMap) Load(key string) (*Job, bool) {
	i, ok := m.data.Load(key)
	if !ok {
		return nil, false
	}
	s, ok := i.(*Job)
	return s, ok
}

// LoadOrStore will return an existing key or
// store the value if not already in the map
func (m *jobMap) LoadOrStore(key string, value *Job) (*Job, bool) {
	i, _ := m.data.LoadOrStore(key, value)
	s, ok := i.(*Job)
	return s, ok
}

// Range over the Listener values in the map
func (m *jobMap) Range(f func(key string, value *Job) bool) {
	m.data.Range(func(k, v interface{}) bool {
		key, ok := k.(string)
		if !ok {
			return false
		}
		value, ok := v.(*Job)
		if !ok {
			return false
		}
		return f(key, value)
	})
}

// Store a Listener in the map
func (m *jobMap) Store(key string, value *Job) {
	m.data.Store(key, value)
}

// Keys returns a list of keys in the map
func (m *jobMap) Keys() []string {
	var keys []string
	m.Range(func(key string, value *Job) bool {
		keys = append(keys, key)
		return true
	})
	sort.Strings(keys)
	return keys
}

func (m *jobMap) Count() int {
	var count int
	m.data.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}
