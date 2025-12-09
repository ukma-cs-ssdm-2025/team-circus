FROM golang:1.25.1 AS builder
WORKDIR /app

# Install task runner
RUN go install github.com/go-task/task/v3/cmd/task@latest

# Copy go mod files first for better caching
COPY backend/go.mod backend/go.sum ./
WORKDIR /app/backend

# Download dependencies
RUN go mod download

# Copy source code
COPY backend/ .

# Ensure vendor directory is up-to-date after copying source
RUN go mod vendor

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ./bin/backend cmd/main.go

# Install migrator
RUN task install:migrator

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
COPY --from=builder /app/backend/taskfile.yml ./taskfile.yml
COPY --from=builder /app/backend/bin/backend ./bin/backend
COPY --from=builder /go/bin/task /usr/local/bin/task
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY backend/db/migrations ./db/migrations

# Create non-root user for security
RUN useradd -m -u 1000 appuser && \
    chown -R appuser:appuser /app

USER appuser

# Health check (using swagger endpoint as health check)
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:${API_PORT:-8080}/swagger/index.html || exit 1
