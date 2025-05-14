FROM golang:1.23-bullseye AS builder


RUN apt-get update && apt-get install -y \
    gcc \
    librdkafka-dev \
    git \
    pkg-config \
    libc6-dev \
    make \
    ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary
WORKDIR /app/cmd/web
RUN go build -o /app/bin/web

WORKDIR /app/cmd/worker
RUN go build -o /app/bin/worker

# Runtime stage
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y librdkafka1 ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/web /web
COPY --from=builder /app/bin/worker /worker

CMD ["/web"]
