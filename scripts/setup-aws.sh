#!/bin/bash

# AWS EC2 setup script for Ubuntu/Debian
# Run this script on a fresh EC2 instance

set -e

echo "Setting up AWS EC2 instance for deployment..."

# Update system
sudo apt-get update
sudo apt-get upgrade -y

# Install Docker
if ! command -v docker &> /dev/null; then
    echo "Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
    rm get-docker.sh
    echo "Docker installed successfully"
else
    echo "Docker is already installed"
fi

# Install Docker Compose
if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo "Installing Docker Compose..."
    sudo apt-get install -y docker-compose-plugin
    echo "Docker Compose installed successfully"
else
    echo "Docker Compose is already installed"
fi

# Install Git
if ! command -v git &> /dev/null; then
    echo "Installing Git..."
    sudo apt-get install -y git
    echo "Git installed successfully"
else
    echo "Git is already installed"
fi

# Install additional utilities
sudo apt-get install -y curl wget vim

# Configure firewall (UFW)
echo "Configuring firewall..."
sudo ufw allow 22/tcp   # SSH
sudo ufw allow 80/tcp   # HTTP
sudo ufw allow 443/tcp  # HTTPS
sudo ufw --force enable

echo "Setup completed!"
echo ""
echo "Next steps:"
echo "1. Logout and login again (or run: newgrp docker) to apply docker group changes"
echo "2. Clone your repository: git clone <your-repo-url>"
echo "3. Copy .env.production.sample to .env and configure it"
echo "4. Run ./deploy.sh to deploy the application"

