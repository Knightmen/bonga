# Go E-commerce Server

A RESTful API server built with Go, Gin, GORM, and Swagger documentation.

## Setup

1. Create a `.env` file in the root directory with the following content:
```env
DB_HOST=your-postgresql-host
DB_USER=your-postgresql-user
DB_PASSWORD=your-postgresql-password
DB_NAME=your-postgresql-dbname
DB_PORT=5432
DB_SSLMODE=require

# AWS Configuration
AWS_REGION=your-aws-region
AWS_ACCESS_KEY_ID=your-aws-access-key
AWS_SECRET_ACCESS_KEY=your-aws-secret-key
AWS_S3_BUCKET=your-s3-bucket-name

# Resume Parse API Configuration
PARSE_API_TOKEN=your-parse-api-token
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
- GET /api/v1/resume - List all resumes (requires Authorization header)
- POST /api/v1/resume - Parse a resume file by calling external service (requires Authorization header and fileName in body)
- GET /api/v1/resume/:id - Get a specific resume (requires Authorization header)
- PUT /api/v1/resume/:id - Update a resume (requires Authorization header)
- DELETE /api/v1/resume/:id - Delete a resume (requires Authorization header)
- GET /api/v1/resume/getSignedUrl - Get a presigned URL for uploading a resume to S3 (requires filename query parameter and Authorization header)

## Development

### Hot Reload
Use air for hot reloading during development:
```bash
cd /Users/magic-kiri/go/bin/air
```

### Swagger Documentation
Generate swagger documentation using:
```bash
~/go/bin/swag init -g main.go
```


 