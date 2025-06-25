# Deployment Guide

This guide covers how to deploy org-roam-woven using the automatically built Docker images.

## 🏗️ CI/CD Pipeline Overview

The project uses GitHub Actions to automatically:

1. **Test** the Go application
2. **Build** Docker images for multiple architectures (amd64, arm64)  
3. **Push** images to GitHub Container Registry (GHCR)
4. **Tag** images based on branches and version tags

### Trigger Conditions

- **Pull Requests**: Run tests only
- **Push to main**: Build and push `latest` tags
- **Version tags (v*)**: Build and push versioned releases

## 📦 Available Images

Images are hosted on GitHub Container Registry:

```bash
# Main application
ghcr.io/your-org/org-roam-woven/org-roam-woven:latest
ghcr.io/your-org/org-roam-woven/org-roam-woven:v1.0.0

# Emacs service  
ghcr.io/your-org/org-roam-woven/emacs:latest
ghcr.io/your-org/org-roam-woven/emacs:v1.0.0
```

### Image Tags

| Tag Pattern | Description | Example |
|-------------|-------------|---------|
| `latest` | Latest stable from main branch | `latest` |
| `v*` | Semantic version releases | `v1.2.3` |
| `main-*` | Commit-specific builds | `main-abc1234` |

## 🚀 Quick Deployment

### Using Docker Compose (Recommended)

1. **Download compose file**
   ```bash
   wget https://raw.githubusercontent.com/your-org/org-roam-woven/main/docker-compose.yml
   wget https://raw.githubusercontent.com/your-org/org-roam-woven/main/.env.example
   cp .env.example .env
   ```

2. **Configure environment**
   ```bash
   # Edit .env file
   GITHUB_REPOSITORY=your-org/org-roam-woven
   TAG=latest
   PORT=18080
   ```

3. **Deploy**
   ```bash
   docker-compose up -d
   ```

### Using Docker Run

```bash
# Create network and volume
docker network create org-roam-network
docker volume create org_roam_data

# Run Emacs service
docker run -d \
  --name org-roam-emacs \
  --network org-roam-network \
  -v org_roam_data:/root/org-roam \
  ghcr.io/your-org/org-roam-woven/emacs:latest \
  emacs --fg-daemon

# Run main application
docker run -d \
  --name org-roam-woven-app \
  --network org-roam-network \
  -p 18080:18080 \
  -v org_roam_data:/root/org-roam:ro \
  -e GIN_MODE=release \
  ghcr.io/your-org/org-roam-woven/org-roam-woven:latest
```

## 🔄 Updates and Rollbacks

### Update to Latest

```bash
# Pull latest images
docker-compose pull

# Restart services  
docker-compose up -d
```

### Deploy Specific Version

```bash
# Set version in environment
export TAG=v1.2.3

# Deploy
docker-compose up -d
```

### Rollback

```bash
# Deploy previous version
export TAG=v1.1.0
docker-compose up -d
```

## 🔧 Production Considerations

### Resource Limits

Add resource limits to docker-compose.yml:

```yaml
services:
  org-roam-woven:
    deploy:
      resources:
        limits:
          memory: 256M
          cpus: '0.5'
```

### Health Checks

The application runs on port 18080. Check health:

```bash
curl http://localhost:18080/
```

### Persistent Data

Ensure org-roam data persists by using named volumes:

```yaml
volumes:
  org_roam_data:
    driver: local
```

### Logs

View application logs:

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f org-roam-woven
```

## 🛡️ Security Notes

- Images run as non-root user for security
- Based on minimal distroless images
- No unnecessary tools or packages included
- Multi-architecture support (amd64, arm64)

## 📊 Monitoring

### Service Status

```bash
docker-compose ps
```

### Resource Usage

```bash
docker stats
```

### Container Health

```bash
docker inspect --format='{{.State.Health.Status}}' org-roam-woven-app
```

## 🐛 Troubleshooting

### Common Issues

**Port already in use:**
```bash
# Change port in .env
PORT=18081
```

**Permission denied:**
```bash
# Check if user can access Docker
docker ps
```

**Image pull fails:**
```bash
# Login to GHCR (if repository is private)
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
```

**Service won't start:**
```bash
# Check logs
docker-compose logs org-roam-woven
```

### Getting Help

1. Check the [GitHub Issues](https://github.com/your-org/org-roam-woven/issues)
2. Review the application logs
3. Verify Docker and Docker Compose versions

## 📋 Deployment Checklist

- [ ] Docker and Docker Compose installed
- [ ] Environment variables configured
- [ ] Network ports available (18080)
- [ ] Persistent storage configured
- [ ] Images pulled successfully
- [ ] Services started and healthy
- [ ] Application accessible via HTTP

---

For development setup and contributing, see [README.md](README.md).