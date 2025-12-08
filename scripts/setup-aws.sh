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

# Install Task (taskfile.dev)
if ! command -v task &> /dev/null; then
    echo "Installing Task..."
    sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.local/bin
    export PATH="$HOME/.local/bin:$PATH"
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
    echo "Task installed successfully"
else
    echo "Task is already installed"
fi

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
echo "2. Source ~/.bashrc or logout/login to add Task to PATH"
echo "3. Clone your repository: git clone <your-repo-url>"
echo "4. Run 'task copy:env' to create .env file and configure it"
echo "5. Run 'task docker:prod:deploy' to deploy the application"

