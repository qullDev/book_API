# Book API

RESTful service for managing books and categories with JWT authentication. The system uses Access Tokens for API access and Refresh Tokens stored in Redis for session management.

## Tech Stack

- **Go + Gin**: Web framework
- **PostgreSQL + GORM**: Database and ORM
- **Redis**: Refresh token storage
- **JWT**: Authentication using HS256 algorithm

## Project Structure

```
book_API/
├── cmd/
│   └── api/
│       └── main.go         # Application entry point
├── internal/
│   ├── config/            # Configuration management
│   ├── cache/             # Redis connection
│   ├── db/               # Database connection
│   ├── domain/           # Domain models
│   │   ├── book/
│   │   ├── category/
│   │   └── user/
│   ├── http/             # HTTP layer
│   │   ├── handlers/     # Request handlers
│   │   ├── middleware/   # HTTP middleware
│   │   └── router/       # Route definitions
│   └── pkg/
│       └── auth/         # Authentication utilities
```

## Prerequisites

- Go 1.24 or higher
- PostgreSQL 13+
- Redis 6+

## Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and configure:

   ```bash
   APP_PORT=8080

   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=books_db
   DB_SSLMODE=disable

   REDIS_ADDR=localhost:6379
   REDIS_PASSWORD=
   REDIS_DB=1

   JWT_SECRET=your_secret_key
   ACCESS_TOKEN_TTL=15m
   REFRESH_TOKEN_TTL=168h

   ENV=dev
   ```

3. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

## Authentication Flow

### 1. Login

```http
POST /api/users/login
Content-Type: application/json

{
    "username": "admin",
    "password": "password"
}
```

Response:

```json
{
  "access_token": "eyJhbG...",
  "refresh_token": "eyJhbG...",
  "token_type": "Bearer",
  "expires_in": 900,
  "refresh_expires_in": 604800,
  "user_id": "...",
  "username": "admin"
}
```

### 2. Refresh Token

```http
POST /api/users/refresh
Content-Type: application/json

{
    "refresh_token": "eyJhbG..."
}
```

### 3. Logout

```http
POST /api/users/logout
Authorization: Bearer eyJhbG...
Content-Type: application/json

{
    "refresh_token": "eyJhbG..."  // Optional, if omitted revokes all tokens
}
```

## Protected Endpoints

All endpoints below require `Authorization: Bearer {access_token}` header.

### Categories

#### List Categories

```http
GET /api/categories
```

#### Create Category

```http
POST /api/categories
Content-Type: application/json

{
    "name": "Fiction"
}
```

#### Get Category

```http
GET /api/categories/:id
```

#### Update Category

```http
PUT /api/categories/:id
Content-Type: application/json

{
    "name": "Non-Fiction"
}
```

#### Delete Category

```http
DELETE /api/categories/:id
```

#### List Books in Category

```http
GET /api/categories/:id/books
```

### Books

#### List Books

```http
GET /api/books
```

#### Create Book

```http
POST /api/books
Content-Type: application/json

{
    "title": "The Go Programming Language",
    "category_id": "uuid",
    "description": "Book description",
    "image_url": "https://example.com/image.jpg",
    "release_year": 2020,
    "price": 59.99,
    "total_page": 150
}
```

#### Get Book

```http
GET /api/books/:id
```

#### Update Book

```http
PUT /api/books/:id
Content-Type: application/json

{
    "title": "Updated Title",
    "price": 49.99
}
```

#### Delete Book

```http
DELETE /api/books/:id
```

## Models

### Book

```go
type Book struct {
    ID          uuid.UUID
    Title       string
    CategoryID  uuid.UUID
    Description string
    ImageURL    string
    ReleaseYear int
    Price       float64
    TotalPage   int
    Thickness   string    // Auto-calculated: "tipis" (<= 100 pages) or "tebal" (> 100 pages)
    CreatedAt   time.Time
    ModifiedAt  time.Time
}
```

### Category

```go
type Category struct {
    ID         uuid.UUID
    Name       string
    CreatedAt  time.Time
    ModifiedAt time.Time
}
```

## Validation Rules

- Release year must be between 1980 and 2024
- Book thickness is automatically set based on total pages:
  - ≤ 100 pages: "tipis"
  - > 100 pages: "tebal"
- Category name must be 1-100 characters
- Book title must be 1-200 characters

## Error Responses

The API returns consistent error responses:

```json
{
  "message": "Error description",
  "error": "Detailed error message (optional)"
}
```

Common status codes:

- 400: Bad Request (invalid input)
- 401: Unauthorized (invalid/expired token)
- 404: Not Found
- 500: Internal Server Error

## Development

1. Run PostgreSQL and Redis (using Docker):

   ```bash
   docker run -d --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres
   docker run -d --name redis -p 6379:6379 redis
   ```

2. Create database:

   ```sql
   CREATE DATABASE books_db;
   ```

3. Run the application in development mode:
   ```bash
   ENV=dev go run cmd/api/main.go
   ```
