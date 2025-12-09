FROM golang:1.23 AS builder
WORKDIR /app/backend

# Use vendored dependencies to avoid network fetches in container builds
ENV GOFLAGS=-mod=vendor

# Copy sources and vendor tree
COPY backend/go.mod backend/go.sum ./
COPY backend/vendor ./vendor
COPY backend/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ./bin/backend cmd/main.go

# Production stage
FROM debian:bookworm-slim
WORKDIR /app

# Install only essential packages
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    wget && \
    rm -rf /var/lib/apt/lists/*

# Copy necessary files from builder
COPY --from=builder /app/backend/bin/backend ./bin/backend
COPY backend/db/migrations ./db/migrations

# Create non-root user for security
RUN useradd -m -u 1000 appuser && \
    chown -R appuser:appuser /app

USER appuser

# Health check (using swagger endpoint as health check)
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD wget --no-verbose --tries=1 -O /dev/null http://localhost:8080/swagger/index.html || exit 1
