package kraken

import (
	"fmt"
	"sync"
	"time"
)

// EpochCounter is the default counter for the module, used for nonce generation.
type EpochCounter struct {
	Granularity time.Duration
	unix        int64
	counter     int
	counterMux  sync.Mutex
}

// NewEpochCounter constructs a NonceGenerator.
func NewEpochCounter() *EpochCounter {
	return &EpochCounter{
		Granularity: time.Second,
		counter:     -1,
	}
}

// Get concatenates the unix epoch value and 3 leading zero counter values.
func (c *EpochCounter) Get() string {
	c.counterMux.Lock()
	defer c.counterMux.Unlock()
	currentTime := time.Now()
	var currentUnix int64
	switch c.Granularity {
	case time.Millisecond:
		currentUnix = currentTime.UnixMilli()
	case time.Microsecond:
		currentUnix = currentTime.UnixMicro()
	case time.Nanosecond:
		currentUnix = currentTime.UnixNano()
	default:
		c.Granularity = time.Second
		currentUnix = currentTime.Unix()
	}
	if currentUnix != c.unix {
		c.counter = -1
		c.unix = currentUnix
	}
	if c.counter >= 999 {
		time.Sleep(time.Until(currentTime.Add(c.Granularity)))
		return c.Get()
	}
	c.counter += 1
	nonce := fmt.Sprintf("%d%03d", currentUnix, c.counter)
	return nonce
}
