# Ecommerce API

A robust REST API built with Go and Gin framework for managing an ecommerce platform. Features include user authentication, product management, shopping cart, orders, payments, and more.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Environment Variables](#environment-variables)
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
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=ecommerce
JWT_SECRET=your_secret_key
```

4. Start the server:
```bash
go run main.go
```


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
    "password": "securepass123",
    "email": "test@example.com",
    "phone": "1234567890",
    "role": "customer"  // Optional: use "admin" for admin account
}
```

#### Login
```http
POST /login
```

Test body:
```json
{
    "username": "testuser",
    "password": "securepass123"
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
    "address_id": 1,
    "payment_mode": "Credit Card"
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
    "comment": "Excellent product!"
}
```

#### List Reviews for Product
```http
GET /reviews/:product_id
```

### Wishlist

#### View Wishlist
```http
GET /wishlist
Authorization: Bearer <token>
```

#### Add to Wishlist
```http
POST /wishlist
Authorization: Bearer <token>
```

Test body:
```json
{
    "product_id": 1
}
```

#### Remove from Wishlist
```http
DELETE /wishlist/:id
Authorization: Bearer <token>
```

### Address

#### Add Address
```http
POST /address
Authorization: Bearer <token>
```

Test body:
```json
{
    "address": "123 Test Street",
    "city": "Test City",
    "zip_code": "12345"
}
```

#### Edit Address
```http
PUT /address/:id
Authorization: Bearer <token>
```

Test body:
```json
{
    "address": "456 Updated Street",
    "city": "New City",
    "zip_code": "54321"
}
```

#### Delete Address
```http
DELETE /address/:id
Authorization: Bearer <token>
```

### Payments

#### Process Payment
```http
POST /payments
Authorization: Bearer <token>
```

Test body:
```json
{
    "order_id": 1,
    "payment_method": "Credit Card",
    "amount": 999.99
}
```

#### Get Payment Status
```http
GET /payments/:order_id
Authorization: Bearer <token>
```

#### Checkout
```http
POST /checkout
Authorization: Bearer <token>
```

### Admin Reports

#### Sales Report (Admin Only)
```http
GET /admin/reports/sales?start_date=2024-01-01&end_date=2025-01-31
Authorization: Bearer <token>
```
#### Inventory Report (Admin Only)
```http
GET /admin/reports/inventory?start_date=2024-01-01&end_date=2025-01-31
Authorization: Bearer <token>
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
