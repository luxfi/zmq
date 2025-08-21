# ZMQ4 Quick Reference

## Build Commands

### Pure Go (Default)
```bash
go build .                    # Build library
go test ./...                 # Run all tests
make test                     # Run tests via Makefile
make test-race               # Test with race detector
make test-cover              # Generate coverage report
```

### With CZMQ Compatibility
```bash
go build -tags czmq4 .        # Build with CZMQ
go test -tags czmq4 ./...     # Test with CZMQ
make test-czmq               # Test via Makefile
make test-czmq-race          # Race detector + CZMQ
```

## File Structure

```
zmq4/
├── Core Implementation (Pure Go)
│   ├── socket.go           # Socket implementation
│   ├── conn.go            # Connection handling
│   ├── msg.go             # Message types
│   └── *.go               # Other pure Go files
│
├── CZMQ Compatibility (Optional)
│   ├── cxx_zmq4_compat.go     # [czmq4 tag] CZMQ wrappers
│   ├── czmq4_test.go          # [czmq4 tag] Compat tests
│   └── zmq4_czmq_integration_test.go # [czmq4 tag] Integration
│
├── Documentation
│   ├── README.md              # Main documentation
│   ├── BUILD_TAGS.md          # Build tag details
│   ├── CZMQ_COMPATIBILITY.md  # Compatibility guide
│   └── QUICK_REFERENCE.md     # This file
│
├── Scripts
│   └── test_czmq_isolation.sh # Validation script
│
└── Examples
    ├── README.md              # Example documentation
    ├── compatibility_test.go  # [czmq4 tag] Compat demo
    └── *.go                   # Various examples
```

## Common Tasks

### Validate CZMQ Isolation
```bash
./scripts/test_czmq_isolation.sh
```

### Run Specific Tests
```bash
# Pure Go only
go test -run TestPubSub ./...

# With CZMQ
go test -tags czmq4 -run TestCZMQIntegration ./...

# Benchmarks
go test -bench . ./...
go test -tags czmq4 -bench BenchmarkCZMQvsPureGo ./...
```

### Check Build Configuration
```bash
# List files in pure Go build
go list -f '{{.GoFiles}}' .

# List files with CZMQ
go list -f '{{.GoFiles}}' -tags=czmq4 .

# Verify no CZMQ imports by default
go list -f '{{.Imports}}' . | grep -c czmq  # Should be 0
```

## Environment Variables

```bash
# For CZMQ builds
export CGO_ENABLED=1

# For pure Go builds
export CGO_ENABLED=0  # Optional, pure Go works with CGO enabled too
```

## Troubleshooting

### Issue: CZMQ build fails
```bash
# Ubuntu/Debian
sudo apt-get install libczmq-dev libzmq3-dev

# macOS
brew install czmq

# Verify installation
pkg-config --libs libczmq
```

### Issue: Tests timeout
```bash
# Increase timeout
go test -timeout 30s ./...

# Run with verbose output
go test -v ./...
```

### Issue: Import cycle
```bash
# Check dependencies
go mod graph | grep zmq

# Clean module cache
go clean -modcache
```

## Socket Types

| Type | Pure Go | CZMQ | Usage |
|------|---------|------|-------|
| REQ-REP | ✓ | ✓ | Request-Reply pattern |
| PUB-SUB | ✓ | ✓ | Publish-Subscribe pattern |
| PUSH-PULL | ✓ | ✓ | Pipeline pattern |
| DEALER-ROUTER | ✓ | ✓ | Advanced routing |
| PAIR | ✓ | ✓ | Exclusive pair |
| XPUB-XSUB | ✓ | ✓ | Extended Pub-Sub |

## Performance Tips

1. **Use pure Go for production** - No CGO overhead
2. **Batch messages** - Reduce syscalls
3. **Set high water marks** - Prevent memory issues
4. **Use appropriate socket types** - Match your pattern
5. **Profile your application** - `go test -cpuprofile`

## CI/CD Integration

The project uses GitHub Actions with three test strategies:
- `test-pure-go` - Pure Go, CGO disabled
- `test-cgo` - CGO enabled, no CZMQ
- `test-czmq4` - Full CZMQ compatibility

## Quick Decision Tree

```
Need ZMQ in Go?
├── Production deployment?
│   └── Use pure Go (default)
├── Testing compatibility?
│   └── Use -tags czmq4
├── Migrating from C?
│   └── Test both, deploy pure Go
└── Development?
    └── Test with both implementations
```

## Links

- [ZeroMQ Documentation](https://zeromq.org/documentation/)
- [Go Package Documentation](https://pkg.go.dev/github.com/luxfi/zmq/v4)
- [GitHub Repository](https://github.com/luxfi/zmq)
- [Issue Tracker](https://github.com/luxfi/zmq/issues)