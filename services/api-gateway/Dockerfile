FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY pkg/ ./pkg/
COPY services/api-gateway/ ./services/api-gateway/

WORKDIR /app/services/api-gateway
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-gateway ./cmd/server

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /api-gateway .

EXPOSE 8080
CMD ["./api-gateway"]
