# schema_service

This service is accountable for the creation and maintenance of the schemas.

## Dependencies

In order to generate the interfaces using `proto`, you will need to install some previous dependencies:

```bash
TODO
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
