#!/bin/bash

# Health check script for monitoring services

set -e

FRONTEND_URL="${FRONTEND_URL:-http://localhost}"
BACKEND_URL="${BACKEND_URL:-http://localhost/api/v1}"

check_service() {
    local name=$1
    local url=$2
    local endpoint=$3
    
    if curl -f -s "$url$endpoint" > /dev/null 2>&1; then
        echo "✓ $name is healthy"
        return 0
    else
        echo "✗ $name is unhealthy"
        return 1
    fi
}

echo "Checking service health..."

FRONTEND_OK=false
BACKEND_OK=false

# Check frontend (health endpoint)
if check_service "Frontend" "$FRONTEND_URL" "/health"; then
    FRONTEND_OK=true
fi

# Check backend (swagger endpoint as health check)
if check_service "Backend" "$BACKEND_URL" "/../swagger/index.html"; then
    BACKEND_OK=true
fi

if [ "$FRONTEND_OK" = true ] && [ "$BACKEND_OK" = true ]; then
    echo "All services are healthy"
    exit 0
else
    echo "Some services are unhealthy"
    exit 1
fi

