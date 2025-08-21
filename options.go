// Copyright 2025 The go-zeromq Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zmq4

import (
	"time"
)

// Option configures a ZeroMQ socket - ONE way to configure
type Option func(s *socket)

// WithID sets socket identity
func WithID(id SocketIdentity) Option {
	return func(s *socket) {
		s.id = id
	}
}

// WithSecurity sets security mechanism (default: NULL)
func WithSecurity(sec Security) Option {
	return func(s *socket) {
		if sec == nil {
			sec = nullSecurity{}
		}
		s.sec = sec
	}
}

// WithTimeout sets socket timeout
func WithTimeout(timeout time.Duration) Option {
	return func(s *socket) {
		s.timeout = timeout
	}
}

// WithLogger is a no-op for compatibility
func WithLogger(logger interface{}) Option {
	return func(s *socket) {}
}

// WithDialerRetry is a no-op for compatibility
func WithDialerRetry(retry time.Duration) Option {
	return func(s *socket) {
		s.retry = retry
	}
}

// WithDialerMaxRetries is a no-op for compatibility
func WithDialerMaxRetries(maxRetries int) Option {
	return func(s *socket) {
		s.maxRetries = maxRetries
	}
}

// WithAutomaticReconnect is a no-op for compatibility
func WithAutomaticReconnect(auto bool) Option {
	return func(s *socket) {
		s.autoReconnect = auto
	}
}

// Socket option constants - only essential ones
const (
	OptionSubscribe   = "SUBSCRIBE"
	OptionUnsubscribe = "UNSUBSCRIBE"
	OptionHWM         = "HWM"
	OptionIdentity    = "IDENTITY"
)
