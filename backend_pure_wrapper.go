//go:build !cgo || !czmq
// +build !cgo !czmq

// Copyright 2025 The go-zeromq Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zmq4

// When CGO is disabled or czmq tag is not set, this file provides the pure Go implementation

// BackendName returns the name of the current backend
func BackendName() string {
	return "pure-go"
}

// IsCZMQAvailable returns false when using pure Go backend
func IsCZMQAvailable() bool {
	return false
}
