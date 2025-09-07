# 🚀 Quick Render Deployment Guide

## Step-by-Step Deployment to Render (Free Plan)

### 1. Push Your Code to GitHub
```bash
git add .
git commit -m "feat: add Render deployment configuration"
git push origin main
```

### 2. Sign Up and Connect to Render
1. Go to [render.com](https://render.com) and sign up
2. Click **"New"** → **"Blueprint"**
3. Connect your GitHub account
4. Select your `Ecommerce-api` repository

### 3. Render Will Automatically Create:
- ✅ **Web Service** (your API) - Free plan
- ✅ **PostgreSQL Database** - Free plan  
- ✅ **Environment Variables** configured automatically

### 4. Set Required Secret (Important!)
In the Render dashboard for your web service:
1. Go to **Environment** tab
2. Find `JWT_SECRET` variable
3. Make sure it's set to: `1234567890poiuytrewqasdfghjklmnbvcxz`

### 5. Your API Will Be Live At:
```
https://your-service-name.onrender.com
```

### 6. Test Your Deployment:
```bash
# Health check
curl https://your-service-name.onrender.com/health

# Create a user
curl -X POST https://your-service-name.onrender.com/signup \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com", 
    "password": "TestPass123",
    "phone": "+1234567890"
  }'
```

## Important Notes

### Free Plan Limitations:
- ⚠️ **Service sleeps after 15 minutes of inactivity**
- ⚠️ **750 hours/month limit** (about 31 days)
- ⚠️ **Database limited to 1GB**
- ⚠️ **First request after sleep takes ~30 seconds**

### Your Environment Variables:
```bash
# These are already configured in render.yaml:
DATABASE_HOST=<auto-configured>
DATABASE_USER=<auto-configured>
DATABASE_PASSWORD=<auto-configured>
DATABASE_NAME=ecommerce
DATABASE_SSLMODE=require
JWT_SECRET=1234567890poiuytrewqasdfghjklmnbvcxz
PORT=8080
ENV=production
LOG_LEVEL=info
```

## Monitoring Your App

### Check Status:
- **Render Dashboard**: Monitor deployments, logs, and metrics
- **Health Endpoint**: `GET /health` for basic status
- **Detailed Health**: `GET /health/detailed` for system info

### View Logs:
1. Go to Render dashboard
2. Select your service
3. Click **"Logs"** tab
4. See real-time application logs

## Troubleshooting

### Common Issues:

1. **Build Fails**:
   - Check Go version compatibility in `render.yaml`
   - Verify all dependencies in `go.mod`

2. **Database Connection Issues**:
   - Render automatically configures database connection
   - Check environment variables in dashboard
   - Database must use SSL in production (`DATABASE_SSLMODE=require`)

3. **App Not Responding**:
   - Free tier sleeps after inactivity
   - First request after sleep takes time
   - Check logs for startup errors

### Getting Help:
- **Render Docs**: [render.com/docs](https://render.com/docs)
- **Community**: [community.render.com](https://community.render.com)
- **Support**: Through Render dashboard

## Success! 🎉

Your Ecommerce API should now be live and accessible from anywhere in the world!

Test all endpoints using the Postman collection in your repository, just change the base URL to your Render deployment URL.
