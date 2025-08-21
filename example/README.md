# ZMQ4 Examples

This directory contains examples demonstrating the usage of the zmq4 package.

## Basic Examples

### Hello World (REQ-REP)
Simple request-reply pattern:
```bash
# Terminal 1
go run hwserver.go

# Terminal 2
go run hwclient.go
```

### Pub-Sub with Envelope
Publisher-subscriber pattern with message envelopes:
```bash
# Terminal 1
go run psenvpub.go

# Terminal 2
go run psenvsub.go
```

### Request-Reply with Workers
Load-balanced request distribution:
```bash
# Terminal 1
go run rrworker.go

# Terminal 2
go run rrclient.go
```

### Router-Dealer
Advanced routing pattern:
```bash
go run rtdealer.go
```

## CZMQ Compatibility Testing

The `compatibility_test.go` example demonstrates interoperability between pure Go and CZMQ implementations.

### Prerequisites
To run the compatibility test, you need CZMQ installed:

**Ubuntu/Debian:**
```bash
sudo apt-get install libczmq-dev libzmq3-dev
```

**macOS:**
```bash
brew install czmq
```

### Running the Compatibility Test
```bash
# Build and run with czmq4 tag
go run -tags czmq4 compatibility_test.go
```

This example creates a CZMQ publisher and a pure Go subscriber to verify that messages can be exchanged between the two implementations.

## Pure Go vs CZMQ

### Pure Go (Default)
All examples run with the pure Go implementation by default:
```bash
go run hwserver.go  # Uses pure Go implementation
```

### With CZMQ
To run any example with CZMQ compatibility:
```bash
go run -tags czmq4 hwserver.go  # Uses CZMQ for socket operations
```

## Testing Interoperability

You can test interoperability by running one component with CZMQ and another without:

```bash
# Terminal 1: Server using CZMQ
go run -tags czmq4 hwserver.go

# Terminal 2: Client using pure Go
go run hwclient.go
```

Both should communicate seamlessly, demonstrating protocol compatibility.