FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor/ ./vendor

COPY . .

RUN GOFLAGS=-mod=vendor go build -o main ./cmd

FROM alpine:latest

COPY --from=builder /app/main /app/main

ENTRYPOINT ["/app/main"]