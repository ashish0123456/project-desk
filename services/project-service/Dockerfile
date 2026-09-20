FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY pkg/ ./pkg/
COPY gen/ ./gen/
COPY services/project-service/ ./services/project-service/

WORKDIR /app/services/project-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /project-service ./cmd/server

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /project-service .
COPY --from=builder /app/services/project-service/migrations ./services/project-service/migrations

EXPOSE 50052
CMD ["./project-service"]
