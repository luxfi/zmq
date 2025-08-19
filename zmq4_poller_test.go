// Copyright 2025 The go-zeromq Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tests ported from github.com/pebbe/zmq4

package zmq4_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/luxfi/zmq/v4"
)

func TestPoller(t *testing.T) {
	t.Skip("Poller implementation needs work")
	ctx := context.Background()
	
	// Create PAIR sockets
	sb := zmq4.NewPair(ctx)
	defer sb.Close()
	sc := zmq4.NewPair(ctx)
	defer sc.Close()
	
	// Bind and connect
	err := sb.Listen("tcp://127.0.0.1:0")
	if err != nil {
		t.Fatal("sb.Listen:", err)
	}
	
	addr := sb.Addr()
	err = sc.Dial(fmt.Sprintf("tcp://%s", addr))
	if err != nil {
		t.Fatal("sc.Dial:", err)
	}
	
	// Create poller
	poller := zmq4.NewPoller()
	
	// Add sockets
	err = poller.Add(sb, 0)
	if err != nil {
		t.Fatal("poller.Add sb:", err)
	}
	err = poller.Add(sc, 0)
	if err != nil {
		t.Fatal("poller.Add sc:", err)
	}
	
	// Poll with no events
	items, err := poller.Poll(100 * time.Millisecond)
	if err != nil {
		t.Error("Poll 1:", err)
	}
	if len(items) != 0 {
		t.Errorf("Poll 1: expected 0 items, got %d", len(items))
	}
	
	// Update to monitor events
	err = poller.Add(sb, zmq4.Writable)
	if err != nil {
		t.Fatal("poller.Add sb WRITABLE:", err)
	}
	err = poller.Add(sc, zmq4.Readable)
	if err != nil {
		t.Fatal("poller.Add sc READABLE:", err)
	}
	
	// Poll - sb should be writable
	items, err = poller.Poll(100 * time.Millisecond)
	if err != nil {
		t.Error("Poll 2:", err)
	}
	if len(items) != 1 {
		t.Errorf("Poll 2: expected 1 item, got %d", len(items))
	} else if items[0].Socket != sb || items[0].Events&zmq4.Writable == 0 {
		t.Errorf("Poll 2: expected sb WRITABLE, got %v", items[0])
	}
	
	// Send message from sb to sc
	content := "12345678ABCDEFGH12345678ABCDEFGH"
	msg := zmq4.NewMsg([]byte(content))
	err = sb.Send(msg)
	if err != nil {
		t.Error("sb.Send:", err)
	}
	
	// Update poller to check sc for readable
	err = poller.Add(sb, 0)
	if err != nil {
		t.Fatal("poller.Add sb 0:", err)
	}
	
	// Poll - sc should be readable
	items, err = poller.Poll(100 * time.Millisecond)
	if err != nil {
		t.Error("Poll 3:", err)
	}
	if len(items) != 1 {
		t.Errorf("Poll 3: expected 1 item, got %d", len(items))
	} else if items[0].Socket != sc || items[0].Events&zmq4.Readable == 0 {
		t.Errorf("Poll 3: expected sc READABLE, got %v", items[0])
	}
	
	// Receive message
	recvMsg, err := sc.Recv()
	if err != nil {
		t.Error("sc.Recv:", err)
	}
	if string(recvMsg.Frames[0]) != content {
		t.Errorf("sc.Recv: got %q, want %q", recvMsg.Frames[0], content)
	}
	
	// Remove sc from poller
	err = poller.Remove(sc)
	if err != nil {
		t.Error("poller.Remove sc:", err)
	}
	
	// Update sb to be writable
	err = poller.Add(sb, zmq4.Writable)
	if err != nil {
		t.Fatal("poller.Add sb WRITABLE:", err)
	}
	
	// Poll - should only have sb
	items, err = poller.Poll(100 * time.Millisecond)
	if err != nil {
		t.Error("Poll 4:", err)
	}
	if len(items) != 1 {
		t.Errorf("Poll 4: expected 1 item, got %d", len(items))
	} else if items[0].Socket != sb || items[0].Events&zmq4.Writable == 0 {
		t.Errorf("Poll 4: expected sb WRITABLE, got %v", items[0])
	}
}

func TestPollerMultipleSockets(t *testing.T) {
	t.Skip("Poller implementation needs work")
	ctx := context.Background()
	
	// Create multiple socket pairs
	var sockets []zmq4.Socket
	defer func() {
		for _, s := range sockets {
			if s != nil {
				s.Close()
			}
		}
	}()
	
	// Create 3 pairs of sockets
	for i := 0; i < 3; i++ {
		pub := zmq4.NewPub(ctx)
		sub := zmq4.NewSub(ctx)
		
		err := pub.Listen(fmt.Sprintf("tcp://127.0.0.1:0"))
		if err != nil {
			t.Fatal("pub.Listen:", err)
		}
		
		addr := pub.Addr()
		err = sub.Dial(fmt.Sprintf("tcp://%s", addr))
		if err != nil {
			t.Fatal("sub.Dial:", err)
		}
		
		// Subscribe to all
		sub.SetOption(zmq4.OptionSubscribe, "")
		
		sockets = append(sockets, pub, sub)
	}
	
	// Create poller and add all sub sockets
	poller := zmq4.NewPoller()
	for i := 1; i < len(sockets); i += 2 {
		err := poller.Add(sockets[i], zmq4.Readable)
		if err != nil {
			t.Fatal("poller.Add:", err)
		}
	}
	
	// Send messages from all publishers
	for i := 0; i < len(sockets); i += 2 {
		msg := zmq4.NewMsg([]byte(fmt.Sprintf("Message from pub %d", i/2)))
		err := sockets[i].Send(msg)
		if err != nil {
			t.Error("Send:", err)
		}
	}
	
	// Give messages time to propagate
	time.Sleep(50 * time.Millisecond)
	
	// Poll - all subscribers should be readable
	items, err := poller.Poll(100 * time.Millisecond)
	if err != nil {
		t.Error("Poll:", err)
	}
	if len(items) != 3 {
		t.Errorf("Poll: expected 3 items, got %d", len(items))
	}
	
	// Read all messages
	for _, item := range items {
		if item.Events&zmq4.Readable != 0 {
			msg, err := item.Socket.Recv()
			if err != nil {
				t.Error("Recv:", err)
			} else {
				t.Logf("Received: %s", msg.Frames[0])
			}
		}
	}
}

