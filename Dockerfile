# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git build-base gcc musl-dev

# Copy go module files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build all binaries
RUN CGO_ENABLED=1 GOOS=linux go build -o /drcp-api ./cmd/api
RUN CGO_ENABLED=1 GOOS=linux go build -o /drcp-worker ./cmd/worker
RUN CGO_ENABLED=1 GOOS=linux go build -o /drcp-xds ./cmd/xds
RUN CGO_ENABLED=1 GOOS=linux go build -o /drcp-anchor ./cmd/anchor

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binaries from builder
COPY --from=builder /drcp-api /app/drcp-api
COPY --from=builder /drcp-worker /app/drcp-worker
COPY --from=builder /drcp-xds /app/drcp-xds
COPY --from=builder /drcp-anchor /app/drcp-anchor

# Copy web frontend assets
COPY --from=builder /app/web /app/web

# Expose port
EXPOSE 8080

# Run
CMD ["/app/drcp-api"]
