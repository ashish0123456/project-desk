FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY pkg/ ./pkg/
COPY gen/ ./gen/
COPY services/auth-service/ ./services/auth-service/

WORKDIR /app/services/auth-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service ./cmd/server

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /auth-service .
COPY --from=builder /app/services/auth-service/migrations ./services/auth-service/migrations

EXPOSE 50051
CMD ["./auth-service"]
