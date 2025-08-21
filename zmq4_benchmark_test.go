//go:build !cgo
// +build !cgo

package zmq4_test

import (
	"context"
	"testing"
	"time"

	"github.com/luxfi/zmq/v4"
)

// Benchmark tests for Pure Go implementation
// Run with: go test -bench=. -benchmem ./...
// Compare with: go test -tags czmq4 -bench=. -benchmem ./...

func BenchmarkPureGoPubSub(b *testing.B) {
	ctx := context.Background()

	b.Run("Throughput-1KB", func(b *testing.B) {
		benchmarkPubSubThroughput(b, ctx, 1024)
	})

	b.Run("Throughput-10KB", func(b *testing.B) {
		benchmarkPubSubThroughput(b, ctx, 10*1024)
	})

	b.Run("Throughput-100KB", func(b *testing.B) {
		benchmarkPubSubThroughput(b, ctx, 100*1024)
	})

	b.Run("Latency", func(b *testing.B) {
		benchmarkPubSubLatency(b, ctx)
	})
}

func BenchmarkPureGoReqRep(b *testing.B) {
	ctx := context.Background()

	b.Run("RoundTrip-Small", func(b *testing.B) {
		benchmarkReqRepRoundTrip(b, ctx, 64)
	})

	b.Run("RoundTrip-Medium", func(b *testing.B) {
		benchmarkReqRepRoundTrip(b, ctx, 1024)
	})

	b.Run("RoundTrip-Large", func(b *testing.B) {
		benchmarkReqRepRoundTrip(b, ctx, 10*1024)
	})
}

func BenchmarkPureGoPushPull(b *testing.B) {
	ctx := context.Background()

	b.Run("Pipeline-Throughput", func(b *testing.B) {
		benchmarkPushPullThroughput(b, ctx)
	})

	b.Run("Pipeline-FanOut", func(b *testing.B) {
		benchmarkPushPullFanOut(b, ctx, 4)
	})
}

func BenchmarkPureGoRouterDealer(b *testing.B) {
	ctx := context.Background()

	b.Run("Routing-Performance", func(b *testing.B) {
		benchmarkRouterDealerPerformance(b, ctx)
	})
}

// Helper functions for benchmarks

func benchmarkPubSubThroughput(b *testing.B, ctx context.Context, msgSize int) {
	pub := zmq4.NewPub(ctx)
	defer pub.Close()
	sub := zmq4.NewSub(ctx)
	defer sub.Close()
	sub.SetOption(zmq4.OptionSubscribe, "")

	endpoint := must(EndPoint("tcp"))
	if err := pub.Listen(endpoint); err != nil {
		b.Fatal(err)
	}
	if err := sub.Dial(endpoint); err != nil {
		b.Fatal(err)
	}

	// Allow connection to establish
	time.Sleep(100 * time.Millisecond)

	// Prepare message
	data := make([]byte, msgSize)
	for i := range data {
		data[i] = byte(i % 256)
	}
	msg := zmq4.NewMsg(data)

	b.ResetTimer()
	b.SetBytes(int64(msgSize))

	for i := 0; i < b.N; i++ {
		if err := pub.Send(msg); err != nil {
			b.Fatal(err)
		}
		if _, err := sub.Recv(); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkPubSubLatency(b *testing.B, ctx context.Context) {
	pub := zmq4.NewPub(ctx)
	defer pub.Close()
	sub := zmq4.NewSub(ctx)
	defer sub.Close()
	sub.SetOption(zmq4.OptionSubscribe, "")

	endpoint := must(EndPoint("tcp"))
	if err := pub.Listen(endpoint); err != nil {
		b.Fatal(err)
	}
	if err := sub.Dial(endpoint); err != nil {
		b.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond)

	msg := zmq4.NewMsgString("latency-test")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		start := time.Now()
		if err := pub.Send(msg); err != nil {
			b.Fatal(err)
		}
		if _, err := sub.Recv(); err != nil {
			b.Fatal(err)
		}
		b.ReportMetric(float64(time.Since(start).Nanoseconds()), "ns/op")
	}
}

func benchmarkReqRepRoundTrip(b *testing.B, ctx context.Context, msgSize int) {
	req := zmq4.NewReq(ctx)
	defer req.Close()
	rep := zmq4.NewRep(ctx)
	defer rep.Close()

	endpoint := must(EndPoint("tcp"))
	if err := rep.Listen(endpoint); err != nil {
		b.Fatal(err)
	}
	if err := req.Dial(endpoint); err != nil {
		b.Fatal(err)
	}

	// Start reply server
	go func() {
		for {
			msg, err := rep.Recv()
			if err != nil {
				return
			}
			if err := rep.Send(msg); err != nil {
				return
			}
		}
	}()

	// Prepare message
	data := make([]byte, msgSize)
	msg := zmq4.NewMsg(data)

	b.ResetTimer()
	b.SetBytes(int64(msgSize * 2)) // Request + Reply

	for i := 0; i < b.N; i++ {
		if err := req.Send(msg); err != nil {
			b.Fatal(err)
		}
		if _, err := req.Recv(); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkPushPullThroughput(b *testing.B, ctx context.Context) {
	push := zmq4.NewPush(ctx)
	defer push.Close()
	pull := zmq4.NewPull(ctx)
	defer pull.Close()

	endpoint := must(EndPoint("tcp"))
	if err := pull.Listen(endpoint); err != nil {
		b.Fatal(err)
	}
	if err := push.Dial(endpoint); err != nil {
		b.Fatal(err)
	}

	msg := zmq4.NewMsgString("pipeline-message")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := push.Send(msg); err != nil {
			b.Fatal(err)
		}
		if _, err := pull.Recv(); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkPushPullFanOut(b *testing.B, ctx context.Context, workers int) {
	push := zmq4.NewPush(ctx)
	defer push.Close()

	endpoint := must(EndPoint("tcp"))
	if err := push.Listen(endpoint); err != nil {
		b.Fatal(err)
	}

	// Create worker pulls
	pulls := make([]zmq4.Socket, workers)
	for i := 0; i < workers; i++ {
		pulls[i] = zmq4.NewPull(ctx)
		defer pulls[i].Close()
		if err := pulls[i].Dial(endpoint); err != nil {
			b.Fatal(err)
		}
	}

	// Start workers
	done := make(chan bool, workers)
	for i := 0; i < workers; i++ {
		go func(pull zmq4.Socket) {
			for {
				if _, err := pull.Recv(); err != nil {
					done <- true
					return
				}
			}
		}(pulls[i])
	}

	msg := zmq4.NewMsgString("work-item")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := push.Send(msg); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkRouterDealerPerformance(b *testing.B, ctx context.Context) {
	router := zmq4.NewRouter(ctx)
	defer router.Close()
	dealer := zmq4.NewDealer(ctx)
	defer dealer.Close()

	endpoint := must(EndPoint("tcp"))
	if err := router.Listen(endpoint); err != nil {
		b.Fatal(err)
	}
	if err := dealer.Dial(endpoint); err != nil {
		b.Fatal(err)
	}

	// Start router handler
	go func() {
		for {
			msg, err := router.Recv()
			if err != nil {
				return
			}
			if err := router.Send(msg); err != nil {
				return
			}
		}
	}()

	msg := zmq4.NewMsgString("routed-message")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := dealer.Send(msg); err != nil {
			b.Fatal(err)
		}
		if _, err := dealer.Recv(); err != nil {
			b.Fatal(err)
		}
	}
}
