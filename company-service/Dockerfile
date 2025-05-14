
FROM golang:1.23-alpine AS builder


RUN apk add --no-cache git librdkafka-dev gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/web
RUN go build -o /app/bin/web

WORKDIR /app/cmd/worker
RUN go build -o /app/bin/worker

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/web /web
COPY --from=builder /app/bin/worker /worker

CMD ["/web"]
