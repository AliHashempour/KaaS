FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -v -o kaas ./cmd/kaas

FROM alpine:3.12

WORKDIR /app

COPY --from=builder /app/kaas ./

CMD ["./kaas"]
