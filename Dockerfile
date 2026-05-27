# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go module files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the API binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /drcp-api ./cmd/api

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /drcp-api /app/drcp-api

# Copy web frontend assets
COPY --from=builder /app/web /app/web

# Expose port
EXPOSE 8080

# Run
CMD ["/app/drcp-api"]
