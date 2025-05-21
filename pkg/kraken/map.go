package kraken

import (
	"fmt"
	"maps"
	"sync"
)

// Map wraps a map and a mutex for concurrent operations.
type Map[K comparable, V any] struct {
	raw map[K]*V
	mux sync.RWMutex
}

// NewMap creates a Map.
func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		raw: make(map[K]*V),
	}
}

// Length returns the size of m.
func (m *Map[K, V]) Length() int {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return len(m.raw)
}

// SetMap sets m's hash map to r.
func (m *Map[K, V]) SetMap(r map[K]*V) {
	if r == nil {
		r = make(map[K]*V)
	}
	m.mux.Lock()
	defer m.mux.Unlock()
	m.raw = r
}

// Get returns the value of k.
func (m *Map[K, V]) Get(k K) (*V, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	value, ok := m.raw[k]
	if !ok {
		return nil, fmt.Errorf("%v not found", k)
	}
	return value, nil
}

// Set assigns v to k.
func (m *Map[K, V]) Set(key K, value *V) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.raw[key] = value
}

// Delete removes k from m.
func (m *Map[K, V]) Delete(key K) {
	m.mux.Lock()
	defer m.mux.Unlock()
	delete(m.raw, key)
}

// Range executes f on all key-value pairs of m.
func (m *Map[K, V]) Range(f func(K, *V)) {
	for key, value := range m.Raw() {
		f(key, value)
	}
}

// Reset assigns an empty hash map to m.
func (m *Map[K, V]) Reset() {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.raw = make(map[K]*V)
}

// Raw returns a copy of the m hash map.
func (m *Map[K, V]) Raw() map[K]*V {
	m.mux.RLock()
	defer m.mux.RUnlock()
	raw := make(map[K]*V)
	maps.Copy(raw, m.raw)
	return raw
}

// Clone returns a copy of m.
func (m *Map[K, V]) Clone() *Map[K, V] {
	c := NewMap[K, V]()
	c.raw = m.Raw()
	return c
}

// Keys returns all keys of m.
func (m *Map[K, V]) Keys() []K {
	keys := make([]K, 0)
	m.Range(func(k K, v *V) {
		keys = append(keys, k)
	})
	return keys
}
