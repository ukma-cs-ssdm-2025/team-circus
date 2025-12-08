#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting deployment...${NC}"

# Check if .env file exists
if [ ! -f .env ]; then
    echo -e "${YELLOW}Warning: .env file not found.${NC}"
    echo -e "${YELLOW}Creating .env from .env.production.sample...${NC}"
    if [ -f .env.production.sample ]; then
        cp .env.production.sample .env
        echo -e "${RED}Please edit .env file with your production values before continuing!${NC}"
        exit 1
    else
        echo -e "${RED}Error: .env.production.sample not found.${NC}"
        exit 1
    fi
fi

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Error: Docker is not installed.${NC}"
    exit 1
fi

# Check if docker-compose is installed
if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo -e "${RED}Error: docker-compose is not installed.${NC}"
    exit 1
fi

# Use docker compose (v2) if available, otherwise docker-compose (v1)
if docker compose version &> /dev/null; then
    COMPOSE_CMD="docker compose"
else
    COMPOSE_CMD="docker-compose"
fi

# Stop existing containers
echo -e "${GREEN}Stopping existing containers...${NC}"
$COMPOSE_CMD -f docker-compose.prod.yml down || true

# Pull latest images (if using pre-built images)
# echo -e "${GREEN}Pulling latest images...${NC}"
# $COMPOSE_CMD -f docker-compose.prod.yml pull

# Build and start containers
echo -e "${GREEN}Building and starting containers...${NC}"
$COMPOSE_CMD -f docker-compose.prod.yml up -d --build

# Wait for services to be healthy
echo -e "${GREEN}Waiting for services to be healthy...${NC}"
sleep 10

# Check service health
echo -e "${GREEN}Checking service health...${NC}"
$COMPOSE_CMD -f docker-compose.prod.yml ps

echo -e "${GREEN}Deployment completed!${NC}"
echo -e "${YELLOW}Services should be available at:${NC}"
echo -e "  - Frontend: http://$(hostname -I | awk '{print $1}')"
echo -e "  - API: http://$(hostname -I | awk '{print $1}')/api"

