# Ecommerce-RESTful-API

# E-commerce Store Backend RESTful API

This project is a backend implementation of an e-commerce store using Golang and the `net/http` package. It provides a RESTful API for managing users, products, shopping carts, and authentication services. This API serves as the core backend for an online store and supports features like user registration, login, cart management, and more.

## Features

- User authentication (SignUp, Login, Token Refresh, Account Deletion).
- Product catalog with filtering and detailed product information.
- Cart management with CRUD operations for adding, updating, and removing items.
- API health check for readiness.
- Middleware-based authentication and authorization.

## Technologies Used

- **Go (Golang)** for the backend server.
- **net/http** package for building RESTful API.
- **JWT** for token-based authentication.
- [goose](https://github.com/pressly/goose) for managing database schema and migrations
- [sqlc](https://github.com/sqlc-dev/sqlc) for sql to golang code generation
- **PostgreSQL**

## Getting Started

### Prerequisites

Ensure you have the following installed:

- **Go** (1.22+)
- **Git**
- **A Database** (e.g., PostgreSQL, MySQL)

### Clone the Repository

```sh
$ git clone https://github.com/Viet-ph/Ecommerce-Restful-API
$ cd Furniture-Store-Server
```

**Environment Variables**:</br>
Ensure that the necessary environment variables are set:
```sh
# Database Configuration
DB_HOST            # Database host
DB_PORT                # Database port (e.g., 5432 for PostgreSQL)
DB_USER           # Database username
DB_PASSWORD   # Database password
DB_NAME      # Database name
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# JWT Configuration
JWT_SECRET    # Secret key used to sign JWT tokens

# Server port
PORT
```


