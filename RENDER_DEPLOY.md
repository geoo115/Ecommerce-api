# 🚀 Manual Render Deployment Guide (UPDATED)

## ⚠️ Important: Use Manual Setup (Blueprint Had Issues)

The Blueprint deployment failed, so we'll do a manual setup which is more reliable.

## Step-by-Step Manual Deployment

### 1. Create PostgreSQL Database First

1. Go to [render.com](https://render.com) and sign in
2. Click **"New"** → **"PostgreSQL"**
3. Configure:
   - **Name**: `ecommerce-db`
   - **Database**: `ecommerce`
   - **User**: `ecommerce_user`  
   - **Plan**: **Free**
4. Click **"Create Database"**
5. **IMPORTANT**: Copy the database connection details:
   ```
   Host: dpg-xxxxx-a.oregon-postgres.render.com
   Port: 5432
   Database: ecommerce
   Username: ecommerce_user
   Password: [generated-password]
   ```

### 2. Create Web Service

1. Click **"New"** → **"Web Service"**
2. **Connect your GitHub repository**: `geoo115/Ecommerce-api`
3. Configure:
   - **Name**: `ecommerce-api`
   - **Runtime**: **Go**
   - **Plan**: **Free**
   - **Build Command**: `go build -o bin/main main.go`
   - **Start Command**: `./bin/main`

### 3. Set Environment Variables

In the **Environment** section of your web service, add these variables:

```bash
ENV=production
PORT=8080
LOG_LEVEL=info
LOG_FORMAT=json
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=1000
BCRYPT_COST=12
CORS_ALLOWED_ORIGINS=*

# Database (use values from step 1)
DATABASE_HOST=dpg-xxxxx-a.oregon-postgres.render.com
DATABASE_PORT=5432
DATABASE_USER=ecommerce_user
DATABASE_PASSWORD=[your-generated-password]
DATABASE_NAME=ecommerce
DATABASE_SSLMODE=require

# JWT Secret
JWT_SECRET=1234567890poiuytrewqasdfghjklmnbvcxz
JWT_EXPIRY=24h
```

### 4. Deploy

1. Click **"Create Web Service"**
2. Render will automatically start building and deploying
3. Wait 5-10 minutes for deployment to complete

### 5. Get Your Live URL

Your API will be available at:
```
https://ecommerce-api-[random].onrender.com
```

## Test Your Deployment

```bash
# Health check (replace with your actual URL)
curl https://your-app-url.onrender.com/health

# Expected response:
{
  "success": true,
  "message": "Health check passed", 
  "data": {
    "status": "healthy",
    "timestamp": "..."
  }
}
```

## Create Your First User

```bash
curl -X POST https://your-app-url.onrender.com/signup \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "TestPass123",
    "phone": "+1234567890"
  }'
```

## Update Postman Collection

1. Open Postman
2. Import your collection from `postman/` folder  
3. Update the base URL to your Render URL
4. Start testing all endpoints!

## Troubleshooting

### If Build Fails:
1. Check the build logs in Render dashboard
2. Verify Go version compatibility
3. Make sure all dependencies are in `go.mod`

### If Database Connection Fails:
1. Double-check all database environment variables
2. Ensure `DATABASE_SSLMODE=require` 
3. Verify database is running in Render dashboard

### If App Won't Start:
1. Check the service logs
2. Verify `PORT=8080` is set
3. Make sure JWT_SECRET is exactly as shown above

## Success! 🎉

Your Ecommerce API should now be live and working!

**Remember**: Free plan sleeps after 15 minutes of inactivity. First request after sleep takes ~30 seconds.
