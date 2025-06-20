package callback

import (
	"maps"
	"sync"
)

// Manager manages the lifecycle of a collection of generic [Callback] structs.
type Manager[T any] struct {
	callbacks map[*Callback[T]]bool
	mux       sync.RWMutex
}

// NewManager constructs a new [Manager] structs.
func NewManager[T any]() *Manager[T] {
	return &Manager[T]{
		callbacks: make(map[*Callback[T]]bool),
	}
}

// Register adds a [Callback] struct to the map.
func (m *Manager[T]) Register(c *Callback[T]) *Callback[T] {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.callbacks[c] = true
	return c
}

// Reregister removes a [Callback] struct from the map.
func (m *Manager[T]) Deregister(c *Callback[T]) *Callback[T] {
	m.mux.Lock()
	defer m.mux.Unlock()
	delete(m.callbacks, c)
	return c
}

// Reset clears all [Callback] struct from the map.
func (m *Manager[T]) Reset() {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.callbacks = make(map[*Callback[T]]bool)
}

type Action[T any] func(*Event[T])

// Recurring adds a recurring [Callback] struct to the map.
func (m *Manager[T]) Recurring(action Action[T]) *Callback[T] {
	return m.Register(&Callback[T]{
		Action:  action,
		Enabled: true,
	})
}

// Once adds a [Callback] to the map that is deregistered after first execution.
func (m *Manager[T]) Once(action Action[T]) *Callback[T] {
	callback := &Callback[T]{Enabled: true}
	callback.Action = func(e *Event[T]) {
		action(e)
		callback.Enabled = false
	}
	return m.Register(callback)
}

// SleepUntilDisabled adds a [Callback] struct to the map and pauses the current goroutine until the callback is disabled.
func (m *Manager[T]) SleepUntilDisabled(action Action[T]) *Callback[T] {
	callback := &Callback[T]{Enabled: true}
	var wg sync.WaitGroup
	wg.Add(1)
	callback.Action = func(e *Event[T]) {
		go func() {
			action(e)
			if !e.Callback.Enabled {
				wg.Done()
			}
		}()
	}
	m.Register(callback)
	wg.Wait()
	return callback
}

// Map returns a clone of the internal [Callback] map.
func (m *Manager[T]) Map() map[*Callback[T]]bool {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return maps.Clone(m.callbacks)
}

// Call fans out [Event] objects across all callbacks in the map.
func (m *Manager[T]) Call(v T) {
	callbacks := m.Map()
	for callback := range callbacks {
		if callback.Enabled {
			callback.Call(v)
		}
	}
	m.mux.Lock()
	for callback := range callbacks {
		if !callback.Enabled {
			delete(m.callbacks, callback)
		}
	}
	m.mux.Unlock()
}

// Callback contains the function reference and parameters.
type Callback[T any] struct {
	Action  Action[T]
	Enabled bool
}

// Call constructs an [Event] object and passes them to the internal Action function.
func (c *Callback[T]) Call(v T) {
	c.Action(&Event[T]{
		Data:     v,
		Callback: c,
	})
}

// Event contains the content and reference to the [Callback] struct.
type Event[T any] struct {
	Data     T
	Callback *Callback[T]
}
