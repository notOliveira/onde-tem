# Onde tem? / Find it

System to find information about establishments, enriched with collaborative informations.

--

## Terminal commands

Build the images and start the containers:
```bash
docker compose up --build
# Or using Makefile
make up-build
# Stop
Ctrl + C
```

Start containers (without rebuild):
```bash
docker compose up -d
# Or
make up
```

View logs:
```bash
docker compose logs -f
# Or
make logs
```

Rebuild images only:
```bash
docker compose build
# Or
make build
```

Restart application (full rebuild):
```bash
docker compose down && docker compose up -d --build
# Or
make restart
```

Stop application:
```bash
docker compose down
# Or
make down
```

Stop and remove containers (including volumes):
```bash
docker compose down -v
# Or
make remove
```

Clean entire Docker environment (project only):
```bash
make remove-all
```

---

## Development commands

Run application locally (without Docker):
```bash
go run ./cmd/api
# Or
make run
```

Run tests:
```bash
go test ./...
# Or
make test
```

Format Go files:
```bash
gofmt -w .
# Or
make fmt
```

Run go vet:
```bash
go vet ./...
# Or
make vet
```

Run linter:
```bash
golangci-lint run
# Or
make lint
```

Generate SQLC files:
```bash
sqlc generate
# Or
make sqlc
```

Tidy go modules:
```bash
go mod tidy
# Or
make tidy
```

---

## Database commands

Access the running PostgreSQL container:
```bash
docker compose exec db psql -U <user> -d <database>
# Or
make db-shell user=<user> db=<database>
```

Make sure the database container is running before executing this command.

---

## Running migrations

Create a new migration:
```bash
docker run --rm -v ./migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq <migration_name>
# Or
make migrate-create name=<migration_name>
```

Run all pending migrations:
```bash
docker compose run migrate up
# Or
make migrate-up
```

Rollback last migration:
```bash
docker compose run migrate down 1
# Or
make migrate-down
```

Reset all migrations (rollback everything):
```bash
docker compose run migrate down -all
# Or
make migrate-reset
```

Check current migration version:
```bash
docker compose run migrate version
# Or
make migrate-version
```

Force database to a specific version (fix dirty state):
```bash
docker compose run migrate force <version>
# Or
make migrate-force v=<version>
```

---

## Production build

Build Linux binary:
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api ./cmd/api
# Or
make build-bin
```
