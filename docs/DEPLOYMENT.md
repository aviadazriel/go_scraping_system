# Deployment Guide

This guide covers deploying the Go Scraping Project in different environments.

## 🏠 Local Development

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL (running locally on port 5432)

### Quick Deployment

1. **Setup Database:**
   ```bash
   # Create database user and database
   psql -h localhost -U aazriel -d postgres -c "CREATE USER scraper WITH PASSWORD 'scraper_password';"
   psql -h localhost -U aazriel -d postgres -c "CREATE DATABASE scraper OWNER scraper;"
   
   # Run migrations
   DATABASE_URL="postgres://scraper:scraper_password@localhost:5432/scraper?sslmode=disable" make migrate-up
   ```

2. **Start Services:**
   ```bash
   docker-compose -f docker-compose.local.yml up -d
   ```

3. **Verify Deployment:**
   ```bash
   # Check service status
   docker-compose -f docker-compose.local.yml ps
   
   # Test API
   curl http://localhost:8082/health
   ```

### Manual Deployment

If you prefer to run services individually:

```bash
# Build services
make build-service SERVICE=url-manager
make build-service SERVICE=api-gateway

# Run services locally
make run-service SERVICE=url-manager &
make run-service SERVICE=api-gateway &
```

## 🐳 Docker Production

### Using Production Docker Compose

```bash
# Start all services including PostgreSQL
docker-compose -f docker-compose.production.yml up -d

# Check status
docker-compose -f docker-compose.production.yml ps
```

### Custom Docker Deployment

1. **Build Images:**
   ```bash
   # Build all images
   docker build -f cmd/url-manager/Dockerfile -t url-manager:latest .
   docker build -f cmd/api-gateway/Dockerfile -t api-gateway:latest .
   ```

2. **Run Containers:**
   ```bash
   # Start infrastructure
   docker-compose -f docker-compose.local.yml up -d zookeeper kafka kafka-ui
   
   # Run services
   docker run -d --name url-manager \
     --network scraping_network \
     -e SCRAPING_DATABASE_HOST=host.docker.internal \
     -e SCRAPING_DATABASE_PORT=5432 \
     -e SCRAPING_DATABASE_USER=scraper \
     -e SCRAPING_DATABASE_PASSWORD=scraper_password \
     -e SCRAPING_DATABASE_NAME=scraper \
     -e SCRAPING_KAFKA_BROKERS=kafka:29092 \
     url-manager:latest
   
   docker run -d --name api-gateway \
     --network scraping_network \
     -p 8082:8080 \
     -e DB_HOST=host.docker.internal \
     -e DB_PORT=5432 \
     -e DB_USER=scraper \
     -e DB_PASSWORD=scraper_password \
     -e DB_NAME=scraper \
     -e URL_MANAGER_URL=http://url-manager:8080 \
     api-gateway:latest
   ```

## ☁️ Cloud Deployment

### AWS ECS

1. **Create ECR Repository:**
   ```bash
   aws ecr create-repository --repository-name go-scraping-project
   ```

2. **Push Images:**
   ```bash
   # Login to ECR
   aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-west-2.amazonaws.com
   
   # Tag and push
   docker tag url-manager:latest <account-id>.dkr.ecr.us-west-2.amazonaws.com/go-scraping-project:url-manager
   docker tag api-gateway:latest <account-id>.dkr.ecr.us-west-2.amazonaws.com/go-scraping-project:api-gateway
   
   docker push <account-id>.dkr.ecr.us-west-2.amazonaws.com/go-scraping-project:url-manager
   docker push <account-id>.dkr.ecr.us-west-2.amazonaws.com/go-scraping-project:api-gateway
   ```

3. **Deploy with ECS:**
   ```bash
   # Create ECS cluster
   aws ecs create-cluster --cluster-name scraping-cluster
   
   # Deploy services (use AWS Console or CLI)
   ```

### Google Cloud Run

1. **Build and Push:**
   ```bash
   # Build images
   docker build -f cmd/url-manager/Dockerfile -t gcr.io/<project-id>/url-manager .
   docker build -f cmd/api-gateway/Dockerfile -t gcr.io/<project-id>/api-gateway .
   
   # Push to GCR
   docker push gcr.io/<project-id>/url-manager
   docker push gcr.io/<project-id>/api-gateway
   ```

2. **Deploy:**
   ```bash
   # Deploy to Cloud Run
   gcloud run deploy url-manager \
     --image gcr.io/<project-id>/url-manager \
     --platform managed \
     --region us-central1 \
     --allow-unauthenticated
   
   gcloud run deploy api-gateway \
     --image gcr.io/<project-id>/api-gateway \
     --platform managed \
     --region us-central1 \
     --allow-unauthenticated
   ```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SCRAPING_DATABASE_HOST` | PostgreSQL host | `localhost` |
| `SCRAPING_DATABASE_PORT` | PostgreSQL port | `5432` |
| `SCRAPING_DATABASE_USER` | Database user | `scraper` |
| `SCRAPING_DATABASE_PASSWORD` | Database password | `scraper_password` |
| `SCRAPING_DATABASE_NAME` | Database name | `scraper` |
| `SCRAPING_KAFKA_BROKERS` | Kafka brokers | `localhost:9092` |
| `SCRAPING_KAFKA_GROUP_ID` | Kafka group ID | `scraper-group` |
| `SCRAPING_SERVER_PORT` | HTTP server port | `8080` |
| `SCRAPING_LOG_LEVEL` | Log level | `info` |

### Secrets Management

For production, use proper secrets management:

**Docker Secrets:**
```bash
# Create secrets
echo "scraper_password" | docker secret create db_password -
echo "kafka_password" | docker secret create kafka_password -

