FROM golang:1.23.2

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/web

CMD ["./main"]
