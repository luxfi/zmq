# zmq4 - LuxFi Fork

[![GitHub release](https://img.shields.io/github/release/luxfi/zmq.svg)](https://github.com/luxfi/zmq/releases)
[![go.dev reference](https://pkg.go.dev/badge/github.com/luxfi/zmq/v4)](https://pkg.go.dev/github.com/luxfi/zmq/v4)
[![CI](https://github.com/luxfi/zmq/workflows/CI/badge.svg)](https://github.com/luxfi/zmq/actions)
[![GoDoc](https://godoc.org/github.com/luxfi/zmq/v4?status.svg)](https://godoc.org/github.com/luxfi/zmq/v4)
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

## License

`zmq4` is released under the `BSD-3` license.

## Documentation

Documentation for `zmq4` is served by [GoDoc](https://godoc.org/github.com/luxfi/zmq/v4).

## Dependencies

This package depends on:
- [github.com/luxfi/czmq/v4](https://github.com/luxfi/czmq) v4.2.0

## Contributing

Contributions are welcome! Please submit issues and pull requests to the [LuxFi ZMQ repository](https://github.com/luxfi/zmq).


## Original Project

This is a fork of [go-zeromq/zmq4](https://github.com/go-zeromq/zmq4). The original project is licensed under BSD-3.

## CI Status
