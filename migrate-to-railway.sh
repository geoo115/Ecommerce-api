#!/bin/bash

# Railway Migration Script
# This script helps migrate from Render to Railway

set -e

echo "🚂 Railway Migration Helper"
echo "=========================="

# Check if Railway CLI is installed
if ! command -v railway &> /dev/null; then
    echo "📦 Installing Railway CLI..."
    curl -fsSL https://railway.app/install.sh | sh
    export PATH="$HOME/.railway/bin:$PATH"
else
    echo "✅ Railway CLI already installed"
fi

echo ""
echo "📋 Migration Checklist:"
echo "1. Sign up at https://railway.app"
echo "2. Create a new project from GitHub repo"
echo "3. Add PostgreSQL database"
echo "4. Set environment variables (see below)"
echo "5. Deploy and test"

echo ""
echo "🔧 Required Environment Variables for Railway:"
echo "============================================"
echo "DATABASE_URL=postgresql://user:pass@host:port/db  # Auto-set by Railway PostgreSQL"
echo "JWT_SECRET=your-super-secure-jwt-secret-32-chars+"
echo "PORT=8080"
echo "GIN_MODE=release"
echo ""
echo "Optional:"
echo "RATE_LIMIT_ENABLED=true"
echo "METRICS_ENABLED=true"
echo "LOG_LEVEL=info"

echo ""
echo "🧪 Quick Test Commands (after deployment):"
echo "=========================================="
echo "# Replace YOUR_APP_NAME with your Railway app name"
echo "export RAILWAY_URL=\"https://YOUR_APP_NAME.railway.app\""
echo ""
echo "# Health check"
echo "curl \$RAILWAY_URL/health"
echo ""
echo "# Test signup"
echo "curl -X POST \$RAILWAY_URL/signup \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"username\":\"test\",\"email\":\"test@example.com\",\"password\":\"testpass123\"}'"

echo ""
echo "📚 Next Steps:"
echo "============="
echo "1. Read RAILWAY_DEPLOY.md for detailed instructions"
echo "2. Update GitHub secrets for CI/CD:"
echo "   - RAILWAY_TOKEN_PRODUCTION"
echo "   - RAILWAY_PROJECT_ID_PRODUCTION"  
echo "   - RAILWAY_PRODUCTION_URL"
echo "3. Test your deployment"
echo "4. Update DNS (if using custom domain)"

echo ""
echo "🎉 Happy deploying with Railway!"
