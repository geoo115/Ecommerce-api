# 🚀 Deployment Guide - Render

This guide covers deploying the Ecommerce API to Render cloud platform.

## Prerequisites

1. **Render Account**: Sign up at [render.com](https://render.com)
2. **GitHub Repository**: Your code should be pushed to GitHub
3. **Environment Variables**: Prepare your production configuration

## Automatic Deployment Setup

### Option 1: Using render.yaml (Recommended)

1. **Connect Repository**:
   - Go to Render dashboard
   - Click "New" → "Blueprint"
   - Connect your GitHub repository
   - Render will automatically detect `render.yaml`

2. **Configure Secrets**:
   ```bash
   # In Render dashboard, set these environment variables:
   JWT_SECRET=your_super_secure_jwt_secret_here_minimum_32_characters
   ```

3. **Deploy**:
   - Click "Apply" to create all services
   - Render will automatically deploy your app with database

### Option 2: Manual Setup

1. **Create PostgreSQL Database**:
   - Go to Render dashboard
   - Click "New" → "PostgreSQL"
   - Name: `ecommerce-db`
   - Plan: Starter (or higher for production)

2. **Create Web Service**:
   - Click "New" → "Web Service"
   - Connect your GitHub repository
   - Configure:
     - **Name**: `ecommerce-api`
     - **Runtime**: Go
     - **Build Command**: `go build -o bin/main main.go`
     - **Start Command**: `./bin/main`
     - **Plan**: Starter (or higher)

3. **Set Environment Variables**:
   ```bash
   ENV=production
   PORT=8080
   LOG_LEVEL=info
   LOG_FORMAT=json
   RATE_LIMIT_ENABLED=true
   RATE_LIMIT_REQUESTS=1000
   BCRYPT_COST=12
   CORS_ALLOWED_ORIGINS=*
   
   # Database (auto-filled from database service)
   DATABASE_HOST=[from database]
   DATABASE_PORT=[from database]
   DATABASE_USER=[from database]
   DATABASE_PASSWORD=[from database]
   DATABASE_NAME=[from database]
   DATABASE_SSLMODE=require
   
   # JWT (set as secret)
   JWT_SECRET=your_super_secure_jwt_secret_here_minimum_32_characters
   JWT_EXPIRY=24h
   ```

## GitHub Actions Integration

The repository includes GitHub Actions workflows that automatically deploy to Render on code changes.

### Required GitHub Secrets

Add these secrets to your GitHub repository settings:

```bash
# Get from Render Account Settings → API Keys
RENDER_API_KEY=your_render_api_key

# Get from your Render service settings
RENDER_SERVICE_ID_STAGING=srv-xxx (for develop branch)
RENDER_SERVICE_ID_PRODUCTION=srv-xxx (for main branch)
```

### Deployment Triggers

- **Staging**: Automatic deployment on push to `develop` branch
- **Production**: Automatic deployment on push to `main` branch
- **Manual**: Use "Run workflow" button in GitHub Actions

## Configuration

### Environment-Specific Settings

#### Production
```yaml
ENV=production
LOG_LEVEL=info
DATABASE_SSLMODE=require
RATE_LIMIT_ENABLED=true
BCRYPT_COST=12
```

#### Staging
```yaml
ENV=staging
LOG_LEVEL=debug
RATE_LIMIT_ENABLED=false
```

### Security Considerations

1. **JWT Secret**: Generate a secure random string (minimum 32 characters)
   ```bash
   # Generate secure JWT secret
   openssl rand -base64 32
   ```

2. **Database SSL**: Always use `DATABASE_SSLMODE=require` in production

3. **CORS**: Configure `CORS_ALLOWED_ORIGINS` with your frontend domains only

4. **Rate Limiting**: Enable in production to prevent abuse

### Health Checks

Render automatically uses the `/health` endpoint for health checks. The API provides:

- **Basic Health**: `GET /health`
- **Detailed Health**: `GET /health/detailed`
- **Readiness**: `GET /ready`
- **Liveness**: `GET /live`

## Monitoring

### Built-in Monitoring

The API includes comprehensive monitoring:

- **Health Checks**: Automatic service health monitoring
- **Metrics**: Performance and system metrics at `/metrics`
- **Logging**: Structured JSON logs
- **Error Tracking**: Centralized error handling

### External Monitoring (Optional)

Consider integrating with:

- **Grafana**: For metrics dashboards
- **Sentry**: For error tracking
- **New Relic**: For APM monitoring

## Scaling

### Vertical Scaling
- Upgrade your Render plan for more CPU/RAM
- Plans: Starter → Standard → Pro

### Horizontal Scaling
- Render automatically handles load balancing
- Consider upgrading to Pro plans for auto-scaling

### Database Scaling
- Monitor database performance in Render dashboard
- Upgrade PostgreSQL plan when needed
- Consider read replicas for high-traffic applications

## Troubleshooting

### Common Issues

1. **Build Failures**:
   ```bash
   # Check Go version compatibility
   go version
   
   # Verify dependencies
   go mod tidy
   go mod verify
   ```

2. **Database Connection Issues**:
   ```bash
   # Verify environment variables are set
   echo $DATABASE_HOST
   
   # Check SSL mode setting
   echo $DATABASE_SSLMODE
   ```

3. **Health Check Failures**:
   ```bash
   # Test health endpoint locally
   curl http://localhost:8080/health
   
   # Check service logs in Render dashboard
   ```

### Getting Help

1. **Render Logs**: Check service logs in Render dashboard
2. **GitHub Actions**: Review workflow logs for deployment issues
3. **Health Endpoints**: Use `/health/detailed` for diagnostic information

## Production Checklist

Before going live:

- [ ] JWT secret is secure and unique
- [ ] Database SSL is enabled
- [ ] CORS origins are properly configured
- [ ] Rate limiting is enabled
- [ ] Environment is set to "production"
- [ ] Health checks are passing
- [ ] Monitoring is configured
- [ ] Backup strategy is in place
- [ ] SSL certificate is configured (automatic with Render)

## Cost Optimization

### Free Tier Limitations
- Render free tier has limitations (sleeps after inactivity)
- Consider Starter plan for production ($7/month)

### Resource Monitoring
- Monitor CPU and memory usage
- Optimize queries and caching
- Use appropriate instance sizes

## Support

- **Render Support**: [render.com/support](https://render.com/support)
- **API Documentation**: See README.md
- **Issues**: GitHub repository issues
