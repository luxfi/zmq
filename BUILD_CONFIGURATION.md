# Build Configuration

This document describes the build configuration for the zmq4 package.

## Overview

The zmq4 package automatically selects between pure Go and CZMQ implementations based on CGO availability. This provides a seamless experience where the package uses pure Go by default but can leverage CZMQ when CGO is enabled.

## Build Modes

### Pure Go Mode (Default)

When `CGO_ENABLED=0` or CGO is not available, the package uses a pure Go implementation with zero C dependencies.

**Files used in Pure Go mode:**
- All standard Go files (socket.go, conn.go, msg.go, etc.)
- No C dependencies required
- Cross-platform compatible

### CZMQ Compatibility Mode

When `CGO_ENABLED=1` and CZMQ is installed, the package automatically uses CZMQ for compatibility testing.

**Files added with CGO enabled:**
- `cxx_zmq4_compat.go` - Provides C-compatible socket implementations
- `czmq4_test.go` - Tests for CZMQ compatibility
- `security/plain/plain_cxx_test.go` - PLAIN security tests with CZMQ

## Usage

### Default Build (Pure Go)

By default, the package uses pure Go:

```bash
# Build
go build .

# Test
go test ./...

# Explicitly disable CGO
CGO_ENABLED=0 go build .
```

### Building with CZMQ

When CGO is enabled and CZMQ is installed:

```bash
# Build with CGO
CGO_ENABLED=1 go build .

# Test with CGO
CGO_ENABLED=1 go test ./...

# Using make
make test-czmq
```

## CI/CD Integration

The GitHub Actions workflow includes three testing strategies:

1. **Pure Go Tests** - Tests the pure-Go implementation with CGO disabled
2. **CGO Tests** - Tests with CGO enabled but without CZMQ dependency
3. **CZMQ4 Tests** - Tests with the CZMQ compatibility layer enabled

This ensures that:
- The pure-Go implementation works correctly on its own
- The package works with CGO enabled
- The implementation maintains compatibility with CZMQ v4

## Dependencies

### Without Build Tags
- No external C dependencies
- Pure Go implementation

### With `czmq4` Build Tag
- Requires CZMQ v4 library
- Requires libzmq3-dev
- Requires CGO_ENABLED=1

On Ubuntu/Debian:
```bash
sudo apt-get install libczmq-dev libzmq3-dev
```

On macOS:
```bash
brew install czmq
```

## Architecture

```
zmq4/
├── Pure Go Implementation (default)
│   ├── socket.go
│   ├── conn.go
│   ├── msg.go
│   └── ... (other pure Go files)
│
└── CZMQ Compatibility Layer (with czmq4 tag)
    ├── cxx_zmq4_compat.go
    └── czmq4_test.go
```

## Testing Strategy

The package maintains two parallel testing approaches:

1. **Unit Tests** - Test the pure-Go implementation functionality
2. **Compatibility Tests** - Verify behavior matches CZMQ implementation

This dual approach ensures:
- The pure-Go implementation is fully functional
- Behavior remains consistent with the ZeroMQ specification
- Interoperability with C implementations is maintained

## Best Practices

1. **Production Use**: Always use the pure-Go implementation (no build tags) for production deployments to avoid C dependencies.

2. **Development**: Test both with and without the `czmq4` tag to ensure compatibility:
   ```bash
   make test         # Pure Go tests
   make test-czmq    # Compatibility tests
   ```

3. **CI/CD**: Include both test suites in your CI pipeline to catch any compatibility issues early.

## Troubleshooting

### CZMQ Tests Failing

If CZMQ tests fail but pure-Go tests pass:
1. Ensure CZMQ v4 is properly installed
2. Check that CGO_ENABLED=1
3. Verify the luxfi/czmq/v4 module is available

### Build Tag Not Working

Ensure you're using the correct syntax:
- Correct: `go build -tags czmq4`
- Incorrect: `go build -tags=czmq4` (note the equals sign)

### Missing Functions

The CZMQ compatibility layer only provides a subset of functions needed for testing. If you need additional CZMQ functions, they should be added to `cxx_zmq4_compat.go` with the appropriate build constraints.

## Future Considerations

The CZMQ compatibility layer is intended solely for testing and validation. Future versions may:
- Expand compatibility test coverage
- Add additional build tags for specific features
- Remove CZMQ dependency entirely once full compatibility is verified