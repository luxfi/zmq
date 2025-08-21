//go:build cgo && czmq
// +build cgo,czmq

// Copyright 2025 The go-zeromq Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zmq4

// #cgo LDFLAGS: -lczmq -lzmq
// #include <czmq.h>
import "C"

// When CGO is enabled and czmq tag is set, this file provides the optimized implementation

// BackendName returns the name of the current backend
func BackendName() string {
	return "czmq"
}

// IsCZMQAvailable returns true when CZMQ is available
func IsCZMQAvailable() bool {
	return true
}
