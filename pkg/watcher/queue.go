package watcher

import (
	"sync"
)

// ItemQueue the queue of Items
type PacketDurationQueue struct {
	items []float32
	lock  sync.RWMutex
}

// New creates a new ItemQueue
func NewItemQueue() *PacketDurationQueue {
	return &PacketDurationQueue{items: make([]float32, MinimumSpeedBlock)}
}

// Enqueue adds an Item to the end of the queue
func (s *PacketDurationQueue) Enqueue(t float32) {
	s.lock.Lock()
	s.items = append(s.items, t)
	s.lock.Unlock()
}

// Dequeue removes an Item from the start of the queue
func (s *PacketDurationQueue) Dequeue() float32 {
	s.lock.Lock()
	item := s.items[0]
	s.items = s.items[1:len(s.items)]
	s.lock.Unlock()
	return item
}

// Front returns the item next in the queue, without removing it
func (s *PacketDurationQueue) Front() float32 {
	s.lock.RLock()
	item := s.items[0]
	s.lock.RUnlock()
	return item
}

// IsEmpty returns true if the queue is empty
func (s *PacketDurationQueue) IsEmpty() bool {
	return len(s.items) == 0
}

// Size returns the number of Items in the queue
func (s *PacketDurationQueue) Size() int {
	return len(s.items)
}
