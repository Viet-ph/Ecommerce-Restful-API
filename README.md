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

```bash
$ git clone https://github.com/Viet-ph/Furniture-Store-Server.git
$ cd Furniture-Store-Server
