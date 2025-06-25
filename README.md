# org-roam-woven

[![CI/CD](https://github.com/your-org/org-roam-woven/actions/workflows/ci.yaml/badge.svg)](https://github.com/your-org/org-roam-woven/actions/workflows/ci.yaml)

A web service for org-roam integration, built with Go and Gin framework.

## 🚀 Quick Start

### Using Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/your-org/org-roam-woven.git
   cd org-roam-woven
   ```

2. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env if needed
   ```

3. **Start services**
   ```bash
   docker-compose up -d
   ```

4. **Access the application**
   ```bash
   curl http://localhost:18080
   ```

### Using Pre-built Images

```bash
# Pull and run the latest image
docker run -d \
  --name org-roam-woven \
  -p 18080:18080 \
  ghcr.io/your-org/org-roam-woven/org-roam-woven:latest
```

## 🛠️ Development

### Prerequisites

- Go 1.24.4+
- Docker and Docker Compose

### Local Development

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Run locally**
   ```bash
   go run main.go
   ```

3. **Run tests**
   ```bash
   go test ./...
   ```

### Building Docker Images

```bash
# Build application image
docker build -f packaging/Dockerfile.org-roam-woven -t org-roam-woven:local .

# Build emacs image
docker build -f packaging/Dockerfile.emacs -t org-roam-emacs:local ./packaging
```

## 📦 CI/CD Pipeline

The project uses GitHub Actions for automated CI/CD:

- **On Pull Request**: Runs tests and code quality checks
- **On Push to Main**: Builds and pushes Docker images to GHCR
- **On Version Tags**: Creates releases with versioned images

### Docker Images

Images are automatically built and pushed to GitHub Container Registry:

- `ghcr.io/your-org/org-roam-woven/org-roam-woven:latest`
- `ghcr.io/your-org/org-roam-woven/emacs:latest`

### Supported Tags

- `latest` - Latest stable release
- `main` - Latest from main branch
- `v1.0.0` - Specific version tags
- `main-abc1234` - Commit-specific builds

## 🔧 Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `18080` | HTTP server port |
| `GIN_MODE` | `debug` | Gin framework mode |
| `GITHUB_REPOSITORY` | - | GitHub repository name |
| `TAG` | `latest` | Docker image tag |

### Docker Compose

The `docker-compose.yml` includes:
- **org-roam-woven**: Main application service
- **emacs**: Emacs daemon for org-roam integration
- **Shared volume**: For org-roam data persistence

## 📋 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Welcome message |

## 🚢 Deployment

### Production Deployment

1. **Set environment variables**
   ```bash
   export GITHUB_REPOSITORY=your-org/org-roam-woven
   export TAG=v1.0.0  # or latest
   ```

2. **Deploy with Docker Compose**
   ```bash
   docker-compose up -d
   ```

3. **Check status**
   ```bash
   docker-compose ps
   curl http://localhost:18080
   ```

### Updating

```bash
# Pull latest images
docker-compose pull

# Restart services
docker-compose up -d
```

## 🔍 Monitoring

### Check Service Status
```bash
docker-compose ps
```

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f org-roam-woven
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Run tests: `go test ./...`
5. Commit: `git commit -m "feat: add my feature"`
6. Push: `git push origin feature/my-feature`
7. Create a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [org-roam](https://github.com/org-roam/org-roam)
- [Docker](https://www.docker.com/)