FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY pkg/ ./pkg/
COPY gen/ ./gen/
COPY services/issue-service/ ./services/issue-service/

WORKDIR /app/services/issue-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /issue-service ./cmd/server

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /issue-service .
COPY --from=builder /app/services/issue-service/migrations ./services/issue-service/migrations

EXPOSE 50053
CMD ["./issue-service"]
