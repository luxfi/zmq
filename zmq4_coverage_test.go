// Copyright 2025 The go-zeromq Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zmq4_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/luxfi/zmq/v4"
)

func TestBackendName(t *testing.T) {
	backend := zmq4.BackendName()
	if backend == "" {
		t.Fatal("BackendName returned empty string")
	}
	
	// Should be either "pure-go" or "czmq"
	if backend != "pure-go" && backend != "czmq" {
		t.Fatalf("unexpected backend name: %s", backend)
	}
	
	t.Logf("Using backend: %s", backend)
}

func TestSocketOptions(t *testing.T) {
	ctx := context.Background()
	
	tests := []struct {
		name   string
		socket zmq4.Socket
	}{
		{"REQ", zmq4.NewReq(ctx)},
		{"REP", zmq4.NewRep(ctx)},
		{"PUB", zmq4.NewPub(ctx)},
		{"SUB", zmq4.NewSub(ctx)},
		{"PUSH", zmq4.NewPush(ctx)},
		{"PULL", zmq4.NewPull(ctx)},
		{"PAIR", zmq4.NewPair(ctx)},
		{"DEALER", zmq4.NewDealer(ctx)},
		{"ROUTER", zmq4.NewRouter(ctx)},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.socket.Close()
			
			// Test setting and getting options
			err := tt.socket.SetOption(zmq4.OptionHWM, 100)
			if err != nil {
				// Option might not be supported
				t.Logf("SetOption HWM not supported: %v", err)
			}
			
			val, err := tt.socket.GetOption(zmq4.OptionHWM)
			if err != nil {
				t.Logf("GetOption HWM not supported: %v", err)
			} else if val != nil {
				t.Logf("HWM value: %v", val)
			}
			
			// Test socket type
			sockType := tt.socket.Type()
			if sockType == "" {
				t.Errorf("socket.Type() returned empty string")
			}
		})
	}
}

func TestMessageClone(t *testing.T) {
	original := zmq4.NewMsgFrom([]byte("frame1"), []byte("frame2"))
	
	cloned := original.Clone()
	
	// Verify frames are copied
	if len(cloned.Frames) != len(original.Frames) {
		t.Fatalf("Clone: frame count mismatch")
	}
	
	for i := range original.Frames {
		if string(cloned.Frames[i]) != string(original.Frames[i]) {
			t.Errorf("Clone: frame %d content mismatch", i)
		}
		
		// Modify cloned frame to ensure it's a copy
		cloned.Frames[i][0] = 'X'
		if original.Frames[i][0] == 'X' {
			t.Error("Clone: frames are not independent copies")
		}
	}
}

