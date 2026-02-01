# Go Kasir API

A simple Cashier/POS (Point of Sale) API built with Go (Golang), adhering to Clean Architecture principles.

## Features

- **Products Management**: CRUD operations for products (database-backed).
- **Categories Management**: CRUD operations for categories (database-backed).
- **Swagger Documentation**: Interactive API documentation.
- **Layered Architecture**: Handlers -> Services -> Repositories -> Database.

## Tech Stack

- **Language**: Go (Golang)
- **Database**: PostgreSQL (Supabase compatible)
- **Driver**: `pgx` (PostgreSQL Driver)
- **Configuration**: Viper
- **Documentation**: Swagger (Swaggo)

## API Endpoints

### Products

| Method   | Endpoint             | Description          |
| :------- | :------------------- | :------------------- |
| `GET`    | `/api/products`      | Get all products     |
| `POST`   | `/api/products`      | Create a new product |
| `GET`    | `/api/products/{id}` | Get product by ID    |
| `PUT`    | `/api/products/{id}` | Update product       |
| `DELETE` | `/api/products/{id}` | Delete product       |

### Categories

| Method   | Endpoint               | Description           |
| :------- | :--------------------- | :-------------------- |
| `GET`    | `/api/categories`      | Get all categories    |
| `POST`   | `/api/categories`      | Create a new category |
| `GET`    | `/api/categories/{id}` | Get category by ID    |
| `PUT`    | `/api/categories/{id}` | Update category       |
| `DELETE` | `/api/categories/{id}` | Delete category       |

### Docs

- Swagger UI: `/swagger/index.html`
- Health Check: `/health`

## Setup & Run

1.  **Clone Repository**
2.  **Environment Variables**
    Create `.env` file:
    ```env
    PORT=8080
    DB_CONN='postgres://user:pass@host:5432/db?sslmode=disable'
    ```
3.  **Run Application**
    ```bash
    go mod tidy
    go run main.go
    ```
