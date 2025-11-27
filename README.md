# Go REST API

A RESTful API built with Go that provides user authentication and management functionality. This project demonstrates a clean architecture approach with proper separation of concerns, database migrations, and secure password handling.

## Features

- **User Registration**: Create new user accounts with validation
- **User Authentication**: Secure login system with password hashing
- **User Management**: Retrieve user information by ID or email
- **Database Migrations**: Automated database schema management
- **Input Validation**: Request payload validation using struct tags
- **Secure Password Storage**: bcrypt password hashing
- **RESTful Design**: Clean API endpoints following REST principles

## Tech Stack

- **Language**: Go 1.25.1
- **Web Framework**: Gorilla Mux for routing
- **Database**: MySQL
- **Password Hashing**: bcrypt
- **Validation**: go-playground/validator
- **Migrations**: golang-migrate
- **Configuration**: Environment variables with godotenv

## Project Structure

```
├── cmd/
│   ├── main.go                 # Application entry point
│   ├── api/
│   │   └── api.go             # API server setup and routing
│   └── migrate/
│       ├── main.go            # Database migration runner
│       └── migrations/        # SQL migration files
├── configs/
│   └── env.go                 # Environment configuration
├── db/
│   └── db.go                  # Database connection setup
├── services/
│   ├── auth/
│   │   └── password.go        # Password hashing utilities
│   └── user/
│       ├── routes.go          # User HTTP handlers
│       ├── routes_test.go     # User route tests
│       └── store.go           # User database operations
├── types/
│   └── user.go                # User types and interfaces
├── utils/
│   └── utils.go               # HTTP utilities (JSON parsing, responses)
└── Makefile                   # Build and development commands
```

## API Endpoints

### Authentication & User Management

- `POST /api/v1/register` - Register a new user
- `POST /api/v1/login` - User login (authentication)
- `GET /api/v1/users/{userID}` - Get user by ID

## Getting Started

### Prerequisites

- Go 1.25.1 or higher
- MySQL database
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/capamir/go-api.git
   cd go-api
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   Create a `.env` file in the root directory:
   ```env
   PUBLIC_HOST=http://localhost
   PORT=8080
   DB_USER=root
   DB_PASSWORD=your_password
   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_NAME=your_database_name
   ```

4. **Run database migrations**
   ```bash
   make migrate-up
   ```

5. **Start the server**
   ```bash
   make run
   # or for development with auto-reload
   make dev
   ```

The API will be available at `http://localhost:8080`

## Usage Examples

### Register a New User

```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "John",
    "lastName": "Doe",
    "email": "john.doe@example.com",
    "password": "securepassword123"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "securepassword123"
  }'
```

### Get User by ID

```bash
curl -X GET http://localhost:8080/api/v1/users/1
```

## Development Commands

The project includes a `Makefile` with useful development commands:

- `make build` - Build the application
- `make run` - Build and run the application
- `make dev` - Run in development mode
- `make test` - Run all tests
- `make migrate-up` - Run database migrations
- `make migrate-down` - Rollback database migrations
- `make migration <name>` - Create a new migration file

## Database Schema

### Users Table

```sql
CREATE TABLE users (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  firstName VARCHAR(255) NOT NULL,
  lastName VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);
```

## Architecture

The project follows a clean architecture pattern with clear separation of concerns:

- **`cmd/`**: Application entry points and executables
- **`configs/`**: Configuration management
- **`db/`**: Database connection and setup
- **`services/`**: Business logic and HTTP handlers
- **`types/`**: Data models and interfaces
- **`utils/`**: Shared utilities and helpers

### Key Design Patterns

- **Repository Pattern**: `UserStore` interface for data access abstraction
- **Handler Pattern**: HTTP request handlers separated from business logic
- **Dependency Injection**: Services receive their dependencies through constructors
- **Interface Segregation**: Small, focused interfaces for better testability

## Security Features

- **Password Hashing**: All passwords are hashed using bcrypt
- **SQL Injection Prevention**: Parameterized queries for all database operations
- **Input Validation**: Request payload validation with comprehensive rules
- **Unique Email Constraint**: Database-level email uniqueness enforcement

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is open source and available under the [MIT License](LICENSE).