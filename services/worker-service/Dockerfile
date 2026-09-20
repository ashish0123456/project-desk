FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY pkg/ ./pkg/
COPY services/worker-service/ ./services/worker-service/

WORKDIR /app/services/worker-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /worker-service ./cmd/worker

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /worker-service .

CMD ["./worker-service"]
