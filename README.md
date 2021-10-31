# simple-grpc-go

## Instructions

1. Install Go & Docker
2. Clone this project
3. Rename `.exampe.env` to `config.env`
4. Build the project with

```
docker-compose up -d
```

5. Access through `localhost:9000/`

If you want to rebuild the project, use

```
docker-compose -d --build
```

To destroy the project, use

```
docker-compose down -v
```

For API Documentation, you can look at [here](https://app.swaggerhub.com/apis-docs/scys12/Todo_API/1.0)

If you want to compile new protobuf, take a look at `build-todo-rpc` in `Makefile`. [protoc](https://grpc.io/docs/protoc-installation/) will compile it into `go` file.

Read documentation of gRPC [here](https://grpc.io/) and how to serve gRPC using http proto [here](https://grpc-ecosystem.github.io/grpc-gateway/)

This project use `grafana` and `prometheus` for system monitoring and alerting. Read official documentation [here](https://prometheus.io/) about prometheus and [here](https://grafana.com) for grafana. There are 3 metrics that are measured, namely latency, total request and response status. Read more [here](https://github.com/scys12/simple-grpc-go/tree/master/pkg/monitoring).

This project also use `jaeger` for distributed tracing system. Read more documentation about [jaeger](https://www.jaegertracing.io/)

TODO :

- Add interceptor for gRPC
- Log, Monitoring, Tracing for gRPC
