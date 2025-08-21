# zmq4 - LuxFi Fork

[![GitHub release](https://img.shields.io/github/release/luxfi/zmq.svg)](https://github.com/luxfi/zmq/releases)
[![go.dev reference](https://pkg.go.dev/badge/github.com/luxfi/zmq/v4)](https://pkg.go.dev/github.com/luxfi/zmq/v4)
[![CI](https://github.com/luxfi/zmq/workflows/CI/badge.svg)](https://github.com/luxfi/zmq/actions)
[![Pure Go Tests](https://img.shields.io/github/actions/workflow/status/luxfi/zmq/ci.yml?label=pure-go&branch=main)](https://github.com/luxfi/zmq/actions)
[![CZMQ Tests](https://img.shields.io/github/actions/workflow/status/luxfi/zmq/ci.yml?label=czmq-compat&branch=main)](https://github.com/luxfi/zmq/actions)
[![Coverage](https://img.shields.io/codecov/c/github/luxfi/zmq/main.svg)](https://codecov.io/gh/luxfi/zmq)
[![GoDoc](https://godoc.org/github.com/luxfi/zmq/v4?status.svg)](https://godoc.org/github.com/luxfi/zmq/v4)
[![Go Report Card](https://goreportcard.com/badge/github.com/luxfi/zmq/v4)](https://goreportcard.com/report/github.com/luxfi/zmq/v4)
[![License](https://img.shields.io/badge/License-BSD--3-blue.svg)](https://github.com/luxfi/zmq/blob/main/LICENSE)

`zmq4` is a pure-Go implementation of ØMQ (ZeroMQ), version 4.

This is the LuxFi fork of the original [go-zeromq/zmq4](https://github.com/go-zeromq/zmq4) project, maintained for use in the Lux Network ecosystem.

See [zeromq.org](http://zeromq.org) for more information about ZeroMQ.

## Installation

```bash
go get github.com/luxfi/zmq/v4
```

## Usage

```go
import "github.com/luxfi/zmq/v4"
```

## Version

This fork maintains compatibility with ZeroMQ v4.2.0 and follows Go module versioning with the `/v4` suffix.

## Development

This fork is maintained by the LuxFi team for use in the Lux Network consensus and networking layers.

### Build Modes

This package automatically selects the implementation based on CGO:

- **CGO_ENABLED=0** (or unset): Pure Go implementation (default, no C dependencies)
- **CGO_ENABLED=1**: CZMQ compatibility layer (requires libczmq)

### Testing

```bash
# Run tests with pure Go implementation
CGO_ENABLED=0 go test ./...
# or simply:
make test

# Run tests with CZMQ compatibility layer
CGO_ENABLED=1 go test ./...
# or:
make test-czmq

# Run tests with verbose output
make test-verbose

# Run tests with race detector
make test-race

# Run tests with coverage
make test-cover
```

## License

`zmq4` is released under the `BSD-3` license.

## Documentation

Documentation for `zmq4` is served by [GoDoc](https://godoc.org/github.com/luxfi/zmq/v4).

## Dependencies

### Pure Go (Default)
By default, this package is a pure-Go implementation with no external dependencies.

### Optional CZMQ Compatibility
For compatibility testing with the C implementation, this package optionally depends on:
- [github.com/luxfi/czmq/v4](https://github.com/luxfi/czmq) v4.2.0 (only when using `czmq4` build tag)

## Contributing

Contributions are welcome! Please submit issues and pull requests to the [LuxFi ZMQ repository](https://github.com/luxfi/zmq).


## Original Project

This is a fork of [go-zeromq/zmq4](https://github.com/go-zeromq/zmq4). The original project is licensed under BSD-3.

## CI Status
