# Deployment Quick Start

This project is now ready for deployment on AWS EC2 free tier.

## Quick Deployment Steps

1. **Launch an EC2 instance** (Ubuntu 22.04, t3.micro)
2. **Connect via SSH**:
   ```bash
   ssh -i your-key.pem ubuntu@your-ec2-ip
   ```
3. **Clone and setup**:
   ```bash
   git clone https://github.com/ukma-cs-ssdm-2025/team-circus.git
   cd team-circus
   chmod +x scripts/setup-aws.sh deploy.sh health-check.sh
   ./scripts/setup-aws.sh
   ```
4. **Configure environment**:
   ```bash
   cp env.production.sample .env
   nano .env  # Edit with your values
   ```
5. **Deploy**:
   ```bash
   ./deploy.sh
   ```

## Files Created for Deployment

- `docker-compose.prod.yml` - Production Docker Compose configuration
- `frontend/Dockerfile` - Frontend production build
- `nginx/` - Nginx reverse proxy configuration
- `deploy.sh` - Deployment script
- `health-check.sh` - Health check script
- `scripts/setup-aws.sh` - AWS EC2 setup script
- `env.sample` / `env.production.sample` - Environment variable templates

## Detailed Documentation

For complete deployment instructions, see:
- **[AWS Deployment Guide](docs/deployment/AWS_DEPLOYMENT.md)** - Complete step-by-step guide

## Architecture

The production deployment consists of:
- **PostgreSQL** - Database (internal network only)
- **Backend** - Go API server (internal network only)
- **Frontend** - React application (served via Nginx)
- **Nginx** - Reverse proxy and static file server (exposed on ports 80/443)

## Environment Variables

Required environment variables are documented in:
- `env.sample` - Development template
- `env.production.sample` - Production template

**Important**: Always change default passwords and secrets in production!

## Health Checks

All services include health checks:
- Backend: `/api/v1/health`
- Frontend: `/health`
- Nginx: `/health`

Run the health check script:
```bash
./health-check.sh
```

## Maintenance

### View Logs
```bash
docker compose -f docker-compose.prod.yml logs -f
```

### Restart Services
```bash
docker compose -f docker-compose.prod.yml restart
```

### Update Application
```bash
git pull
./deploy.sh
```

## Security Notes

- Change all default passwords in `.env`
- Use strong `SECRET_TOKEN` (64+ characters)
- Configure firewall (UFW) properly
- Use HTTPS in production (Let's Encrypt)
- Restrict SSH access to your IP

## Support

For issues, refer to the [AWS Deployment Guide](docs/deployment/AWS_DEPLOYMENT.md) or project documentation.

