// Copyright 2025 The go-zeromq Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zmq4_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/luxfi/zmq/v4"
)

func TestStreamSocket(t *testing.T) {
	// STREAM sockets are used for TCP connections without ZeroMQ framing
	// They're useful for connecting to non-ZeroMQ TCP services
	
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Test basic STREAM socket creation
	stream := zmq4.NewStream(ctx)
	if stream == nil {
		t.Fatal("Failed to create STREAM socket")
	}
	defer stream.Close()

	if got, want := stream.Type(), zmq4.Stream; got != want {
		t.Fatalf("socket type: got %q, want %q", got, want)
	}
}

func TestStreamTCPConnection(t *testing.T) {
	t.Skip("STREAM socket implementation needs work")
	// Test STREAM socket connecting to a regular TCP server
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Start a regular TCP server
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	addr := listener.Addr().String()
	
	// TCP server that echoes data
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		// Echo server
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				return
			}
			if _, err := conn.Write(buf[:n]); err != nil {
				return
			}
		}
	}()

	// Create STREAM socket and connect
	stream := zmq4.NewStream(ctx)
	defer stream.Close()

	endpoint := fmt.Sprintf("tcp://%s", addr)
	if err := stream.Dial(endpoint); err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}

	// STREAM sockets need to handle identity frames
	// First message from STREAM socket includes identity frame
	time.Sleep(100 * time.Millisecond)

	// Send a message
	testData := []byte("Hello from STREAM socket")
	msg := zmq4.NewMsg(testData)
	if err := stream.Send(msg); err != nil {
		t.Fatalf("Failed to send: %v", err)
	}

	// Receive echo
	reply, err := stream.Recv()
	if err != nil {
		t.Fatalf("Failed to receive: %v", err)
	}

	// STREAM sockets may include identity frames
	// Check if we got our data back
	found := false
	for _, frame := range reply.Frames {
		if bytes.Contains(frame, testData) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Echo not received correctly, got frames: %v", reply.Frames)
	}
}

func TestStreamToStream(t *testing.T) {
	t.Skip("STREAM socket implementation needs work")
	// Test two STREAM sockets communicating
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Server STREAM socket
	server := zmq4.NewStream(ctx)
	defer server.Close()

	if err := server.Listen("tcp://127.0.0.1:0"); err != nil {
		t.Fatal(err)
	}

	addr := server.Addr()
	if addr == nil {
		t.Fatal("Failed to get server address")
	}

	// Client STREAM socket
	client := zmq4.NewStream(ctx)
	defer client.Close()

	endpoint := fmt.Sprintf("tcp://%s", addr.String())
	if err := client.Dial(endpoint); err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}

	// Allow connection to establish
	time.Sleep(100 * time.Millisecond)

	// Send from client
	clientMsg := []byte("Hello from client")
	if err := client.Send(zmq4.NewMsg(clientMsg)); err != nil {
		t.Fatalf("Client send failed: %v", err)
	}

	// Receive on server
	msg, err := server.Recv()
	if err != nil {
		t.Fatalf("Server receive failed: %v", err)
	}

	// STREAM sockets include identity frames
	if len(msg.Frames) < 2 {
		t.Fatalf("Expected at least 2 frames (identity + data), got %d", len(msg.Frames))
	}

	// First frame is identity, second should be our data
	identity := msg.Frames[0]
	data := msg.Frames[1]

	if !bytes.Equal(data, clientMsg) {
		t.Errorf("Data mismatch: got %q, want %q", data, clientMsg)
	}

	// Send reply from server (must include identity)
	serverReply := []byte("Hello from server")
	replyMsg := zmq4.NewMsgFrom(identity, serverReply)
	if err := server.Send(replyMsg); err != nil {
		t.Fatalf("Server send failed: %v", err)
	}

	// Receive reply on client
	reply, err := client.Recv()
	if err != nil {
		t.Fatalf("Client receive failed: %v", err)
	}

	// Check reply
	found := false
	for _, frame := range reply.Frames {
		if bytes.Equal(frame, serverReply) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Reply not found in frames: %v", reply.Frames)
	}
}

func TestStreamCompatibility(t *testing.T) {
	// Test IsCompatible for STREAM sockets
	if !zmq4.Stream.IsCompatible(zmq4.Stream) {
		t.Error("STREAM should be compatible with STREAM")
	}

	// STREAM should not be compatible with other socket types
	incompatible := []zmq4.SocketType{
		zmq4.Req, zmq4.Rep, zmq4.Dealer, zmq4.Router,
		zmq4.Pub, zmq4.Sub, zmq4.XPub, zmq4.XSub,
		zmq4.Push, zmq4.Pull, zmq4.Pair,
	}

	for _, typ := range incompatible {
		if zmq4.Stream.IsCompatible(typ) {
			t.Errorf("STREAM should not be compatible with %s", typ)
		}
		if typ.IsCompatible(zmq4.Stream) {
			t.Errorf("%s should not be compatible with STREAM", typ)
		}
	}
}

func TestStreamHTTPServer(t *testing.T) {
	t.Skip("STREAM socket implementation needs work")
	// Test STREAM socket connecting to HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Start a simple HTTP server
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	addr := listener.Addr().String()
	
	// Simple HTTP server
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				
				// Read request (we don't parse it, just read it)
				buf := make([]byte, 1024)
				n, err := c.Read(buf)
				if err != nil && err != io.EOF {
					return
				}
				
				// Check if it looks like HTTP
				if n > 4 && bytes.HasPrefix(buf[:n], []byte("GET ")) {
					// Send simple HTTP response
					response := "HTTP/1.1 200 OK\r\nContent-Length: 12\r\n\r\nHello World!"
					c.Write([]byte(response))
				}
			}(conn)
		}
	}()

	// Connect with STREAM socket
	stream := zmq4.NewStream(ctx)
	defer stream.Close()

	endpoint := fmt.Sprintf("tcp://%s", addr)
	if err := stream.Dial(endpoint); err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}

	// Allow connection to establish
	time.Sleep(100 * time.Millisecond)

	// Send HTTP request
	httpReq := []byte("GET / HTTP/1.1\r\nHost: localhost\r\n\r\n")
	msg := zmq4.NewMsg(httpReq)
	if err := stream.Send(msg); err != nil {
		t.Fatalf("Failed to send HTTP request: %v", err)
	}

	// Receive HTTP response
	reply, err := stream.Recv()
	if err != nil {
		t.Fatalf("Failed to receive HTTP response: %v", err)
	}

	// Check if we got HTTP response
	found := false
	for _, frame := range reply.Frames {
		if bytes.Contains(frame, []byte("HTTP/1.1 200 OK")) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("HTTP response not received, got frames: %v", reply.Frames)
	}
}

func BenchmarkStream(t *testing.B) {
	ctx := context.Background()
	
	// Server
	server := zmq4.NewStream(ctx)
	defer server.Close()
	server.Listen("tcp://127.0.0.1:0")
	
	addr := server.Addr()
	
	// Client
	client := zmq4.NewStream(ctx)
	defer client.Close()
	client.Dial(fmt.Sprintf("tcp://%s", addr.String()))
	
	// Allow connection
	time.Sleep(100 * time.Millisecond)
	
	// Benchmark message exchange
	data := []byte("benchmark data")
	msg := zmq4.NewMsg(data)
	
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		if err := client.Send(msg); err != nil {
			t.Fatal(err)
		}
		if _, err := server.Recv(); err != nil {
			t.Fatal(err)
		}
	}
}