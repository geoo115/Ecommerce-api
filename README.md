# Ecommerce API

A robust REST API built with Go and Gin framework for managing an ecommerce platform. Features include user authentication, product management, shopping cart, orders, payments, and more.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Environment Variables](#environment-variables)
- [Security Features](#security-features)
- [Monitoring & Health Checks](#monitoring--health-checks)
- [Rate Limiting](#rate-limiting)
- [Logging](#logging)
- [API Documentation](#api-documentation)
  - [Authentication](#authentication)
  - [Categories](#categories)
  - [Products](#products)
  - [Cart](#cart)
  - [Orders](#orders)
  - [Reviews](#reviews)
  - [Wishlist](#wishlist)
  - [Payments](#payments)
  - [Address](#address)
  - [Admin Reports](#admin-reports)
  - [Health Checks](#health-checks)

## Prerequisites

- Go 1.16 or higher
- PostgreSQL
- Postman for testing

## Installation

1. Clone the repository:
```bash
git clone https://github.com/geoo115/Ecommerce.git
cd Ecommerce
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
cp env.example .env
# Edit .env with your actual values
```

4. Start the server:
```bash
go run main.go
```

## Environment Variables

Copy `env.example` to `.env` and configure the following variables:

```env
# Database Configuration
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=your_password
DATABASE_NAME=ecommerce
DATABASE_SSLMODE=disable

# JWT Configuration (REQUIRED - must be set)
JWT_SECRET=your_super_secret_jwt_key_here_make_it_long_and_random

# Server Configuration
PORT=8080

# Logging Configuration
LOG_LEVEL=info
# Available levels: debug, info, warn, error, fatal
```

## Security Features

This API includes several security measures:

### Authentication & Authorization
- JWT-based authentication with secure token generation
- Role-based access control (admin/customer)
- Password hashing using bcrypt
- Secure token validation

### Input Validation & Sanitization
- Comprehensive input validation for all endpoints
- SQL injection prevention using parameterized queries
- Input sanitization to prevent XSS attacks
- Request size limits and validation

### Security Headers
- CORS middleware for cross-origin requests
- Security headers (X-Frame-Options, X-Content-Type-Options, etc.)
- Content Security Policy (CSP)
- XSS protection headers

### Error Handling
- Standardized error responses
- No sensitive information in error messages
- Proper HTTP status codes
- Secure error logging

### Database Security
- Parameterized queries to prevent SQL injection
- Proper database connection handling
- Input validation before database operations

## Monitoring & Health Checks

The API provides comprehensive health check endpoints for monitoring:

### Health Check Endpoints

#### Basic Health Check
```http
GET /health
```
Returns basic application status.

#### Detailed Health Check
```http
GET /health/detailed
```
Returns detailed health status including database connectivity and system information.

#### Readiness Check
```http
GET /ready
```
Checks if the application is ready to serve traffic (database connectivity + startup time).

#### Liveness Check
```http
GET /live
```
Simple check to verify the application is alive.

#### Metrics
```http
GET /metrics
```
Returns application metrics including system information and memory usage.

### Health Check Response Format
```json
{
  "success": true,
  "message": "Health check passed",
  "data": {
    "status": "healthy",
    "timestamp": "2024-01-01T12:00:00Z",
    "uptime": "1h30m45s",
    "version": "1.0.0",
    "services": {
      "api": {"status": "healthy"},
      "database": {"status": "healthy"},
      "system": {
        "go_version": "go1.22.0",
        "architecture": "amd64",
        "os": "linux",
        "num_cpu": 8,
        "num_goroutine": 15
      }
    }
  }
}
```

## Rate Limiting

The API implements rate limiting to prevent abuse:

### Rate Limit Tiers

- **General API**: 100 requests per minute
- **Authentication endpoints**: 10 requests per minute (signup/login)
- **Admin endpoints**: 50 requests per minute

### Rate Limit Headers

When rate limited, the API returns:
- `X-RateLimit-Limit`: Maximum requests allowed
- `X-RateLimit-Remaining`: Remaining requests
- `X-RateLimit-Reset`: Time when limits reset

### Rate Limit Response
```json
{
  "success": false,
  "error": "Rate limit exceeded. Please try again later.",
  "code": 429
}
```

## Logging

The API includes structured logging with configurable levels:

### Log Levels
- `debug`: Detailed debug information
- `info`: General information (default)
- `warn`: Warning messages
- `error`: Error messages
- `fatal`: Fatal errors (exits application)

### Log Format
```
[2024-01-01 12:00:00] INFO: HTTP Request: GET /products from 192.168.1.1 - Status: 200 - Duration: 15ms
```

### Logged Events
- HTTP requests with timing and status
- Database operations
- Security events
- Application errors
- System metrics

## API Documentation and Testing Guide

This section provides detailed instructions for testing all API endpoints using Postman or similar tools.

### Authentication

#### Sign Up
```http
POST /signup
```

Test body:
```json
{
    "username": "testuser",
    "password": "SecurePass123",
    "email": "test@example.com",
    "phone": "1234567890",
    "role": "customer"  // Optional: use "admin" for admin account
}
```

**Validation Rules:**
- Username: 3-30 characters, alphanumeric and underscores only
- Password: Minimum 8 characters, must contain uppercase, lowercase, and numeric
- Email: Valid email format
- Phone: 10-15 digits

#### Login
```http
POST /login
```

Test body:
```json
{
    "username": "testuser",
    "password": "SecurePass123"
}
```

#### Logout
```http
POST /logout
Authorization: Bearer <token>
```

### Categories

#### List Categories
```http
GET /categories
```

#### Add Category (Admin Only)
```http
POST /categories
Authorization: Bearer <token>
```

Test body:
```json
{
    "name": "Electronics"
}
```

#### Delete Category (Admin Only)
```http
DELETE /categories/:id
Authorization: Bearer <token>
```

### Products

#### List Products
```http
GET /products
```

Query parameters:
- `category_id=1`
- `page=1`
- `limit=10`

#### Get Single Product
```http
GET /product/:id
```

#### Add Product (Admin Only)
```http
POST /product
Authorization: Bearer <token>
```

Test body:
```json
{
    "name": "Test Product",
    "price": 999.99,
    "category_id": 1,
    "description": "Test description",
    "stock": 50
}
```

#### Edit Product (Admin Only)
```http
PUT /product/:id
Authorization: Bearer <token>
```

Test body:
```json
{
    "name": "Updated Product",
    "price": 899.99,
    "description": "Updated description",
    "stock": 45
}
```

#### Delete Product (Admin Only)
```http
DELETE /product/:id
Authorization: Bearer <token>
```

#### Search Products
```http
GET /products/search?query=laptop
```

### Cart

#### View Cart
```http
GET /cart
Authorization: Bearer <token>
```

#### Add to Cart
```http
POST /cart
Authorization: Bearer <token>
```

Test body:
```json
{
    "product_id": 1,
    "quantity": 2
}
```

#### Remove from Cart
```http
DELETE /cart/:id
Authorization: Bearer <token>
```

### Orders

#### Place Order
```http
POST /orders
Authorization: Bearer <token>
```

Test body:
```json
{
    "items": [
        {
            "product_id": 1,
            "quantity": 2
        }
    ]
}
```

#### List Orders
```http
GET /orders
Authorization: Bearer <token>
```

#### Get Single Order
```http
GET /orders/:id
Authorization: Bearer <token>
```

#### Cancel Order
```http
PUT /orders/:id/cancel
Authorization: Bearer <token>
```

### Reviews

#### Add Review
```http
POST /reviews
Authorization: Bearer <token>
```

Test body:
```json
{
    "product_id": 1,
    "rating": 5,
    "comment": "Great product!"
}
```

#### List Reviews
```http
GET /reviews/:product_id
```

### Health Checks

#### Basic Health Check
```http
GET /health
```

#### Detailed Health Check
```http
GET /health/detailed
```

#### Readiness Check
```http
GET /ready
```

#### Liveness Check
```http
GET /live
```

#### Metrics
```http
GET /metrics
```

## Testing Steps

1. Start by creating a new user account using the signup endpoint
2. Login to get the JWT token
3. For admin operations, create an admin account and use its token
4. Add the token to your request headers for authenticated endpoints
5. Test each endpoint with both valid and invalid data to ensure proper error handling
6. For testing order flow:
   - Add products to cart
   - Create address
   - Place order
   - Process payment
   - Check order status

## Testing with Postman

1. Import the Postman collection from the `postman` directory
2. Set up environment variables in Postman:
   - `BASE_URL`: `http://localhost:8080`
   - `TOKEN`: After login, set this to the received JWT token

## Error Handling

The API returns standard HTTP status codes:

- 200: Successful operation
- 201: Resource created
- 400: Bad request (invalid input)
- 401: Unauthorized (invalid/missing token)
- 403: Forbidden (insufficient permissions)
- 404: Resource not found
- 500: Internal server error

Error Response Format:
```json
{
    "error": "Error message here"
}
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request


## Contact

Geoo115 - [GitHub Profile](https://github.com/geoo115)

Project Link: [https://github.com/geoo115/Ecommerce](https://github.com/geoo115/Ecommerce)
