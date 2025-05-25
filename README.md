# Go E-commerce Server

A RESTful API server built with Go, Gin, GORM, and Swagger documentation.

## Setup

1. Create a `.env` file in the root directory with the following content:
```env
DB_HOST=your-supabase-host
DB_USER=your-supabase-user
DB_PASSWORD=your-supabase-password
DB_NAME=your-supabase-dbname
DB_PORT=5432
DB_SSLMODE=require
```


2. Install dependencies:
```bash
go mod tidy
```

3. Run the server:
```bash
go run main.go
```

## API Documentation

Once the server is running, you can access the Swagger documentation at:
http://localhost:8080/swagger/index.html

## Available Endpoints

- GET /health - Health check endpoint
- GET /api/v1/products - List all products
- GET /api/v1/products/:id - Get a specific product
- POST /api/v1/products - Create a new product
- PUT /api/v1/products/:id - Update a product
- DELETE /api/v1/products/:id - Delete a product 
