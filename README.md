# DevPrep - Development Preparation Platform

A robust Go application for user management with comprehensive testing suite.

## Features

- User registration and authentication
- Session management
- User profile management
- Role-based access control
- RESTful API endpoints
- Comprehensive test coverage

## Tech Stack

- **Backend**: Go 1.24+
- **Web Framework**: Fiber v2
- **Database**: PostgreSQL
- **Testing**: Testify, SQLMock, Dockertest
- **Authentication**: Session-based with cookies
- **Password Hashing**: bcrypt

## Project Structure

```
.
├── cmd/devprep/           # Application entrypoint
├── internal/
│   ├── config/            # Configuration management
│   ├── database/          # Database connection
│   ├── dto/               # Data Transfer Objects
│   ├── handlers/          # HTTP handlers
│   ├── middleware/        # HTTP middleware
│   ├── models/            # Data models
│   ├── repository/        # Data access layer
│   ├── service/           # Business logic layer
│   ├── routes/            # Route definitions
│   └── utils/             # Utility functions
├── test/
│   ├── unit/              # Unit tests
│   ├── integration/       # Integration tests
│   ├── e2e/               # End-to-end tests
│   ├── helpers/           # Test utilities
│   ├── fixtures/          # Test data
│   └── testdata/          # Test database setup
├── database/migrations/   # Database migrations
└── .github/workflows/     # CI/CD pipelines
```

## Getting Started

### Prerequisites

- Go 1.24 or higher
- PostgreSQL 15+
- Docker (for tests)
- Make

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd devprep
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Start PostgreSQL and create databases:
```bash
make db-setup
```

5. Run migrations:
```bash
make migrate-up
```

### Running the Application

```bash
# Development mode
make dev

# Build and run
make build
make run
```

The server will start on port 3000 (configurable via SERVER_PORT environment variable).

## Testing

This project includes comprehensive testing at multiple levels:

### Unit Tests
Test business logic in isolation using mocks:
```bash
make test-unit
```

### Integration Tests
Test API endpoints with real database:
```bash
make test-integration
```

### End-to-End Tests
Test complete user flows:
```bash
make test-e2e
```

### All Tests
Run the complete test suite:
```bash
make test-all
```

### Test Coverage
Generate coverage reports:
```bash
make test-coverage
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout

### User Management
- `GET /api/v1/users/profile` - Get user profile
- `PUT /api/v1/users/profile` - Update user profile
- `DELETE /api/v1/users/profile` - Delete user account
- `GET /api/v1/users/` - List all users (authenticated)

### Health Checks
- `GET /healthz` - Health check
- `GET /readyz` - Readiness check

## Testing Strategy

The testing approach follows industry best practices:

1. **Unit Tests** - Fast, isolated tests for business logic
   - Mock external dependencies
   - Test edge cases and error conditions
   - High code coverage target (>80%)

2. **Integration Tests** - Test component interactions
   - Real database connections
   - HTTP request/response testing
   - Session management testing

3. **End-to-End Tests** - Complete user journey testing
   - Full application stack
   - User registration to deletion flows
   - Multi-user scenarios

### Test Database Setup

Tests use a separate PostgreSQL database that can be configured via:
- Docker containers (automated via dockertest)
- Local PostgreSQL instance
- In-memory database for unit tests

## Development Commands

```bash
# Code quality
make fmt          # Format code
make vet          # Vet code
make lint         # Lint code

# Dependencies
make deps         # Update dependencies

# Database
make migrate-up   # Run migrations
make migrate-down # Rollback migrations
make migrate-create # Create new migration

# Docker
make docker-up    # Start services
make docker-down  # Stop services

# Build
make build        # Build binary
make clean        # Clean artifacts
```

## Environment Variables

- `DATABASE_URL` - PostgreSQL connection string
- `SERVER_PORT` - Server port (default: 3000)
- `SESSION_SECRET` - Session encryption key
- `ENVIRONMENT` - Application environment (development/production)

## Contributing

1. Follow Go coding standards
2. Write tests for new features
3. Ensure all tests pass
4. Update documentation
5. Submit pull request

## License

This project is licensed under the MIT License.