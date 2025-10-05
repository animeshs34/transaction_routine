# Setup & Usage

## Prerequisites
- **Go 1.21+**
- **Docker & Docker Compose** (for containerized run)

---

## Run with Docker Compose

```bash
docker compose up --build
```

This starts:
- `app` (the transaction API)
- `postgres` (database)

App will be available at **http://localhost:8080**

---

## Run Locally with Go

1. Start Postgres:
   ```bash
   docker run -d --name txn-db -p 5432:5432 \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=postgres \
     -e POSTGRES_DB=transaction_routine \
     postgres:16-alpine
   ```

2. Edit `config/config.yaml`:
   ```yaml
   database:
     type: postgres
     host: localhost
     port: 5432
     user: postgres
     password: postgres
     dbname: transaction_routine
     sslmode: disable
   ```

3. Run the API:
   ```bash
   go mod tidy
   go run ./cmd/api
   ```

---

## Endpoints

### Health Check
```bash
curl http://localhost:8080/healthz
# → {"status":"ok"}
```

### Create Account
```bash
curl -X POST http://localhost:8080/accounts \
  -H 'Content-Type: application/json' \
  -d '{"document_number":"12345678900"}'
```

### Get Account
```bash
curl http://localhost:8080/accounts/1
```

### Create Transaction
```bash
curl -X POST http://localhost:8080/transactions \
  -H 'Content-Type: application/json' \
  -d '{"account_id":1,"operation_type_id":4,"amount":123.45}'
```

---

## Configuration

Defaults come from `config/config.yaml`.  
Environment variables (`APP_*`) override values.

| Config | Env Var | Default |
|--------|---------|---------|
| Server Port | `APP_SERVER_PORT` | 8080 |
| DB Type | `APP_DATABASE_TYPE` | memory |
| DB Host | `APP_DATABASE_HOST` | localhost |
| DB Port | `APP_DATABASE_PORT` | 5432 |
| DB User | `APP_DATABASE_USER` | postgres |
| DB Password | `APP_DATABASE_PASSWORD` | postgres |
| DB Name | `APP_DATABASE_DBNAME` | transaction_routine |
