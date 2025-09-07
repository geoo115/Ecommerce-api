# Railway Deployment Guide

This guide will help you deploy your E-commerce API to Railway.

## Prerequisites

1. Sign up at [railway.app](https://railway.app)
2. Install Railway CLI (optional): `npm install -g @railway/cli`

## Quick Deploy (Recommended)

### Option 1: Deploy Button (Easiest)
[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/go-api)

### Option 2: GitHub Integration

1. Go to [railway.app](https://railway.app)
2. Click "New Project"
3. Select "Deploy from GitHub repo"
4. Choose your `Ecommerce-api` repository
5. Railway will automatically detect it's a Go project

## Database Setup

Railway provides managed PostgreSQL:

1. In your Railway project dashboard:
   - Click "New" → "Database" → "Add PostgreSQL"
   - Railway will create a database and provide connection details

2. Environment variables will be automatically set:
   - `DATABASE_URL` (complete connection string)
   - `PGHOST`, `PGPORT`, `PGDATABASE`, `PGUSER`, `PGPASSWORD`

## Environment Variables

Set these in your Railway project settings:

### Required Variables
```bash
# Database (automatically set by Railway PostgreSQL)
DATABASE_URL=postgresql://username:password@host:port/database

# JWT Configuration
JWT_SECRET=your-super-secure-jwt-secret-key-here

# Server Configuration
PORT=8080
GIN_MODE=release
```

### Optional Variables
```bash
# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60

# Redis (if using Railway Redis)
REDIS_URL=redis://default:password@host:port

# Metrics
METRICS_ENABLED=true
```

## Project Structure

Your project includes:
- `railway.json` - Railway configuration
- `main.go` - Application entry point with health checks
- `/health` endpoint - For Railway health monitoring
- Automatic database connection with Railway PostgreSQL

## Deployment Process

1. **Automatic Deployment**: 
   - Push to your `main` branch
   - Railway automatically builds and deploys
   - No additional CI/CD setup needed

2. **Build Process**:
   - Railway uses Nixpacks to detect Go project
   - Runs: `go build -o bin/main main.go`
   - Starts with: `./bin/main`

3. **Health Monitoring**:
   - Railway monitors `/health` endpoint
   - Automatic restarts on failure
   - Healthcheck timeout: 100s

## Domain and URLs

- Railway provides a free domain: `https://your-app-name.railway.app`
- Custom domains available on paid plans
- HTTPS enabled by default

## Cold Start Prevention

Railway has better cold start handling than Render, but you can still use:

```bash
# Keep warm with Railway
*/10 * * * * curl https://your-app-name.railway.app/health
```

## Monitoring

Railway dashboard provides:
- Real-time logs
- Resource usage metrics  
- Deployment history
- Database metrics

## Pricing

- **Hobby Plan**: $5/month per service
- **Pro Plan**: Usage-based pricing
- More predictable than Render's pricing

## Migration from Render

1. Export your Render database (if needed)
2. Deploy to Railway using this guide
3. Import data to Railway PostgreSQL
4. Update DNS (if using custom domain)
5. Test your endpoints

## Support

- Railway Documentation: [docs.railway.app](https://docs.railway.app)
- Railway Discord: [railway.app/discord](https://railway.app/discord)
- Railway Support: support@railway.app

## Quick Test Commands

```bash
# Health check
curl https://your-app-name.railway.app/health

# Test signup
curl -X POST https://your-app-name.railway.app/signup \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"testpass123"}'

# Test login
curl -X POST https://your-app-name.railway.app/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"testpass123"}'
```