func TestMultipartMessage(t *testing.T) {
	ctx := context.Background()
	
	// Create PAIR sockets for testing
	s1 := zmq4.NewPair(ctx)
	defer s1.Close()
	s2 := zmq4.NewPair(ctx)
	defer s2.Close()
	
	// Bind and connect
	err := s1.Listen("tcp://127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	
	addr := s1.Addr()
	err = s2.Dial(fmt.Sprintf("tcp://%s", addr))
	if err != nil {
		t.Fatal(err)
	}
	
	// Allow connection to establish
	time.Sleep(100 * time.Millisecond)
	
	// Send multipart message
	msg := zmq4.NewMsgFrom([]byte("part1"), []byte("part2"), []byte("part3"))
	err = s1.SendMulti(msg)
	if err != nil {
		t.Fatal("SendMulti:", err)
	}
	
	// Receive multipart message
	received, err := s2.Recv()
	if err != nil {
		t.Fatal("Recv:", err)
	}
	
	// Verify all parts received
	if len(received.Frames) != 3 {
		t.Fatalf("Expected 3 frames, got %d", len(received.Frames))
	}
	
	for i, expected := range []string{"part1", "part2", "part3"} {
		if string(received.Frames[i]) != expected {
			t.Errorf("Frame %d: got %q, want %q", i, received.Frames[i], expected)
		}
	}
}

func TestAuthFunctions(t *testing.T) {
	// Test key generation
	public, secret, err := zmq4.NewCurveKeypair()
	if err != nil {
		t.Fatal("NewCurveKeypair:", err)
	}
	
	if public == "" || secret == "" {
		t.Fatal("Empty keys generated")
	}
	
	// Test Z85 encoding/decoding
	data := []byte("test data for encoding")
	encoded := zmq4.Z85encode(data)
	if encoded == "" {
		t.Fatal("Z85encode returned empty string")
	}
	
	decoded, err := zmq4.Z85decode(encoded)
	if err != nil {
		t.Fatal("Z85decode:", err)
	}
	
	if string(decoded) != string(data) {
		t.Errorf("Z85 round-trip failed: got %q, want %q", decoded, data)
	}
	
	// Test auth functions
	err = zmq4.AuthStart()
	if err != nil {
		t.Log("AuthStart:", err)
	} else {
		defer zmq4.AuthStop()
		
		// Test allow/deny
		zmq4.AuthAllow("test-domain", "127.0.0.1")
		zmq4.AuthDeny("test-domain", "192.168.1.1")
		
		// Test CURVE operations
		zmq4.AuthCurveAdd("test-domain", public)
		zmq4.AuthCurveRemove("test-domain", public)
		
		// Test verbose setting
		zmq4.AuthSetVerbose(true)
		zmq4.AuthSetVerbose(false)
	}
}

func TestXPubXSub(t *testing.T) {
	ctx := context.Background()
	
	xpub := zmq4.NewXPub(ctx)
	defer xpub.Close()
	xsub := zmq4.NewXSub(ctx)
	defer xsub.Close()
	
	// Test binding
	err := xpub.Listen("tcp://127.0.0.1:0")
	if err != nil {
		t.Fatal("xpub.Listen:", err)
	}
	
	addr := xpub.Addr()
	err = xsub.Dial(fmt.Sprintf("tcp://%s", addr))
	if err != nil {
		t.Fatal("xsub.Dial:", err)
	}
	
	// XPUB and XSUB are specialized PUB/SUB sockets
	if xpub.Type() != "XPUB" {
		t.Errorf("xpub.Type(): got %q, want %q", xpub.Type(), "XPUB")
	}
	
	if xsub.Type() != "XSUB" {
		t.Errorf("xsub.Type(): got %q, want %q", xsub.Type(), "XSUB")
	}
}

func TestInprocTransport(t *testing.T) {
	ctx := context.Background()
	
	// Create PUSH/PULL pair for inproc testing
	push := zmq4.NewPush(ctx)
	defer push.Close()
	pull := zmq4.NewPull(ctx)
	defer pull.Close()
	
	// Use inproc transport
	endpoint := "inproc://test-endpoint"
	
	err := pull.Listen(endpoint)
	if err != nil {
		t.Fatal("pull.Listen:", err)
	}
	
	err = push.Dial(endpoint)
	if err != nil {
		t.Fatal("push.Dial:", err)
	}
	
	// Send message
	msg := zmq4.NewMsg([]byte("inproc test"))
	err = push.Send(msg)
	if err != nil {
		t.Fatal("push.Send:", err)
	}
	
	// Receive message
	received, err := pull.Recv()
	if err != nil {
		t.Fatal("pull.Recv:", err)
	}
	
	if string(received.Frames[0]) != "inproc test" {
		t.Errorf("Got %q, want %q", received.Frames[0], "inproc test")
	}
}

func TestErrorConditions(t *testing.T) {
	ctx := context.Background()
	
	// Test invalid endpoint
	socket := zmq4.NewReq(ctx)
	defer socket.Close()
	
	err := socket.Dial("invalid://endpoint")
	if err == nil {
		t.Error("Expected error for invalid endpoint")
	}
	
	// Test double close
	socket.Close()
	err = socket.Close()
	if err == nil {
		t.Log("Double close allowed (might be acceptable)")
	}
	
	// Test send after close
	msg := zmq4.NewMsg([]byte("test"))
	err = socket.Send(msg)
	if err == nil {
		t.Error("Expected error for send after close")
	}
}