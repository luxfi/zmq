//go:build cgo
// +build cgo

package zmq4_test

import (
	"context"
	"testing"
	"time"

	"github.com/luxfi/zmq/v4"
)

// Benchmark tests for CZMQ implementation
// These are the same benchmarks but using CZMQ sockets

func BenchmarkCZMQPubSub(b *testing.B) {
	ctx := context.Background()

	b.Run("Throughput-1KB", func(b *testing.B) {
		benchmarkCZMQPubSubThroughput(b, ctx, 1024)
	})

	b.Run("Throughput-10KB", func(b *testing.B) {
		benchmarkCZMQPubSubThroughput(b, ctx, 10*1024)
	})

	b.Run("Throughput-100KB", func(b *testing.B) {
		benchmarkCZMQPubSubThroughput(b, ctx, 100*1024)
	})

	b.Run("Latency", func(b *testing.B) {
		benchmarkCZMQPubSubLatency(b, ctx)
	})
}

func BenchmarkCZMQReqRep(b *testing.B) {
	ctx := context.Background()

	b.Run("RoundTrip-Small", func(b *testing.B) {
		benchmarkCZMQReqRepRoundTrip(b, ctx, 64)
	})

	b.Run("RoundTrip-Medium", func(b *testing.B) {
		benchmarkCZMQReqRepRoundTrip(b, ctx, 1024)
	})

	b.Run("RoundTrip-Large", func(b *testing.B) {
		benchmarkCZMQReqRepRoundTrip(b, ctx, 10*1024)
	})
}

func benchmarkCZMQPubSubThroughput(b *testing.B, ctx context.Context, msgSize int) {
	pub := zmq4.NewCPub(ctx)
	defer pub.Close()
	sub := zmq4.NewCSub(ctx)
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

func benchmarkCZMQPubSubLatency(b *testing.B, ctx context.Context) {
	pub := zmq4.NewCPub(ctx)
	defer pub.Close()
	sub := zmq4.NewCSub(ctx)
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

func benchmarkCZMQReqRepRoundTrip(b *testing.B, ctx context.Context, msgSize int) {
	req := zmq4.NewCReq(ctx)
	defer req.Close()
	rep := zmq4.NewCRep(ctx)
	defer rep.Close()

	endpoint := must(EndPoint("tcp"))
	if err := rep.Listen(endpoint); err != nil {
		b.Fatal(err)
	}
	if err := req.Dial(endpoint); err != nil {
		b.Fatal(err)
	}

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

	data := make([]byte, msgSize)
	msg := zmq4.NewMsg(data)

	b.ResetTimer()
	b.SetBytes(int64(msgSize * 2))

	for i := 0; i < b.N; i++ {
		if err := req.Send(msg); err != nil {
			b.Fatal(err)
		}
		if _, err := req.Recv(); err != nil {
			b.Fatal(err)
		}
	}
}