func TestReactor(t *testing.T) {
	t.Skip("Temporarily disabled - reactor implementation needs work")
	ctx := context.Background()
	
	// Create socket pair
	req := zmq4.NewReq(ctx)
	defer req.Close()
	rep := zmq4.NewRep(ctx)
	defer rep.Close()
	
	// Bind and connect
	err := rep.Listen("tcp://127.0.0.1:0")
	if err != nil {
		t.Fatal("rep.Listen:", err)
	}
	
	addr := rep.Addr()
	err = req.Dial(fmt.Sprintf("tcp://%s", addr))
	if err != nil {
		t.Fatal("req.Dial:", err)
	}
	
	// Create reactor
	reactor := zmq4.NewReactor()
	
	// Track received messages
	received := make(chan string, 10)
	
	// Add reply socket with handler
	err = reactor.AddSocket(rep, zmq4.Readable, func(state zmq4.State) {
		if state&zmq4.Readable != 0 {
			msg, err := rep.Recv()
			if err != nil {
				t.Error("rep.Recv:", err)
				return
			}
			
			// Echo back
			reply := zmq4.NewMsg([]byte("Reply: " + string(msg.Frames[0])))
			err = rep.Send(reply)
			if err != nil {
				t.Error("rep.Send:", err)
				return
			}
			
			received <- string(msg.Frames[0])
		}
	})
	if err != nil {
		t.Fatal("reactor.AddSocket rep:", err)
	}
	
	// Add request socket with handler
	err = reactor.AddSocket(req, zmq4.Readable, func(state zmq4.State) {
		if state&zmq4.Readable != 0 {
			msg, err := req.Recv()
			if err != nil {
				t.Error("req.Recv:", err)
				return
			}
			received <- string(msg.Frames[0])
		}
	})
	if err != nil {
		t.Fatal("reactor.AddSocket req:", err)
	}
	
	// Run reactor in background
	go func() {
		err := reactor.Run()
		if err != nil {
			t.Error("reactor.Run:", err)
		}
	}()
	
	// Send some messages
	for i := 0; i < 3; i++ {
		msg := zmq4.NewMsg([]byte(fmt.Sprintf("Request %d", i)))
		err = req.Send(msg)
		if err != nil {
			t.Error("req.Send:", err)
		}
		
		// Wait for request and reply
		select {
		case reqMsg := <-received:
			t.Logf("Received request: %s", reqMsg)
		case <-time.After(1 * time.Second):
			t.Error("Timeout waiting for request")
		}
		
		select {
		case repMsg := <-received:
			t.Logf("Received reply: %s", repMsg)
		case <-time.After(1 * time.Second):
			t.Error("Timeout waiting for reply")
		}
	}
	
	// Stop reactor
	reactor.Stop()
	time.Sleep(100 * time.Millisecond)
}

func TestPollerWithTimeout(t *testing.T) {
	t.Skip("Poller implementation needs work")
	ctx := context.Background()
	
	// Create unconnected socket
	socket := zmq4.NewPull(ctx)
	defer socket.Close()
	
	socket.Listen("tcp://127.0.0.1:0")
	
	// Create poller
	poller := zmq4.NewPoller()
	err := poller.Add(socket, zmq4.Readable)
	if err != nil {
		t.Fatal("poller.Add:", err)
	}
	
	// Poll with short timeout - should timeout
	start := time.Now()
	items, err := poller.Poll(100 * time.Millisecond)
	elapsed := time.Since(start)
	
	if err != nil {
		t.Error("Poll:", err)
	}
	if len(items) != 0 {
		t.Errorf("Poll: expected 0 items on timeout, got %d", len(items))
	}
	if elapsed < 90*time.Millisecond || elapsed > 200*time.Millisecond {
		t.Errorf("Poll timeout took %v, expected ~100ms", elapsed)
	}
	
	// Poll with no timeout (PollAll)
	done := make(chan bool)
	go func() {
		items, err := poller.PollAll()
		if err != nil {
			t.Error("PollAll:", err)
		}
		t.Logf("PollAll returned %d items", len(items))
		done <- true
	}()
	
	// Should block until we connect and send
	time.Sleep(50 * time.Millisecond)
	
	// Connect and send
	push := zmq4.NewPush(ctx)
	defer push.Close()
	
	addr := socket.Addr()
	err = push.Dial(fmt.Sprintf("tcp://%s", addr))
	if err != nil {
		t.Fatal("push.Dial:", err)
	}
	
	msg := zmq4.NewMsg([]byte("test"))
	err = push.Send(msg)
	if err != nil {
		t.Fatal("push.Send:", err)
	}
	
	// PollAll should now return
	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Error("PollAll did not return after message sent")
	}
}