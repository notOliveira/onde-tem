# Onde tem? / Find it

System to find information about establishments, enriched with collaborative informations.

## Terminal commands

Build the images and start the containers:
```bash
docker compose up --build
# Or using Makefile
make up-build
# Stop
Ctrl + C
```

Format all Go files in the project:
```bash
gofmt -w .
# Or
make fmt
```

Access the running PostgreSQL container:
```bash
docker exec -it <container_name> psql -U <user> -d <database>
# Example: docker exec -it app-db psql -U postgres -d app
# Or using Makefile:
make db-shell user=<user> db=<database>
# Make sure the database container is running before executing this command.
```

Stop and Remove Containers
```bash
docker compose down -v
# Or
make remove
# Clean Entire Docker Environment (Project Only)
make remove-all
```
## Running migrations

Create a new migration
```bash
docker run --rm -v ./migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq {migration_name}
# Or
make migrate-create name={migration_name}
```

Run all pending migrations
```bash
docker compose run migrate up
# Or
make migrate-up
```

Rollback last migration
```bash
docker compose run migrate down 1
# Or
make migrate-down
```

Reset all migrations (rollback everything)
```bash
docker compose run migrate down -all
# Or
make migrate-reset
```

Check current migration version
```bash
docker compose run migrate version
# Or
make migrate-version
```

Force database to a specific version (fix dirty state)
```bash
docker compose run migrate force <version>
# Or
make migrate-force v=<version>
```

Create a new migration
```bash
docker run --rm -v ./migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq <name>
# Or
make migrate-create name=<migration_name>
```
