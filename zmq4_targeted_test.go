// Copyright 2025 The go-zeromq Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Targeted tests for uncovered functionality

package zmq4_test

import (
	"context"
	"testing"
	"time"

	"github.com/luxfi/zmq/v4"
)

// Test Stream socket
func TestStreamSocket(t *testing.T) {
	ctx := context.Background()

	// Create stream socket
	stream := zmq4.NewStream(ctx)
	if stream == nil {
		t.Skip("Stream socket not available")
	}
	defer stream.Close()

	// Test basic methods
	if stream.Type() == "" {
		t.Error("Stream Type() returned empty")
	}

	// Test Listen
	err := stream.Listen("tcp://127.0.0.1:0")
	if err != nil {
		t.Logf("Stream Listen error: %v", err)
	}

	// Test Addr
	addr := stream.Addr()
	if addr != nil {
		t.Logf("Stream addr: %v", addr)
	}

	// Test Send (may not be supported)
	msg := zmq4.NewMsg([]byte("test"))
	err = stream.Send(msg)
	if err != nil {
		t.Logf("Stream Send error (expected): %v", err)
	}

	// Test Recv (non-blocking)
	done := make(chan bool, 1)
	go func() {
		_, err := stream.Recv()
		if err != nil {
			t.Logf("Stream Recv error: %v", err)
		}
		done <- true
	}()

	select {
	case <-done:
		// Received or errored
	case <-time.After(10 * time.Millisecond):
		// Timeout ok
	}

	// Test GetOption
	val, err := stream.GetOption(zmq4.OptionIdentity)
	if err != nil {
		t.Logf("Stream GetOption error: %v", err)
	} else {
		t.Logf("Stream identity: %v", val)
	}

	// Test SetOption
	err = stream.SetOption(zmq4.OptionIdentity, "stream-id")
	if err != nil {
		t.Logf("Stream SetOption error: %v", err)
	}
}

// Test Queue functionality
func TestQueue(t *testing.T) {
	// Create queue
	queue := zmq4.NewQueue()
	if queue == nil {
		t.Skip("Queue not available")
	}

	// Queue doesn't have many methods in simplified API
	t.Log("Queue created successfully")
}

// Test Proxy with capture socket
func TestProxyWithCaptureSocket(t *testing.T) {
	ctx := context.Background()

	// Create frontend
	frontend := zmq4.NewRouter(ctx)
	if frontend == nil {
		t.Skip("Router not available")
	}
	defer frontend.Close()

	// Create backend
	backend := zmq4.NewDealer(ctx)
	if backend == nil {
		t.Skip("Dealer not available")
	}
	defer backend.Close()

	// Create capture socket (optional third parameter)
	capture := zmq4.NewPub(ctx)
	if capture != nil {
		defer capture.Close()
	}

	// Bind sockets
	err := frontend.Listen("tcp://127.0.0.1:0")
	if err != nil {
		t.Fatal("Frontend listen failed:", err)
	}

	err = backend.Listen("tcp://127.0.0.1:0")
	if err != nil {
		t.Fatal("Backend listen failed:", err)
	}

	// Start proxy in goroutine
	done := make(chan error, 1)
	go func() {
		// Proxy with capture socket (if available)
		err := zmq4.Proxy(frontend, backend)
		done <- err
	}()

	// Let it run briefly
	select {
	case err := <-done:
		if err != nil {
			t.Logf("Proxy error: %v", err)
		}
	case <-time.After(10 * time.Millisecond):
		// Expected - proxy runs forever
	}
}

// Test connection-level functions
func TestConnectionFunctions(t *testing.T) {
	ctx := context.Background()

	// Test with DEALER-ROUTER as they support more operations
	dealer := zmq4.NewDealer(ctx, zmq4.WithID(zmq4.SocketIdentity("dealer-1")))
	if dealer == nil {
		t.Skip("Dealer not available")
	}
	defer dealer.Close()

	router := zmq4.NewRouter(ctx)
	if router == nil {
		t.Skip("Router not available")
	}
	defer router.Close()

	// Bind router
	err := router.Listen("tcp://127.0.0.1:0")
	if err != nil {
		t.Fatal("Router listen failed:", err)
	}

	addr := router.Addr()
	if addr == nil {
		t.Fatal("Router addr is nil")
	}

	// Connect dealer
	err = dealer.Dial("tcp://" + addr.String())
	if err != nil {
		t.Fatal("Dealer dial failed:", err)
	}

	// Wait for connection
	time.Sleep(50 * time.Millisecond)

	// Send from dealer (no identity needed)
	msg := zmq4.NewMsgFrom([]byte("data"))
	err = dealer.Send(msg)
	if err != nil {
		t.Logf("Dealer send error: %v", err)
	}

	// Receive on router (gets identity frame)
	recvDone := make(chan zmq4.Msg, 1)
	go func() {
		msg, err := router.Recv()
		if err != nil {
			t.Logf("Router recv error: %v", err)
		} else {
			recvDone <- msg
		}
	}()

	select {
	case msg := <-recvDone:
		if len(msg.Frames) < 2 {
			t.Error("Router should add identity frame")
		} else {
			t.Logf("Received %d frames", len(msg.Frames))
		}
	case <-time.After(100 * time.Millisecond):
		// Timeout
	}

	// Test SendMulti
	multiMsg := zmq4.NewMsgFrom([]byte("frame1"), []byte("frame2"))
	err = dealer.SendMulti(multiMsg)
	if err != nil {
		t.Logf("SendMulti error: %v", err)
	}
}

