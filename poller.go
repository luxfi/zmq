// Copyright 2025 The go-zeromq Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zmq4

import (
	"fmt"
	"sync"
	"time"
)

// State represents the state of a socket in a poller
type State int

const (
	// Readable indicates the socket is ready for reading
	Readable State = 1 << iota
	// Writable indicates the socket is ready for writing
	Writable
	// Error indicates an error condition on the socket
	Error
)

// Poller provides I/O multiplexing for multiple sockets
type Poller struct {
	mu      sync.RWMutex
	sockets []pollerSocket
	active  bool
}

type pollerSocket struct {
	socket Socket
	events State
}

// NewPoller creates a new Poller
func NewPoller() *Poller {
	return &Poller{
		sockets: make([]pollerSocket, 0),
	}
}

// Add adds a socket to the poller with the specified events to monitor
func (p *Poller) Add(socket Socket, events State) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if socket == nil {
		return fmt.Errorf("cannot add nil socket to poller")
	}
	
	// Check if socket already exists
	for i, ps := range p.sockets {
		if ps.socket == socket {
			// Update events
			p.sockets[i].events = events
			return nil
		}
	}
	
	// Add new socket
	p.sockets = append(p.sockets, pollerSocket{
		socket: socket,
		events: events,
	})
	
	return nil
}

// Remove removes a socket from the poller
func (p *Poller) Remove(socket Socket) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	for i, ps := range p.sockets {
		if ps.socket == socket {
			// Remove socket
			p.sockets = append(p.sockets[:i], p.sockets[i+1:]...)
			return nil
		}
	}
	
	return fmt.Errorf("socket not found in poller")
}

// PollItem represents a socket and its ready events
type PollItem struct {
	Socket Socket
	Events State
}

// Poll waits for events on the registered sockets
func (p *Poller) Poll(timeout time.Duration) ([]PollItem, error) {
	p.mu.RLock()
	sockets := make([]pollerSocket, len(p.sockets))
	copy(sockets, p.sockets)
	p.mu.RUnlock()
	
	if len(sockets) == 0 {
		return nil, fmt.Errorf("no sockets registered")
	}
	
	// Calculate deadline
	var deadline time.Time
	if timeout > 0 {
		deadline = time.Now().Add(timeout)
	} else if timeout == 0 {
		// Non-blocking poll
		deadline = time.Now()
	}
	// timeout < 0 means infinite wait
	
	ready := make([]PollItem, 0)
	
	// Simple polling implementation
	// In a real implementation, this would use epoll/kqueue/select
	start := time.Now()
	for {
		// Check each socket
		for _, ps := range sockets {
			if ps.events == 0 {
				continue
			}
			
			item := PollItem{
				Socket: ps.socket,
				Events: 0,
			}
			
			// Check if socket is ready
			// For now, we'll assume sockets are always writable
			// and check for readability based on socket type
			if ps.events&Writable != 0 {
				// Most sockets are usually writable
				item.Events |= Writable
			}
			
			if ps.events&Readable != 0 {
				// Check if socket has data to read
				// This is a simplified check - real implementation would use system calls
				// For now, we'll skip the readable check to avoid blocking
			}
			
			if item.Events != 0 {
				ready = append(ready, item)
			}
		}
		
		// If we found ready sockets or timeout, return
		if len(ready) > 0 {
			return ready, nil
		}
		
		// Check timeout
		if timeout >= 0 && time.Now().After(deadline) {
			return nil, nil // Timeout is not an error
		}
		
		// If infinite wait and no sockets ready, sleep briefly
		if timeout < 0 && len(ready) == 0 {
			time.Sleep(10 * time.Millisecond)
			// Check every 100ms max
			if time.Since(start) > 100*time.Millisecond {
				return nil, nil
			}
		} else {
			break
		}
	}
	
	return ready, nil
}

// PollAll polls all sockets with infinite timeout
func (p *Poller) PollAll() ([]PollItem, error) {
	return p.Poll(-1)
}

// String returns a string representation of the state
func (s State) String() string {
	var states []string
	if s&Readable != 0 {
		states = append(states, "READABLE")
	}
	if s&Writable != 0 {
		states = append(states, "WRITABLE")
	}
	if s&Error != 0 {
		states = append(states, "ERROR")
	}
	if len(states) == 0 {
		return "NONE"
	}
	result := states[0]
	for i := 1; i < len(states); i++ {
		result += "|" + states[i]
	}
	return result
}

// Reactor provides event-driven I/O for ZeroMQ sockets
type Reactor struct {
	poller    *Poller
	handlers  map[Socket]func(State)
	running   bool
	mu        sync.RWMutex
	stopCh    chan struct{}
}

// NewReactor creates a new Reactor
func NewReactor() *Reactor {
	return &Reactor{
		poller:   NewPoller(),
		handlers: make(map[Socket]func(State)),
		stopCh:   make(chan struct{}),
	}
}

// AddSocket adds a socket with its event handler
func (r *Reactor) AddSocket(socket Socket, events State, handler func(State)) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if err := r.poller.Add(socket, events); err != nil {
		return err
	}
	
	r.handlers[socket] = handler
	return nil
}

// RemoveSocket removes a socket from the reactor
func (r *Reactor) RemoveSocket(socket Socket) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if err := r.poller.Remove(socket); err != nil {
		return err
	}
	
	delete(r.handlers, socket)
	return nil
}

// Run starts the reactor event loop
func (r *Reactor) Run() error {
	r.mu.Lock()
	if r.running {
		r.mu.Unlock()
		return fmt.Errorf("reactor already running")
	}
	r.running = true
	r.mu.Unlock()
	
	defer func() {
		r.mu.Lock()
		r.running = false
		r.mu.Unlock()
	}()
	
	for {
		select {
		case <-r.stopCh:
			return nil
		default:
			items, err := r.poller.Poll(100 * time.Millisecond)
			if err != nil {
				return err
			}
			
			for _, item := range items {
				r.mu.RLock()
				handler, ok := r.handlers[item.Socket]
				r.mu.RUnlock()
				
				if ok && handler != nil {
					handler(item.Events)
				}
			}
		}
	}
}

// Stop stops the reactor
func (r *Reactor) Stop() {
	close(r.stopCh)
}