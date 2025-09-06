# Postman Collection for Ecommerce API

This directory contains comprehensive Postman collection and environment files for testing the Ecommerce API.

## Files Included

- **`Ecommerce-API.postman_collection.json`** - Complete API collection with all endpoints
- **`Ecommerce-API.postman_environment.json`** - Environment variables for the collection

## Quick Setup

### 1. Import Collection and Environment

1. Open Postman
2. Click **Import** button
3. Select both JSON files from this directory
4. The collection will appear in your workspace

### 2. Configure Environment

1. Select **"Ecommerce API Environment"** from the environment dropdown
2. Verify the `BASE_URL` is set to `http://localhost:8080` (or your server URL)
3. Other variables will be populated automatically during testing

### 3. Start Testing

1. **Start your API server** first:
   ```bash
   go run main.go
   ```

2. **Test the flow**:
   - Start with **Authentication > Sign Up** to create a user
   - Use **Authentication > Login** to get JWT token (automatically saved to `TOKEN`)
   - Create an admin user and login to get `ADMIN_TOKEN` for admin endpoints
   - Test other endpoints in logical order

## Collection Structure

### 🔐 Authentication
- Sign Up (creates new user account)
- Login (automatically saves JWT token)
- Logout

### 🏥 Health Checks  
- Basic Health Check
- Detailed Health Check
- Readiness Check
- Liveness Check
- Metrics

### 📂 Categories
- List Categories
- Add Category (Admin only)
- Delete Category (Admin only)

### 🛍️ Products
- List Products (with filtering)
- Get Single Product
- Search Products
- Add Product (Admin only)
- Edit Product (Admin only) 
- Delete Product (Admin only)

### 🛒 Shopping Cart
- View Cart
- Add to Cart
- Remove from Cart

### 📦 Orders
- Place Order
- List Orders
- Get Single Order
- Cancel Order

### 📍 Address Management
- Add Address
- Edit Address
- Delete Address

### ❤️ Wishlist
- View Wishlist
- Add to Wishlist
- Remove from Wishlist

### 💳 Payment Processing
- Process Payment
- Get Payment Status

### 🛒 Checkout
- Process Checkout (complete order flow)

### ⭐ Reviews
- Add Review
- List Reviews

### 📊 Admin Reports
- Sales Report
- Inventory Report

## Environment Variables

The collection uses these environment variables:

| Variable | Description | Auto-populated |
|----------|-------------|----------------|
| `BASE_URL` | API server URL | Manual |
| `TOKEN` | User JWT token | ✅ After login |
| `ADMIN_TOKEN` | Admin JWT token | Manual |
| `USER_ID` | Current user ID | Manual |
| `PRODUCT_ID` | Test product ID | Manual |
| `CATEGORY_ID` | Test category ID | Manual |
| `ORDER_ID` | Test order ID | Manual |

## Testing Workflow

### For Regular Users:
1. **Authentication** → Sign Up → Login
2. **Browse Products** → List Products → Get Product Details
3. **Shopping** → Add to Cart → View Cart
4. **Address** → Add Address
5. **Order** → Place Order → Check Order Status
6. **Reviews** → Add Product Review

### For Admin Users:
1. **Create Admin Account** → Sign Up with `"role": "admin"`
2. **Admin Login** → Save token as `ADMIN_TOKEN`
3. **Manage Catalog** → Add/Edit/Delete Products
4. **Manage Categories** → Add/Delete Categories
5. **View Reports** → Sales/Inventory Reports

## Advanced Features

### Automatic Token Management
- The Login request automatically extracts and saves the JWT token
- All authenticated requests use the saved token
- No manual token copying required!

### Pre-configured Test Data
- All requests include realistic test data
- Easy to modify for your specific testing needs
- Covers various scenarios and edge cases

### Environment Flexibility
- Easy to switch between development, staging, and production
- Just change the `BASE_URL` in environment settings

## Tips for Testing

1. **Start Fresh**: Use Health Check to verify API is running
2. **Follow Order**: Authentication first, then other endpoints
3. **Check Responses**: Verify success/error responses match documentation
4. **Test Edge Cases**: Try invalid data to test error handling
5. **Admin vs User**: Remember to use appropriate tokens for different roles

## Troubleshooting

### Common Issues:

**401 Unauthorized Error**
- Ensure you're logged in and token is saved
- Check if token has expired
- Verify correct token for user/admin endpoints

**404 Not Found**
- Verify API server is running on correct port
- Check endpoint URLs match your routes
- Ensure database has required test data

**Connection Refused**
- API server not running
- Wrong `BASE_URL` in environment
- Port conflicts

### Getting Help:
1. Check API server logs for detailed error messages
2. Verify your `.env` configuration
3. Test with curl commands to isolate Postman issues
4. Review the main README.md for API documentation

## Contributing

When adding new endpoints:
1. Add the request to appropriate collection folder
2. Include realistic test data in request body
3. Add any new environment variables needed
4. Update this README with the new endpoint

Happy Testing! 🚀