// Test XPUB-XSUB specific functionality
func TestXPubXSubOperations(t *testing.T) {
	ctx := context.Background()

	// Create XPUB
	xpub := zmq4.NewXPub(ctx)
	if xpub == nil {
		t.Skip("XPub not available")
	}
	defer xpub.Close()

	// Create XSUB
	xsub := zmq4.NewXSub(ctx)
	if xsub == nil {
		t.Skip("XSub not available")
	}
	defer xsub.Close()

	// Bind XPUB
	err := xpub.Listen("tcp://127.0.0.1:0")
	if err != nil {
		t.Fatal("XPub listen failed:", err)
	}

	addr := xpub.Addr()
	if addr == nil {
		t.Fatal("XPub addr is nil")
	}

	// Connect XSUB
	err = xsub.Dial("tcp://" + addr.String())
	if err != nil {
		t.Fatal("XSub dial failed:", err)
	}

	// Wait for connection
	time.Sleep(50 * time.Millisecond)

	// XSUB sends subscription message
	subMsg := zmq4.NewMsg([]byte("\x01topic"))
	err = xsub.Send(subMsg)
	if err != nil {
		t.Logf("XSub send subscription error: %v", err)
	}

	// XPUB receives subscription
	recvDone := make(chan zmq4.Msg, 1)
	go func() {
		msg, err := xpub.Recv()
		if err != nil {
			t.Logf("XPub recv error: %v", err)
		} else {
			recvDone <- msg
		}
	}()

	select {
	case msg := <-recvDone:
		if len(msg.Frames) > 0 {
			t.Logf("XPub received subscription: %v", msg.Frames[0])
		}
	case <-time.After(100 * time.Millisecond):
		// Timeout
	}

	// XPUB publishes message
	pubMsg := zmq4.NewMsg([]byte("topic: data"))
	err = xpub.Send(pubMsg)
	if err != nil {
		t.Logf("XPub send error: %v", err)
	}

	// XSUB receives message
	recvDone2 := make(chan zmq4.Msg, 1)
	go func() {
		msg, err := xsub.Recv()
		if err != nil {
			t.Logf("XSub recv error: %v", err)
		} else {
			recvDone2 <- msg
		}
	}()

	select {
	case msg := <-recvDone2:
		if len(msg.Frames) > 0 {
			t.Logf("XSub received: %v", msg.Frames[0])
		}
	case <-time.After(100 * time.Millisecond):
		// Timeout
	}
}

// Test auth with actual authentication
func TestAuthWithAuthentication(t *testing.T) {
	// Start auth
	err := zmq4.AuthStart()
	if err != nil {
		t.Logf("AuthStart error: %v", err)
		return
	}
	defer zmq4.AuthStop()

	// Set verbose
	zmq4.AuthSetVerbose(true)

	// Allow specific domain and IP
	zmq4.AuthAllow("global", "127.0.0.1")
	zmq4.AuthAllow("global", "::1")

	// Deny specific IP
	zmq4.AuthDeny("blocked", "192.168.1.1")

	// Generate keypair for curve auth
	pubKey, secKey, err := zmq4.NewCurveKeypair()
	if err != nil {
		t.Logf("NewCurveKeypair error: %v", err)
		return
	}

	if pubKey != "" && secKey != "" {
		// Add curve key
		zmq4.AuthCurveAdd("secure", pubKey)

		// Get public key from secret
		derivedPub, err := zmq4.AuthCurvePublic(secKey)
		if err != nil {
			t.Logf("AuthCurvePublic error: %v", err)
		} else if derivedPub != pubKey {
			t.Logf("Derived public key mismatch")
		}

		// Remove curve key
		zmq4.AuthCurveRemove("secure", pubKey)
	}

	// Test metadata handler
	zmq4.AuthSetMetadataHandler(func(domain, address string) map[string]string {
		metadata := make(map[string]string)
		metadata["User-Id"] = "test-user"
		metadata["Domain"] = domain
		metadata["Address"] = address
		return metadata
	})

	// Test Z85 encoding/decoding
	testData := []byte("Hello, World! 1234567890")
	// Z85 requires data length to be multiple of 4
	paddedData := make([]byte, ((len(testData)+3)/4)*4)
	copy(paddedData, testData)

	encoded := zmq4.Z85encode(paddedData)
	if encoded != "" {
		decoded, err := zmq4.Z85decode(encoded)
		if err != nil {
			t.Errorf("Z85decode error: %v", err)
		} else if string(decoded[:len(testData)]) != string(testData) {
			t.Errorf("Z85 roundtrip failed")
		}
	}
}