# Use in docker-compose
services:
  url-manager:
    secrets:
      - db_password
      - kafka_password
```

**AWS Secrets Manager:**
```bash
# Store secrets
aws secretsmanager create-secret \
  --name /scraping/database \
  --secret-string '{"password":"scraper_password","user":"scraper"}'
```

## 📊 Monitoring

### Health Checks

All services expose health endpoints:
- API Gateway: `http://localhost:8082/health`
- URL Manager: Background service (no HTTP)

### Logging

```bash
# View logs
docker-compose -f docker-compose.local.yml logs -f

# Specific service
docker-compose -f docker-compose.local.yml logs -f url-manager
```

### Metrics

The services expose metrics endpoints:
- API Gateway: `http://localhost:8082/api/v1/metrics`
- Scraping stats: `http://localhost:8082/api/v1/metrics/scraping`

### Kafka Monitoring

- **Kafka UI**: `http://localhost:8080`
- **Kafka CLI**:
  ```bash
  # List topics
  docker exec scraping_kafka kafka-topics --bootstrap-server localhost:9092 --list
  
  # Monitor messages
  docker exec scraping_kafka kafka-console-consumer \
    --bootstrap-server localhost:9092 \
    --topic scraping-tasks \
    --from-beginning
  ```

## 🔒 Security

### Network Security

```bash
# Create custom network
docker network create scraping-network

# Use internal network for service communication
services:
  url-manager:
    networks:
      - scraping-network
  api-gateway:
    networks:
      - scraping-network
    ports:
      - "8082:8080"  # Only expose necessary ports
```

### SSL/TLS

For production, add SSL termination:

```bash
# Using nginx as reverse proxy
docker run -d --name nginx-proxy \
  -p 443:443 \
  -v /path/to/ssl:/etc/nginx/ssl \
  nginx:alpine
```

### Authentication

Add authentication middleware to the API Gateway:

```go
// Example middleware
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Validate JWT token or API key
        token := r.Header.Get("Authorization")
        if !validateToken(token) {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

## 🚀 Scaling

### Horizontal Scaling

```bash
# Scale API Gateway
docker-compose -f docker-compose.local.yml up -d --scale api-gateway=3

# Scale URL Manager
docker-compose -f docker-compose.local.yml up -d --scale url-manager=2
```

### Load Balancing

Use nginx or HAProxy for load balancing:

```nginx
upstream api_gateway {
    server api-gateway-1:8080;
    server api-gateway-2:8080;
    server api-gateway-3:8080;
}

server {
    listen 80;
    location / {
        proxy_pass http://api_gateway;
    }
}
```

## 🔄 CI/CD

### GitHub Actions

```yaml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Build and push images
        run: |
          docker build -f cmd/url-manager/Dockerfile -t ${{ secrets.ECR_REGISTRY }}/url-manager:${{ github.sha }} .
          docker build -f cmd/api-gateway/Dockerfile -t ${{ secrets.ECR_REGISTRY }}/api-gateway:${{ github.sha }} .
          docker push ${{ secrets.ECR_REGISTRY }}/url-manager:${{ github.sha }}
          docker push ${{ secrets.ECR_REGISTRY }}/api-gateway:${{ github.sha }}
      
      - name: Deploy to ECS
        run: |
          aws ecs update-service --cluster scraping-cluster --service url-manager --force-new-deployment
          aws ecs update-service --cluster scraping-cluster --service api-gateway --force-new-deployment
```

## 🛠️ Troubleshooting

### Common Issues

1. **Database Connection Failed**
   ```bash
   # Check PostgreSQL
   pg_isready -h localhost -p 5432
   
   # Check credentials
   psql -h localhost -U scraper -d scraper -c "SELECT 1;"
   ```

2. **Kafka Connection Issues**
   ```bash
   # Check Kafka health
   docker-compose -f docker-compose.local.yml logs kafka
   
   # Verify topics
   docker exec scraping_kafka kafka-topics --bootstrap-server localhost:9092 --list
   ```

3. **Service Won't Start**
   ```bash
   # Check logs
   docker-compose -f docker-compose.local.yml logs -f url-manager
   
   # Check configuration
   docker exec scraping_url_manager env | grep SCRAPING
   ```

### Debug Mode

Enable debug logging:

```bash
# Set debug log level
export SCRAPING_LOG_LEVEL=debug

# Restart services
docker-compose -f docker-compose.local.yml restart url-manager api-gateway
```

### Performance Tuning

```bash
# Increase Kafka partitions
docker exec scraping_kafka kafka-topics --bootstrap-server localhost:9092 \
  --alter --topic scraping-tasks --partitions 10

# Tune database connections
export SCRAPING_DATABASE_MAX_OPEN_CONNS=50
export SCRAPING_DATABASE_MAX_IDLE_CONNS=10
``` 