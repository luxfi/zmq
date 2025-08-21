// Copyright 2025 The go-zeromq Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zmq4

import (
	"context"
	"net"
)

// NewStream returns a new STREAM ZeroMQ socket.
// The returned socket value is initially unbound.
// STREAM sockets are used to send and receive TCP data
// from a non-ZeroMQ peer when using the tcp:// transport.
func NewStream(ctx context.Context, opts ...Option) Socket {
	stream := &streamSocket{sck: newSocket(ctx, Stream, opts...)}
	return stream
}

// streamSocket is a STREAM ZeroMQ socket.
type streamSocket struct {
	sck *socket
}

// Close closes the open Socket
func (stream *streamSocket) Close() error {
	return stream.sck.Close()
}

// Send puts the message on the outbound send queue.
// Send blocks until the message can be queued or the send deadline expires.
func (stream *streamSocket) Send(msg Msg) error {
	return stream.sck.Send(msg)
}

// SendMulti puts the message on the outbound send queue.
// SendMulti blocks until the message can be queued or the send deadline expires.
// The message will be sent as a multipart message.
func (stream *streamSocket) SendMulti(msg Msg) error {
	return stream.sck.SendMulti(msg)
}

// Recv receives a complete message.
func (stream *streamSocket) Recv() (Msg, error) {
	return stream.sck.Recv()
}

// Listen connects a local endpoint to the Socket.
func (stream *streamSocket) Listen(ep string) error {
	return stream.sck.Listen(ep)
}

// Dial connects a remote endpoint to the Socket.
func (stream *streamSocket) Dial(ep string) error {
	return stream.sck.Dial(ep)
}

// Type returns the type of this Socket (PUB, SUB, ...)
func (stream *streamSocket) Type() SocketType {
	return stream.sck.Type()
}

// Addr returns the listener's address.
// Addr returns nil if the socket isn't a listener.
func (stream *streamSocket) Addr() net.Addr {
	return stream.sck.Addr()
}

// GetOption is used to retrieve an option for a socket.
func (stream *streamSocket) GetOption(name string) (interface{}, error) {
	return stream.sck.GetOption(name)
}

// SetOption is used to set an option for a socket.
func (stream *streamSocket) SetOption(name string, value interface{}) error {
	return stream.sck.SetOption(name, value)
}

var (
	_ Socket = (*streamSocket)(nil)
)