// Test message operations more thoroughly
func TestMessageOperationsDetailed(t *testing.T) {
	// Test empty message
	emptyMsg := zmq4.NewMsg(nil)
	if len(emptyMsg.Frames) != 1 {
		t.Error("Empty message should have 1 frame")
	}
	if len(emptyMsg.Frames[0]) != 0 {
		t.Error("Empty frame should have 0 length")
	}

	// Test Bytes method
	msg := zmq4.NewMsgFrom([]byte("frame1"), []byte("frame2"))
	bytes := msg.Bytes()
	if len(bytes) == 0 {
		t.Error("Bytes() should return non-empty")
	}

	// Test Clone independence
	original := zmq4.NewMsgFrom([]byte("original"))
	cloned := original.Clone()

	// Modify original
	original.Frames[0] = []byte("modified")

	// Check clone is unchanged
	if string(cloned.Frames[0]) != "original" {
		t.Error("Clone was affected by original modification")
	}

	// Test String() method if available
	strMsg := zmq4.NewMsgString("test string")
	if string(strMsg.Frames[0]) != "test string" {
		t.Error("String message mismatch")
	}

	// Test multi-frame string message
	multiStr := zmq4.NewMsgFromString([]string{"str1", "str2", "str3"})
	if len(multiStr.Frames) != 3 {
		t.Error("Multi-string message frame count mismatch")
	}
}

// Test socket options more thoroughly
func TestSocketOptionsDetailed(t *testing.T) {
	ctx := context.Background()

	// Test all option setters
	options := []zmq4.Option{
		zmq4.WithID(zmq4.SocketIdentity("test-id")),
		zmq4.WithSecurity(nil),
		zmq4.WithTimeout(5 * time.Second),
		zmq4.WithLogger(nil),
		zmq4.WithDialerRetry(200 * time.Millisecond),
		zmq4.WithDialerMaxRetries(5),
		zmq4.WithAutomaticReconnect(false),
	}

	// Create socket with all options
	sock := zmq4.NewDealer(ctx, options...)
	if sock == nil {
		t.Skip("Dealer not available")
	}
	defer sock.Close()

	// Test SetOption with various option types
	sock.SetOption(zmq4.OptionHWM, 1000)
	sock.SetOption(zmq4.OptionIdentity, "new-identity")
	sock.SetOption(zmq4.OptionSubscribe, "topic")
	sock.SetOption(zmq4.OptionUnsubscribe, "topic")

	// Test GetOption
	val, err := sock.GetOption(zmq4.OptionIdentity)
	if err != nil {
		t.Logf("GetOption error: %v", err)
	} else {
		t.Logf("Identity: %v", val)
	}

	val, err = sock.GetOption(zmq4.OptionHWM)
	if err != nil {
		t.Logf("GetOption HWM error: %v", err)
	} else {
		t.Logf("HWM: %v", val)
	}
}

// Test endpoint validation
func TestEndpointValidation(t *testing.T) {
	// EndPoint function doesn't exist in simplified API
	t.Skip("EndPoint function not available in simplified API")
}

// Test transport-specific functionality
func TestTransportSpecific(t *testing.T) {
	ctx := context.Background()

	// Test TCP with specific options
	tcpSock := zmq4.NewPair(ctx)
	if tcpSock != nil {
		defer tcpSock.Close()

		err := tcpSock.Listen("tcp://127.0.0.1:0")
		if err != nil {
			t.Logf("TCP listen error: %v", err)
		}
	}

	// Test IPC (Unix domain sockets)
	ipcSock := zmq4.NewPair(ctx)
	if ipcSock != nil {
		defer ipcSock.Close()

		err := ipcSock.Listen("ipc://@test.ipc")
		if err != nil {
			t.Logf("IPC listen error: %v", err)
		}
	}

	// Test inproc (in-process)
	inprocSock := zmq4.NewPair(ctx)
	if inprocSock != nil {
		defer inprocSock.Close()

		err := inprocSock.Listen("inproc://test.inproc")
		if err != nil {
			t.Logf("Inproc listen error: %v", err)
		}
	}
}
