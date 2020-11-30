[![Go Report Card](https://goreportcard.com/badge/github.com/vontikov/stoa)](https://goreportcard.com/report/github.com/vontikov/stoa)

# Stoa

Distributed in-memory collections based on
[Hashicorp's Raft implementation](https://github.com/hashicorp/raft)

## Build

1. Install [Go](https://golang.org/doc/install)

2. Make sure you have GOPATH environment variable defined in your system.

3. Add $GOPATH/bin to the PATH environment variable.

4. Run `make deps` to install the following dependencies:

* [protoc](https://grpc.io/docs/protoc-installation/)
* [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
* [protoc-gen-validate](https://github.com/envoyproxy/protoc-gen-validate)
* [gomock](https://github.com/golang/mock)

5. `make build` builds Stoa executable

6. `make image` builds Stoa Docker image
