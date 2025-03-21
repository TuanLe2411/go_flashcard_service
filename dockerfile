# Build stage
FROM golang:1.24.0-alpine AS builder

# Set necessary environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

# Add labels for better maintenance
LABEL maintainer="Flashcard Service Team" \
      description="Flashcard Service API" \
      version="1.0"

WORKDIR /app

# Copy go mod and sum files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application with optimizations
RUN go build -ldflags="-s -w" -o flashcard_service cmd/main/main.go

# Run stage
FROM alpine:3.21

# Add necessary runtime packages and security configurations
RUN apk add --no-cache ca-certificates tzdata curl && \
    addgroup -S appgroup && adduser -S appuser -G appgroup && \
    mkdir -p /app/log && chown appuser:appgroup /app/log && chmod 777 /app/log

# Set the timezone to Asia/Ho_Chi_Minh
ENV TZ=Asia/Ho_Chi_Minh
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app

# Copy binary from build stage
COPY --from=builder /app/flashcard_service .

# Copy environment files
COPY .env.* ./

# Use non-root user
USER appuser

# Expose the port the service runs on
EXPOSE 9090

ENV ENV=production

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:9090/health || exit 1

# Command to run
ENTRYPOINT ["./flashcard_service"]