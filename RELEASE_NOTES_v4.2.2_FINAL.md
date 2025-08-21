# Release v4.2.2 - Automatic CGO-based Implementation Selection

## Overview

This patch release simplifies the build configuration by automatically selecting between pure Go and CZMQ implementations based on CGO availability. No API or protocol changes - fully backward compatible.

## 🎯 Key Changes

### Automatic Implementation Selection
- **CGO_ENABLED=0** (default): Pure Go implementation with zero C dependencies
- **CGO_ENABLED=1**: Automatically uses CZMQ if available for compatibility testing
- No more manual build tags required!

### What Changed
- Replaced `czmq4` build tag with automatic CGO detection
- CZMQ compatibility layer now activates when `CGO_ENABLED=1`
- Simplified build process - just works out of the box

## 📦 Installation

### Pure Go (Default)
```bash
go get github.com/luxfi/zmq/v4@v4.2.2
```

### With CZMQ Support
```bash
# Install CZMQ first (if needed)
sudo apt-get install libczmq-dev  # Ubuntu/Debian
brew install czmq                 # macOS

# Then build with CGO
CGO_ENABLED=1 go build
```

## 🔧 Usage

### Default Build (Pure Go)
```bash
go build .           # Pure Go, no dependencies
go test ./...        # Run tests with pure Go
```

### With CZMQ Compatibility
```bash
CGO_ENABLED=1 go build .    # Uses CZMQ if available
CGO_ENABLED=1 go test ./... # Test with CZMQ
```

### Using Makefile
```bash
make test           # Pure Go tests
make test-czmq      # CZMQ tests (sets CGO_ENABLED=1)
```

## ✨ Benefits

1. **Simpler**: No need to remember build tags
2. **Automatic**: Detects and uses CZMQ when CGO is enabled
3. **Backward Compatible**: No breaking changes
4. **Flexible**: Easy switch between implementations

## 📚 Documentation

- Updated README with new build instructions
- Simplified BUILD_TAGS.md → BUILD_CONFIGURATION.md
- Enhanced examples and documentation
- Added comprehensive testing tools

## 🔄 Migration

No migration needed! The package automatically selects the appropriate implementation:
- Existing pure Go builds continue to work
- CGO builds automatically get CZMQ compatibility
- All APIs remain unchanged

## 🧪 Testing

The package includes comprehensive tests for both modes:
```bash
# Test pure Go
CGO_ENABLED=0 go test ./...

# Test with CZMQ
CGO_ENABLED=1 go test ./...

# Run benchmarks
./scripts/benchmark_comparison.sh
```

## 📝 Technical Details

- **Build constraints**: Files use `//go:build cgo` instead of custom tags
- **Automatic selection**: CGO presence determines implementation
- **Zero overhead**: Pure Go remains the default with no dependencies
- **Full compatibility**: CZMQ layer ensures protocol compliance

## 🐛 Fixes

- Simplified build configuration
- Removed dependency on manual build tags
- Fixed import paths in examples
- Cleaned up go.mod dependencies

## 📊 Stats

- **Files changed**: 15
- **Build modes**: 2 (Pure Go / CZMQ)
- **Default mode**: Pure Go (zero dependencies)
- **Backward compatibility**: 100%

---

This is a patch release (v4.2.2) focused on build improvements. The ZMQ protocol remains v4.2.x compatible with no breaking changes.