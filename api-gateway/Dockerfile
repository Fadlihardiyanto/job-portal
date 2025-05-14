FROM golang:1.23-bullseye AS builder

RUN apt-get update && apt-get install -y git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/web
RUN go build -o /app/bin/web

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/web /web

CMD ["/web"]
