# Build stage
FROM golang:1.24.0-alpine AS builder

# Set necessary environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm64

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

# Add necessary packages
RUN apk --no-cache add ca-certificates tzdata && \
    mkdir /app

# Set timezone to UTC by default
ENV TZ=UTC

# Create a non-root user
RUN adduser -D -g '' appuser && \
    chown -R appuser:appuser /app

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