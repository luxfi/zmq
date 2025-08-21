// Copyright 2025 The go-zeromq Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zmq4

import (
	"fmt"
)

// Proxy starts a proxy that forwards messages between frontend and backend.
// This is the ONLY way to create a proxy - no complex variations.
func Proxy(frontend, backend Socket) error {
	if frontend == nil || backend == nil {
		return fmt.Errorf("frontend and backend sockets are required")
	}

	errChan := make(chan error, 2)

	// Frontend to backend
	go func() {
		for {
			msg, err := frontend.Recv()
			if err != nil {
				errChan <- err
				return
			}
			if err := backend.Send(msg); err != nil {
				errChan <- err
				return
			}
		}
	}()

	// Backend to frontend
	go func() {
		for {
			msg, err := backend.Recv()
			if err != nil {
				errChan <- err
				return
			}
			if err := frontend.Send(msg); err != nil {
				errChan <- err
				return
			}
		}
	}()

	// Wait for first error
	return <-errChan
}
