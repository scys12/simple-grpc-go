# simple-grpc-go

## Instructions

1. Install Go
2. Clone this project
3. Create database `todo` and copy `db.sql` for migrating table
4. Download dependency using command:

```go
go mod download
```

5. Rename `.exampe.env` to `config.env`

6. Build application using `go build`

7. Run with

```go
./simple-grpc-go
```

8. Access through `localhost:9000/`

If you want to compile new protobuf, take a look at `build-todo-rpc` in `Makefile`. [protoc](https://grpc.io/docs/protoc-installation/) will compile it into `go` file.

Read documentation of gRPC [here](https://grpc.io/) and how to serve gRPC using http proto [here](https://grpc-ecosystem.github.io/grpc-gateway/)
