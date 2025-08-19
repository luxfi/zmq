// Copyright 2025 The go-zeromq Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zmq4

import (
	"context"
	"fmt"
	"sync"
)

// Proxy starts a proxy in the current goroutine.
// It connects a frontend socket to a backend socket.
// Messages received on the frontend are sent to the backend, and vice versa.
// If capture socket is not nil, all messages are also sent to the capture socket.
func Proxy(frontend, backend, capture Socket) error {
	if frontend == nil || backend == nil {
		return fmt.Errorf("frontend and backend sockets are required")
	}
	
	// Create channels for bidirectional message passing
	frontToBack := make(chan Msg, 100)
	backToFront := make(chan Msg, 100)
	errChan := make(chan error, 2)
	done := make(chan struct{})
	
	var wg sync.WaitGroup
	wg.Add(2)
	
	// Frontend to backend forwarding
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				msg, err := frontend.Recv()
				if err != nil {
					select {
					case errChan <- err:
					case <-done:
					}
					return
				}
				
				// Send to capture if provided
				if capture != nil {
					capture.Send(msg.Clone())
				}
				
				select {
				case frontToBack <- msg:
				case <-done:
					return
				}
			}
		}
	}()
	
	// Backend to frontend forwarding
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				msg, err := backend.Recv()
				if err != nil {
					select {
					case errChan <- err:
					case <-done:
					}
					return
				}
				
				// Send to capture if provided
				if capture != nil {
					capture.Send(msg.Clone())
				}
				
				select {
				case backToFront <- msg:
				case <-done:
					return
				}
			}
		}
	}()
	
	// Message forwarding loops
	go func() {
		for msg := range frontToBack {
			if err := backend.Send(msg); err != nil {
				select {
				case errChan <- err:
				case <-done:
				}
				return
			}
		}
	}()
	
	go func() {
		for msg := range backToFront {
			if err := frontend.Send(msg); err != nil {
				select {
				case errChan <- err:
				case <-done:
				}
				return
			}
		}
	}()
	
	// Wait for error or completion
	err := <-errChan
	close(done)
	wg.Wait()
	close(frontToBack)
	close(backToFront)
	
	return err
}


