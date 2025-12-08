# AWS Free Tier Deployment Guide

This guide will help you deploy the Markdown Circus Docs (MCD) application on AWS EC2 free tier.

## Prerequisites

- AWS account with free tier eligibility
- Basic knowledge of Linux command line
- SSH access to EC2 instance

## Step 1: Launch EC2 Instance

1. **Log in to AWS Console** and navigate to EC2
2. **Launch Instance** with the following settings:
   - **AMI**: Ubuntu Server 22.04 LTS (Free Tier eligible)
   - **Instance Type**: t3.micro (Free Tier eligible)
   - **Key Pair**: Create or select an existing key pair
   - **Security Group**: Create a new security group with the following rules:
     - SSH (22) from your IP
     - HTTP (80) from anywhere (0.0.0.0/0)
     - HTTPS (443) from anywhere (0.0.0.0/0)
   - **Storage**: 8 GB gp3 (Free Tier includes 30 GB)

3. **Launch the instance**

## Step 2: Connect to EC2 Instance

```bash
ssh -i your-key.pem ubuntu@your-ec2-public-ip
```

## Step 3: Initial Setup

Run the setup script to install Docker and required tools:

```bash
# Clone the repository
git clone https://github.com/ukma-cs-ssdm-2025/team-circus.git
cd team-circus

# Make scripts executable
chmod +x scripts/setup-aws.sh
chmod +x deploy.sh
chmod +x health-check.sh

# Run setup script
./scripts/setup-aws.sh
```

**Important**: After running the setup script, logout and login again (or run `newgrp docker`) to apply Docker group changes.

## Step 4: Configure Environment Variables

1. **Copy the production environment template**:
   ```bash
   cp .env.production.sample .env
   ```

2. **Edit the .env file** with your production values:
   ```bash
   nano .env
   ```

3. **Generate secure passwords and tokens**:
   ```bash
   # Generate a strong database password
   openssl rand -base64 32
   
   # Generate a strong secret token
   openssl rand -base64 64
   ```

4. **Update the following critical values**:
   - `DB_PASSWORD`: Strong password for PostgreSQL
   - `SECRET_TOKEN`: Strong secret for JWT tokens
   - `CORS_ALLOW_ORIGINS`: Your domain (e.g., `http://your-ec2-ip` or `https://your-domain.com`)
   - `VITE_API_BASE_URL`: Your API URL (e.g., `http://your-ec2-ip/api` or `https://your-domain.com/api`)

## Step 5: Deploy the Application

```bash
./deploy.sh
```

This script will:
- Build Docker images
- Start all services (PostgreSQL, Backend, Frontend, Nginx)
- Run database migrations
- Set up the reverse proxy

## Step 6: Verify Deployment

1. **Check service status**:
   ```bash
   docker compose -f docker-compose.prod.yml ps
   ```

2. **Check logs**:
   ```bash
   # All services
   docker compose -f docker-compose.prod.yml logs -f
   
   # Specific service
   docker compose -f docker-compose.prod.yml logs -f backend
   ```

3. **Run health check**:
   ```bash
   ./health-check.sh
   ```

4. **Access the application**:
   - Open your browser and navigate to: `http://your-ec2-public-ip`
   - API should be available at: `http://your-ec2-public-ip/api`

## Step 7: (Optional) Set Up Domain and SSL

### Using a Domain Name

1. **Point your domain to EC2 IP**:
   - Add an A record in your DNS provider pointing to your EC2 public IP

2. **Update CORS and API URLs** in `.env`:
   ```
   CORS_ALLOW_ORIGINS=https://your-domain.com,https://www.your-domain.com
   VITE_API_BASE_URL=https://your-domain.com/api
   ```

3. **Set up SSL with Let's Encrypt** (recommended):
   ```bash
   # Install certbot
   sudo apt-get update
   sudo apt-get install -y certbot python3-certbot-nginx
   
   # Get certificate (adjust for your domain)
   sudo certbot certonly --standalone -d your-domain.com -d www.your-domain.com
   
   # Copy certificates to nginx/ssl directory
   sudo mkdir -p nginx/ssl
   sudo cp /etc/letsencrypt/live/your-domain.com/fullchain.pem nginx/ssl/
   sudo cp /etc/letsencrypt/live/your-domain.com/privkey.pem nginx/ssl/
   sudo chmod 644 nginx/ssl/*.pem
   ```

4. **Update nginx configuration**:
   - Uncomment the HTTPS server block in `nginx/conf.d/default.conf`
   - Update the `server_name` with your domain
   - Uncomment the HTTP to HTTPS redirect

5. **Restart services**:
   ```bash
   docker compose -f docker-compose.prod.yml restart nginx
   ```

## Maintenance Commands

### View Logs
```bash
# All services
docker compose -f docker-compose.prod.yml logs -f

# Specific service
docker compose -f docker-compose.prod.yml logs -f backend
docker compose -f docker-compose.prod.yml logs -f frontend
docker compose -f docker-compose.prod.yml logs -f postgres
```

### Restart Services
```bash
# All services
docker compose -f docker-compose.prod.yml restart

# Specific service
docker compose -f docker-compose.prod.yml restart backend
```

### Stop Services
```bash
docker compose -f docker-compose.prod.yml down
```

### Update Application
```bash
# Pull latest code
git pull

# Rebuild and restart
./deploy.sh
```

### Backup Database
```bash
docker compose -f docker-compose.prod.yml exec postgres pg_dump -U mcd_user mcd_db > backup_$(date +%Y%m%d_%H%M%S).sql
```

### Restore Database
```bash
docker compose -f docker-compose.prod.yml exec -T postgres psql -U mcd_user mcd_db < backup_file.sql
```

## Monitoring and Troubleshooting

### Check Resource Usage
```bash
# Docker stats
docker stats

# Disk usage
df -h
docker system df
```

### Common Issues

1. **Port already in use**:
   ```bash
   sudo lsof -i :80
   sudo kill -9 <PID>
   ```

2. **Database connection errors**:
   - Check if postgres container is running: `docker compose -f docker-compose.prod.yml ps postgres`
   - Check postgres logs: `docker compose -f docker-compose.prod.yml logs postgres`

3. **Out of memory**:
   - Free tier t3.micro has 1 GB RAM
   - Monitor with: `free -h`
   - Consider stopping unnecessary services

4. **Disk space**:
   - Clean Docker: `docker system prune -a`
   - Remove old images: `docker image prune -a`

## Security Considerations

1. **Change default passwords** in `.env`
2. **Use strong SECRET_TOKEN** (at least 64 characters)
3. **Restrict SSH access** to your IP only
4. **Enable firewall** (UFW is configured in setup script)
5. **Use HTTPS** in production (Let's Encrypt)
6. **Regular updates**: Keep system and Docker images updated
7. **Backup database** regularly

## Cost Optimization (Free Tier)

- **Instance**: Use t3.micro (750 hours/month free)
- **Storage**: Use 8-10 GB (30 GB free tier)
- **Data Transfer**: First 1 GB/month free
- **Monitor usage** in AWS Cost Explorer

## Next Steps

- Set up automated backups
- Configure monitoring (CloudWatch or similar)
- Set up CI/CD pipeline
- Configure log rotation
- Set up alerts for service failures

## Support

For issues or questions, please refer to:
- Project README: [README.md](../README.md)
- Backend documentation: [backend/README.md](../backend/README.md)
- Frontend documentation: [frontend/README.md](../frontend/README.md)

