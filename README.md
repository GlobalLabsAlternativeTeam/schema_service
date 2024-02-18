# schema_service

This service is accountable for the creation and maintenance of the schemas.

## Dependencies

- **Go**, any one of the three latest major releases of Go ([installation guide](https://go.dev/doc/install)).
- **Protocol buffer compiler**, `protoc`, version 3 ([installation guide](https://grpc.io/docs/protoc-installation/)).
- **Go plugins** for the protocol compiler:

  1. Install the protocol compiler plugins for Go using the following commands:

  ```bash
  $ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
  $ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
  ```

  2. Update your PATH so that the protoc compiler can find the plugins:

  ```bash
  $ export PATH="$PATH:$(go env GOPATH)/bin"
  ```

## Usage

### Generating proto interfaces

In order to run the service, we first need to **generate the interfaces** using protoc.

To do this, run the following command from the root directory:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/schema_service.proto
```

### Running the service

You can now launch the service using:

```bash
go run cmd/main.go
```

### Extra: generating example data

We provide a script to generate some example data located in `~/cmd/scripts/gen_data.go`.
To execute this, run the following command from the root directory:

```bash
go run cmd/scripts/gen_data.go
```
