FROM golang:1.16 AS build

# Environment variables will be put here
ENV GO111MODULE=on
# end

RUN apt-get update && \
  apt-get upgrade -y && \
  apt-get install -y git

WORKDIR /go/src/simple-grpc-go/

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main .

FROM alpine
COPY --from=build /go/src/simple-grpc-go/config.env go/bin/simple-grpc-go/
COPY --from=build /go/src/simple-grpc-go/main go/bin/simple-grpc-go/

EXPOSE 9000
WORKDIR /go/bin/simple-grpc-go/

CMD ["./main"]