package internal

import (
	"sync"
)

type StationCheckCounter struct {
	value int
	mu    sync.Mutex
}

func (c *StationCheckCounter) Increment() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

func (c *StationCheckCounter) Value() int {
	return c.value
}
