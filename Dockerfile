
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o main ./cmd/api

FROM alpine:latest


WORKDIR /app

COPY --from=builder /app/main .

RUN adduser -D appuser
USER appuser

CMD ["./main"]
