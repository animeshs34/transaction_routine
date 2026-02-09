# Transaction Routine API

A  Go-based microservice for managing customer accounts and financial transactions.

## Core Features

- **Account Management**: Create and retrieve customer profiles via document numbers.
- **Transaction Processing**: Register financial operations (Purchases, Withdrawals, Credits) with automatic normalization based on business rules.
- **Clean Architecture**: Strict separation of concerns between Domain, Service, Repository, and API layers.
- **Production-Grade Observability**: Structured JSON logging using Uber-Zap.
- **Resilient Infrastructure**: PostgreSQL with healthchecks and automated initialization.
- **Comprehensive Testing**: 90%+ coverage across unit and integration tests (using Testcontainers).

## Tech Stack

- **Go 1.24**: Leveraging latest language features.
- **Gin**: High-performance HTTP web framework.
- **PostgreSQL**: Stable and reliable relational data storage.
- **Docker**: Containerized deployment for consistent environments.

## Getting Started

### 1. Quick Start with Docker
```bash
make docker-up
```
This will start the API at `http://localhost:8080` and a PostgreSQL database.

### 2. Manual Execution
Ensure you have Go 1.24 installed:
```bash
make run
```

## Testing

Run the full test suite (requires Docker for integration tests):
```bash
make test
```

## API Documentation

### Create Account
`POST /accounts`
```json
{ "document_number": "12345678900" }
```

### Get Account
`GET /accounts/:id`

### Create Transaction
`POST /transactions`
```json
{
  "account_id": 1,
  "operation_type_id": 1,
  "amount": 100.50
}
```
*Note: Operation types 1, 2, and 3 result in negative amounts; type 4 results in positive.*

## Project Structure
- `cmd/api`: Entry point and dependency wiring.
- `internal/api`: HTTP handlers and request/response models.
- `internal/service`: Core business logic (Service interfaces and private implementations).
- `internal/repository`: Persistence logic using the Repository pattern.
- `internal/domain`: Pure business entities and domain rules.
- `pkg/errors`: Centralized error management with idiomatic wrapping.
