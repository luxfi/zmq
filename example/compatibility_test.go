//go:build cgo
// +build cgo

// This example demonstrates using the CZMQ compatibility layer.
// Build and run with: CGO_ENABLED=1 go run compatibility_test.go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/luxfi/zmq/v4"
)

func main() {
	fmt.Println("Running with CZMQ compatibility layer enabled")
	fmt.Println("This example tests interoperability between pure Go and CZMQ sockets")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a CZMQ publisher
	pub := zmq4.NewCPub(ctx)
	defer pub.Close()

	// Create a pure Go subscriber
	sub := zmq4.NewSub(ctx)
	defer sub.Close()
	sub.SetOption(zmq4.OptionSubscribe, "")

	// Bind publisher and connect subscriber
	if err := pub.Listen("tcp://127.0.0.1:5555"); err != nil {
		log.Fatalf("could not listen: %v", err)
	}

	if err := sub.Dial("tcp://127.0.0.1:5555"); err != nil {
		log.Fatalf("could not dial: %v", err)
	}

	// Give sockets time to connect
	time.Sleep(100 * time.Millisecond)

	// Send from CZMQ publisher
	msg := zmq4.NewMsgString("Hello from CZMQ!")
	if err := pub.Send(msg); err != nil {
		log.Fatalf("could not send: %v", err)
	}

	// Receive in pure Go subscriber
	received, err := sub.Recv()
	if err != nil {
		log.Fatalf("could not receive: %v", err)
	}

	fmt.Printf("Pure Go subscriber received: %s\n", string(received.Bytes()))
	fmt.Println("✓ CZMQ compatibility test passed!")
}
