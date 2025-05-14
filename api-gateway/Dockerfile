FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/web
RUN go build -o /app/bin/web

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/web /web

CMD ["/web"]
