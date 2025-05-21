package kraken

import (
	"maps"
	"sync"
)

// CallbackManager manages the lifecycle of a collection of generic [Callback] structs.
type CallbackManager[T any] struct {
	callbacks map[*Callback[T]]bool
	mux       sync.RWMutex
}

// NewCallbackManager constructs a new [CallbackManager] structs.
func NewCallbackManager[T any]() *CallbackManager[T] {
	return &CallbackManager[T]{
		callbacks: make(map[*Callback[T]]bool),
	}
}

// Register adds a [Callback] struct to the map.
func (cg *CallbackManager[T]) Register(c *Callback[T]) *Callback[T] {
	cg.mux.Lock()
	defer cg.mux.Unlock()
	cg.callbacks[c] = true
	return c
}

// Reregister removes a [Callback] struct from the map.
func (cg *CallbackManager[T]) Deregister(c *Callback[T]) *Callback[T] {
	cg.mux.Lock()
	defer cg.mux.Unlock()
	delete(cg.callbacks, c)
	return c
}

// Reset clears all [Callback] struct from the map.
func (cg *CallbackManager[T]) Reset() {
	cg.mux.Lock()
	defer cg.mux.Unlock()
	cg.callbacks = make(map[*Callback[T]]bool)
}

type Action[T any] func(*Event[T])

// Recurring adds a recurring [Callback] struct to the map.
func (cg *CallbackManager[T]) Recurring(action Action[T]) *Callback[T] {
	return cg.Register(&Callback[T]{
		Action:  action,
		Enabled: true,
	})
}

// Once adds a [Callback] to the map that is deregistered after first execution.
func (cg *CallbackManager[T]) Once(action Action[T]) *Callback[T] {
	callback := &Callback[T]{Enabled: true}
	callback.Action = func(e *Event[T]) {
		action(e)
		callback.Enabled = false
	}
	return cg.Register(callback)
}

// SleepUntilDisabled adds a [Callback] struct to the map and pauses the current goroutine until the callback is disabled.
func (cg *CallbackManager[T]) SleepUntilDisabled(action Action[T]) *Callback[T] {
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
	cg.Register(callback)
	wg.Wait()
	return callback
}

// Map returns a clone of the internal [Callback] map.
func (cg *CallbackManager[T]) Map() map[*Callback[T]]bool {
	cg.mux.RLock()
	defer cg.mux.RUnlock()
	return maps.Clone(cg.callbacks)
}

// Call fans out [Event] objects across all callbacks in the map.
func (cg *CallbackManager[T]) Call(v T) {
	callbacks := cg.Map()
	for callback := range callbacks {
		if callback.Enabled {
			callback.Call(v)
		}
	}
	cg.mux.Lock()
	for callback := range callbacks {
		if !callback.Enabled {
			delete(cg.callbacks, callback)
		}
	}
	cg.mux.Unlock()
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
