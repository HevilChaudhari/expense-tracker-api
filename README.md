# 💰 Expense Tracker REST API

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![Database](https://img.shields.io/badge/PostgreSQL-14+-4169E1?style=flat-square&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Auth](https://img.shields.io/badge/Auth-JWT-black?style=flat-square&logo=jsonwebtokens)](https://jwt.io)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)

A robust, modular, and secure RESTful Expense Tracker API built in **Go (Golang)** using standard library HTTP routing, **PostgreSQL** with high-performance connection pooling via `pgx/v5`, and **JWT (JSON Web Token)** authentication.

Designed following **Clean / Layered Architecture** principles (Handlers → Services → Repositories) ensuring clear separation of concerns, scalability, and maintainability.

---

## 📑 Table of Contents

- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [Project Architecture](#-project-architecture)
- [Database Schema](#-database-schema)
- [Getting Started](#-getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation & Setup](#installation--setup)
  - [Environment Variables](#environment-variables)
  - [Running the Server](#running-the-server)
- [API Documentation](#-api-documentation)
  - [Authentication](#authentication)
    - [Register User](#1-register-a-new-user)
    - [Login User](#2-user-login)
  - [Expenses Management](#expenses-management)
    - [Create Expense](#3-create-an-expense)
    - [Get All Expenses / Filter](#4-get-all-expenses-with-filters)
    - [Get Expense by ID](#5-get-expense-by-id)
    - [Update Expense](#6-update-an-expense)
    - [Delete Expense](#7-delete-an-expense)
    - [Get Expense Summary](#8-get-financial-summary)
- [How to Push to GitHub](#-how-to-push-to-github)
- [License](#-license)

---

## ✨ Features

- **User Authentication & Authorization**: Secure signup and login with hashed passwords (`bcrypt`) and stateless JWT token authentication.
- **Data Isolation**: Multi-tenant data segregation ensuring users can only view, edit, or delete their own expenses.
- **Full CRUD for Expenses**: Create, read, update, and delete expense records.
- **Advanced Query Filtering**: Filter expenses flexibly using query parameters:
  - By category (`?category=Food`)
  - By minimum amount (`?minAmount=50`)
  - By maximum amount (`?maxAmount=500`)
- **Financial Analytics & Summary**: Real-time aggregation endpoint calculating:
  - Total number of expenses
  - Total expenditure
  - Average expense amount
  - Highest & lowest expense recorded
- **Clean Layered Architecture**: Decoupled codebase divided into Handlers, Business Services, and Repositories with PostgreSQL connection pooling.

---

## 🛠 Tech Stack

- **Language**: [Go (Golang)](https://go.dev/)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **Database Driver / Connection Pool**: [`jackc/pgx/v5`](https://github.com/jackc/pgx)
- **Authentication**: [`golang-jwt/jwt/v5`](https://github.com/golang-jwt/jwt)
- **Password Hashing**: [`golang.org/x/crypto/bcrypt`](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- **Configuration**: [`joho/godotenv`](https://github.com/joho/godotenv)

---

## 🏛 Project Architecture

```
expense-tracker-api/
├── cmd/
│   └── server/
│       └── main.go               # Application entry point, dependency wiring, routes
├── internal/
│   ├── auth/
│   │   └── jwt.go                # JWT creation and validation
│   ├── database/
│   │   └── postgres.go           # PostgreSQL connection pool initializer (pgxpool)
│   ├── handlers/
│   │   ├── expense_handler.go    # HTTP handlers for expenses & summary
│   │   └── user_handler.go       # HTTP handlers for user registration & login
│   ├── middleware/
│   │   └── auth_middleware.go    # JWT validation middleware & context injection
│   ├── models/
│   │   ├── expense.go            # Expense domain entity
│   │   ├── expense_summary.go    # Analytics summary model
│   │   └── user.go               # User domain entity
│   ├── repositories/
│   │   ├── expense_repository.go # Database queries for expenses table
│   │   └── user_repository.go    # Database queries for users table
│   └── services/
│       ├── expense_service.go    # Business logic & calculations for expenses
│       └── user_service.go       # Business logic for auth & bcrypt verification
├── .env.example                  # Environment configuration template
├── .gitignore                    # Git ignore file
├── go.mod                        # Go module dependencies
├── go.sum                        # Checksums for dependencies
└── README.md                     # Project documentation
```

---

## 🗄 Database Schema

Ensure PostgreSQL is running and execute the following SQL statements to create the database schema:

```sql
-- 1. Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Create expenses table
CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    category VARCHAR(100) NOT NULL,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- 3. (Optional) Create index for user_id on expenses for faster lookups
CREATE INDEX IF NOT EXISTS idx_expenses_user_id ON expenses(user_id);
```

---

## 🚀 Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) installed (version 1.22 or higher)
- [PostgreSQL](https://www.postgresql.org/download/) installed and running
- Git

### Installation & Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/HevilChaudhari/expense-tracker-api.git
   cd expense-tracker-api
   ```

2. **Download Go modules:**
   ```bash
   go mod download
   ```

3. **Configure Environment Variables:**
   Copy the example `.env.example` file to `.env`:
   ```bash
   # Linux / macOS
   cp .env.example .env

   # Windows (PowerShell)
   Copy-Item .env.example .env
   ```

4. **Update `.env` values:**
   ```env
   DATABASE_URL=postgres://your_postgres_user:your_password@localhost:5432/expense_tracker?sslmode=disable
   JWT_SECRET=your_super_secret_jwt_key_here
   ```

### Running the Server

Start the application:

```bash
go run cmd/server/main.go
```

The server will start at `http://localhost:8080`.

---

## 📡 API Documentation

### Base URL
`http://localhost:8080`

### Endpoints Overview

| Method | Endpoint | Auth Required | Description |
| :--- | :--- | :---: | :--- |
| `POST` | `/register` | No | Register a new user |
| `POST` | `/login` | No | Login and obtain JWT token |
| `GET` | `/expenses` | **Yes** | Get all expenses (supports filtering) |
| `POST` | `/expenses` | **Yes** | Add a new expense |
| `GET` | `/expenses/{id}` | **Yes** | Get expense details by ID |
| `PUT` | `/expenses/{id}` | **Yes** | Update an existing expense |
| `DELETE` | `/expenses/{id}` | **Yes** | Delete an expense |
| `GET` | `/expenses/summary`| **Yes** | Get analytics summary (totals, averages, min/max) |

---

### Authentication

#### 1. Register a New User
- **URL**: `/register`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "name": "Jane Doe",
    "email": "jane@example.com",
    "password": "securepassword123"
  }
  ```
- **Response**: `201 Created`
  ```json
  {
    "id": 1,
    "name": "Jane Doe",
    "email": "jane@example.com",
    "created_at": "2026-09-30T10:15:30Z"
  }
  ```

#### 2. User Login
- **URL**: `/login`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "email": "jane@example.com",
    "password": "securepassword123"
  }
  ```
- **Response**: `200 OK`
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "name": "Jane Doe",
      "email": "jane@example.com",
      "created_at": "2026-09-30T10:15:30Z"
    }
  }
  ```

> 💡 **Important:** Include the token in subsequent requests via the `Authorization` header:  
> `Authorization: Bearer <your_jwt_token>`

---

### Expenses Management

#### 3. Create an Expense
- **URL**: `/expenses`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer <TOKEN>`
- **Request Body**:
  ```json
  {
    "title": "Groceries at Supermarket",
    "amount": 84.50,
    "category": "Food"
  }
  ```
- **Response**: `201 Created`

#### 4. Get All Expenses (with Filters)
- **URL**: `/expenses`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <TOKEN>`
- **Optional Query Parameters**:
  - `category`: Filter by exact category (e.g. `?category=Food`)
  - `minAmount`: Filter by minimum expense amount (e.g. `?minAmount=50`)
  - `maxAmount`: Filter by maximum expense amount (e.g. `?maxAmount=200`)
- **Example**: `/expenses?category=Food&minAmount=20&maxAmount=100`
- **Response**: `200 OK`
  ```json
  [
    {
      "id": 1,
      "title": "Groceries at Supermarket",
      "amount": 84.5,
      "category": "Food"
    }
  ]
  ```

#### 5. Get Expense by ID
- **URL**: `/expenses/{id}`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <TOKEN>`
- **Response**: `200 OK`
  ```json
  {
    "id": 1,
    "title": "Groceries at Supermarket",
    "amount": 84.5,
    "category": "Food"
  }
  ```

#### 6. Update an Expense
- **URL**: `/expenses/{id}`
- **Method**: `PUT`
- **Headers**: `Authorization: Bearer <TOKEN>`
- **Request Body**:
  ```json
  {
    "title": "Weekly Organic Groceries",
    "amount": 95.00,
    "category": "Food"
  }
  ```
- **Response**: `200 OK`
  ```json
  {
    "id": 1,
    "title": "Weekly Organic Groceries",
    "amount": 95,
    "category": "Food"
  }
  ```

#### 7. Delete an Expense
- **URL**: `/expenses/{id}`
- **Method**: `DELETE`
- **Headers**: `Authorization: Bearer <TOKEN>`
- **Response**: `204 No Content`

#### 8. Get Financial Summary
- **URL**: `/expenses/summary`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <TOKEN>`
- **Response**: `200 OK`
  ```json
  {
    "totalExpenses": 5,
    "totalAmount": 342.75,
    "averageAmount": 68.55,
    "highestAmount": 150.00,
    "lowestAmount": 12.50
  }
  ```

---

## 📤 How to Push to GitHub

To commit and push the latest changes (including `README.md` and `.env.example`) to your GitHub repository:

```bash
# 1. Check changed files
git status

# 2. Stage changes
git add README.md .env.example

# 3. Commit changes
git commit -m "docs: add comprehensive README and .env.example"

# 4. Push to GitHub
git push origin main
```

---

## 📄 License

This project is licensed under the [MIT License](LICENSE) - see the LICENSE file for details.
