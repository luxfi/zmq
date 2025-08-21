//go:build cgo
// +build cgo

package zmq4_test

import (
	"context"
	"testing"
	"time"

	"github.com/luxfi/zmq/v4"
)

// TestCZMQIntegration verifies that CZMQ sockets can communicate with pure Go sockets
func TestCZMQIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Run("CZMQ-Pub-Go-Sub", func(t *testing.T) {
		testCZMQPubGoSub(t, ctx)
	})

	t.Run("Go-Pub-CZMQ-Sub", func(t *testing.T) {
		testGoPubCZMQSub(t, ctx)
	})

	t.Run("CZMQ-Req-Go-Rep", func(t *testing.T) {
		testCZMQReqGoRep(t, ctx)
	})

	t.Run("Go-Req-CZMQ-Rep", func(t *testing.T) {
		testGoReqCZMQRep(t, ctx)
	})
}

func testCZMQPubGoSub(t *testing.T, ctx context.Context) {
	// Create CZMQ publisher
	pub := zmq4.NewCPub(ctx)
	defer pub.Close()

	// Create Go subscriber
	sub := zmq4.NewSub(ctx)
	defer sub.Close()
	sub.SetOption(zmq4.OptionSubscribe, "")

	// Setup connection
	endpoint := must(EndPoint("tcp"))
	if err := pub.Listen(endpoint); err != nil {
		t.Fatalf("could not listen: %v", err)
	}
	if err := sub.Dial(endpoint); err != nil {
		t.Fatalf("could not dial: %v", err)
	}

	// Allow connection to establish
	time.Sleep(100 * time.Millisecond)

	// Send from CZMQ
	msg := zmq4.NewMsgString("test-czmq-pub")
	if err := pub.Send(msg); err != nil {
		t.Fatalf("could not send: %v", err)
	}

	// Receive in Go
	received, err := sub.Recv()
	if err != nil {
		t.Fatalf("could not receive: %v", err)
	}

	if string(received.Bytes()) != "test-czmq-pub" {
		t.Errorf("expected 'test-czmq-pub', got '%s'", string(received.Bytes()))
	}
}

func testGoPubCZMQSub(t *testing.T, ctx context.Context) {
	// Create Go publisher
	pub := zmq4.NewPub(ctx)
	defer pub.Close()

	// Create CZMQ subscriber
	sub := zmq4.NewCSub(ctx)
	defer sub.Close()
	sub.SetOption(zmq4.OptionSubscribe, "")

	// Setup connection
	endpoint := must(EndPoint("tcp"))
	if err := pub.Listen(endpoint); err != nil {
		t.Fatalf("could not listen: %v", err)
	}
	if err := sub.Dial(endpoint); err != nil {
		t.Fatalf("could not dial: %v", err)
	}

	// Allow connection to establish
	time.Sleep(100 * time.Millisecond)

	// Send from Go
	msg := zmq4.NewMsgString("test-go-pub")
	if err := pub.Send(msg); err != nil {
		t.Fatalf("could not send: %v", err)
	}

	// Receive in CZMQ
	received, err := sub.Recv()
	if err != nil {
		t.Fatalf("could not receive: %v", err)
	}

	if string(received.Bytes()) != "test-go-pub" {
		t.Errorf("expected 'test-go-pub', got '%s'", string(received.Bytes()))
	}
}

