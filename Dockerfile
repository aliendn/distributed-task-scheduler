# syntax=docker/dockerfile:1

# Build Stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Add build dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o distributed-task-scheduler ./cmd/distributed-task-scheduler

# Final Image
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/distributed-task-scheduler .

EXPOSE 8080

ENTRYPOINT ["./distributed-task-scheduler"]
