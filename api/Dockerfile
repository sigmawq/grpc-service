# syntax=docker/dockerfile:1

FROM golang:1.21.11-alpine3.20
WORKDIR /build

COPY api  ./api
COPY shared  ./shared
COPY grpc  ./grpc
COPY go.mod ./go.mod
COPY go.sum ./go.sum

RUN go mod download
RUN go build -v -o api.run ./api
EXPOSE 8080
CMD [ "/build/api.run" ]