# syntax=docker/dockerfile:1

FROM golang:1.21.11-alpine3.20
WORKDIR /build

COPY worker  ./worker
COPY shared  ./shared
COPY grpc  ./grpc
COPY go.mod ./go.mod
COPY go.sum ./go.sum
COPY data ./data

RUN go mod download
RUN go build -v -o worker.run ./worker
CMD [ "/build/worker.run" ]