func testCZMQReqGoRep(t *testing.T, ctx context.Context) {
	// Create CZMQ requester
	req := zmq4.NewCReq(ctx)
	defer req.Close()

	// Create Go replier
	rep := zmq4.NewRep(ctx)
	defer rep.Close()

	// Setup connection
	endpoint := must(EndPoint("tcp"))
	if err := rep.Listen(endpoint); err != nil {
		t.Fatalf("could not listen: %v", err)
	}
	if err := req.Dial(endpoint); err != nil {
		t.Fatalf("could not dial: %v", err)
	}

	// Send request from CZMQ
	reqMsg := zmq4.NewMsgString("request-from-czmq")
	if err := req.Send(reqMsg); err != nil {
		t.Fatalf("could not send request: %v", err)
	}

	// Receive in Go
	received, err := rep.Recv()
	if err != nil {
		t.Fatalf("could not receive request: %v", err)
	}

	if string(received.Bytes()) != "request-from-czmq" {
		t.Errorf("expected 'request-from-czmq', got '%s'", string(received.Bytes()))
	}

	// Send reply from Go
	repMsg := zmq4.NewMsgString("reply-from-go")
	if err := rep.Send(repMsg); err != nil {
		t.Fatalf("could not send reply: %v", err)
	}

	// Receive reply in CZMQ
	reply, err := req.Recv()
	if err != nil {
		t.Fatalf("could not receive reply: %v", err)
	}

	if string(reply.Bytes()) != "reply-from-go" {
		t.Errorf("expected 'reply-from-go', got '%s'", string(reply.Bytes()))
	}
}

func testGoReqCZMQRep(t *testing.T, ctx context.Context) {
	// Create Go requester
	req := zmq4.NewReq(ctx)
	defer req.Close()

	// Create CZMQ replier
	rep := zmq4.NewCRep(ctx)
	defer rep.Close()

	// Setup connection
	endpoint := must(EndPoint("tcp"))
	if err := rep.Listen(endpoint); err != nil {
		t.Fatalf("could not listen: %v", err)
	}
	if err := req.Dial(endpoint); err != nil {
		t.Fatalf("could not dial: %v", err)
	}

	// Send request from Go
	reqMsg := zmq4.NewMsgString("request-from-go")
	if err := req.Send(reqMsg); err != nil {
		t.Fatalf("could not send request: %v", err)
	}

	// Receive in CZMQ
	received, err := rep.Recv()
	if err != nil {
		t.Fatalf("could not receive request: %v", err)
	}

	if string(received.Bytes()) != "request-from-go" {
		t.Errorf("expected 'request-from-go', got '%s'", string(received.Bytes()))
	}

	// Send reply from CZMQ
	repMsg := zmq4.NewMsgString("reply-from-czmq")
	if err := rep.Send(repMsg); err != nil {
		t.Fatalf("could not send reply: %v", err)
	}

	// Receive reply in Go
	reply, err := req.Recv()
	if err != nil {
		t.Fatalf("could not receive reply: %v", err)
	}

	if string(reply.Bytes()) != "reply-from-czmq" {
		t.Errorf("expected 'reply-from-czmq', got '%s'", string(reply.Bytes()))
	}
}

// BenchmarkCZMQvsPureGo compares performance between CZMQ and pure Go implementations
func BenchmarkCZMQvsPureGo(b *testing.B) {
	ctx := context.Background()

	b.Run("CZMQ-Pub-Sub", func(b *testing.B) {
		pub := zmq4.NewCPub(ctx)
		defer pub.Close()
		sub := zmq4.NewCSub(ctx)
		defer sub.Close()
		sub.SetOption(zmq4.OptionSubscribe, "")

		endpoint := must(EndPoint("tcp"))
		pub.Listen(endpoint)
		sub.Dial(endpoint)
		time.Sleep(100 * time.Millisecond)

		msg := zmq4.NewMsgString("benchmark")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pub.Send(msg)
			sub.Recv()
		}
	})

	b.Run("PureGo-Pub-Sub", func(b *testing.B) {
		pub := zmq4.NewPub(ctx)
		defer pub.Close()
		sub := zmq4.NewSub(ctx)
		defer sub.Close()
		sub.SetOption(zmq4.OptionSubscribe, "")

		endpoint := must(EndPoint("tcp"))
		pub.Listen(endpoint)
		sub.Dial(endpoint)
		time.Sleep(100 * time.Millisecond)

		msg := zmq4.NewMsgString("benchmark")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pub.Send(msg)
			sub.Recv()
		}
	})
}
