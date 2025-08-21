# CZMQ Compatibility Layer

## Overview

The zmq4 package includes an optional CZMQ (C ZeroMQ) compatibility layer that allows for:
- Testing interoperability between pure Go and C implementations
- Validating protocol compliance
- Performance benchmarking
- Migration support for systems transitioning from C to Go

## Architecture

```
┌─────────────────────────────────────────────┐
│             Application Code                 │
└─────────────┬───────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────┐
│            zmq4 Public API                   │
│         (Socket, Msg, Options)               │
└──────┬──────────────────────┬────────────────┘
       │                      │
       │ Default              │ With czmq4 tag
       ▼                      ▼
┌──────────────┐      ┌──────────────────────┐
│  Pure Go     │      │  CZMQ Compatibility  │
│Implementation│      │     Layer            │
│  socket.go   │      │ cxx_zmq4_compat.go   │
└──────────────┘      └───────────┬──────────┘
                                   │
                                   ▼
                          ┌────────────────┐
                          │   libczmq.so   │
                          │   (C library)  │
                          └────────────────┘
```

## Build Configuration

### Default Build (Pure Go)
```bash
go build .                    # Pure Go, no C dependencies
go test ./...                 # Run tests with pure Go
make test                     # Using Makefile
```

### With CZMQ Compatibility
```bash
go build -tags czmq4 .        # Include CZMQ compatibility
go test -tags czmq4 ./...     # Run tests with CZMQ
make test-czmq                # Using Makefile
```

## Files Affected by Build Tags

| File | Build Tag | Purpose |
|------|-----------|---------|
| `cxx_zmq4_compat.go` | `czmq4` | CZMQ socket implementations |
| `czmq4_test.go` | `czmq4` | Compatibility tests |
| `security/plain/plain_cxx_test.go` | `czmq4` | PLAIN security tests |
| `zmq4_czmq_integration_test.go` | `czmq4` | Integration tests |

## Testing Strategy

### 1. Unit Tests (Pure Go)
Standard unit tests that verify the pure Go implementation:
```bash
make test
```

### 2. Compatibility Tests (CZMQ)
Tests that verify interoperability with CZMQ:
```bash
make test-czmq
```

### 3. Integration Tests
End-to-end tests that verify communication between Go and CZMQ sockets:
```bash
go test -tags czmq4 -run TestCZMQIntegration ./...
```

### 4. Performance Benchmarks
Compare performance between implementations:
```bash
go test -tags czmq4 -bench BenchmarkCZMQvsPureGo ./...
```

## CI/CD Pipeline

The GitHub Actions workflow includes three test strategies:

1. **test-pure-go**: Tests pure Go implementation (CGO_ENABLED=0)
2. **test-cgo**: Tests with CGO enabled but without CZMQ dependency
3. **test-czmq4**: Tests with CZMQ compatibility layer enabled

## Installation Requirements

### Pure Go (Default)
No additional requirements - works out of the box.

### With CZMQ Support

#### Ubuntu/Debian
```bash
sudo apt-get update
sudo apt-get install libczmq-dev libzmq3-dev
```

#### macOS
```bash
brew install czmq
```

#### RHEL/CentOS/Fedora
```bash
sudo yum install czmq-devel zeromq-devel
```

## Usage Examples

### Interoperability Testing
```go
// Create a CZMQ publisher
pub := zmq4.NewCPub(ctx)
defer pub.Close()

// Create a pure Go subscriber
sub := zmq4.NewSub(ctx)
defer sub.Close()

// They can communicate seamlessly
pub.Listen("tcp://127.0.0.1:5555")
sub.Dial("tcp://127.0.0.1:5555")
```

### Running Examples
```bash
# Pure Go example
go run example/hwserver.go

# With CZMQ compatibility
go run -tags czmq4 example/hwserver.go

# Test interoperability
go run -tags czmq4 example/compatibility_test.go
```

## Validation

Run the validation script to ensure proper isolation:
```bash
./scripts/test_czmq_isolation.sh
```

This script verifies:
- Pure Go builds work without CZMQ
- Build tags correctly isolate CZMQ code
- Test files are properly tagged
- No CZMQ imports in default build
- Makefile targets are properly configured

## Migration Guide

### From CZMQ to Pure Go

1. **Remove build tags**: Simply build without the `czmq4` tag
2. **No code changes needed**: The API remains the same
3. **Test thoroughly**: Run your test suite without the tag
4. **Deploy**: No C library dependencies needed

### Testing During Migration

```bash
# Test with CZMQ (current system)
go test -tags czmq4 ./...

# Test with pure Go (target system)
go test ./...

# Compare performance
go test -tags czmq4 -bench . ./...
go test -bench . ./...
```

## Troubleshooting

### Build Errors with czmq4 Tag

**Error**: `undefined reference to czmq functions`
**Solution**: Install CZMQ development libraries

**Error**: `cannot find -lczmq`
**Solution**: Ensure CZMQ is in your library path:
```bash
export LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH
```

### Test Failures

**Issue**: Tests pass without czmq4 but fail with it
**Likely Cause**: Version mismatch between CZMQ and ZeroMQ
**Solution**: Ensure compatible versions:
- CZMQ v4.2.0+
- ZeroMQ v4.x

### Performance Differences

The pure Go implementation may have different performance characteristics:
- **Latency**: Generally comparable
- **Throughput**: May vary based on message patterns
- **Memory usage**: Go's garbage collector vs C's manual memory management

## Best Practices

1. **Production Use**: Always use pure Go implementation for production
2. **Development**: Test with both implementations during development
3. **CI/CD**: Include both test suites in your pipeline
4. **Documentation**: Document which implementation your system uses
5. **Monitoring**: Track performance metrics when migrating

## Future Roadmap

- [ ] Expand compatibility test coverage
- [ ] Add more socket patterns to compatibility layer
- [ ] Performance optimization for pure Go implementation
- [ ] Remove CZMQ dependency once full compatibility is verified

## Support

For issues related to:
- **Pure Go implementation**: Open issue on [github.com/luxfi/zmq](https://github.com/luxfi/zmq)
- **CZMQ compatibility**: Check CZMQ version and library installation
- **Protocol compliance**: Run compatibility tests with `-v` flag for details