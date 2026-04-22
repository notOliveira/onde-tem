include .env
export

APP_NAME=onde-tem
DB_URL=postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}

# =========================
# SANITY CHECKS
# =========================

sanity:
	@echo ================================
	@echo   🔍 SANITY CHECK STARTING
	@echo ================================

	@echo.
	@echo 🐳 Checking Docker...
	@docker ps > NUL 2>&1 || (echo ❌ Docker is NOT running & exit 1)
	@echo ✅ Docker is running

	@echo.
	@echo 🐘 Checking Postgres container...
	@docker compose ps db | findstr "Up" > NUL || (echo ❌ Postgres container is NOT running & exit 1)
	@echo ✅ Postgres container is up

	@echo.
	@echo 📡 Checking Postgres connection...
	@docker compose exec db pg_isready -U ${DB_USER} || (echo ❌ Postgres not accepting connections & exit 1)

	@echo.
	@echo 🧱 Checking migration version...
	@docker compose run --rm migrate version || (echo ❌ Migration check failed & exit 1)

	@echo.
	@echo ⚙️ Running sqlc generate...
	@sqlc generate || (echo ❌ sqlc failed & exit 1)
	@echo ✅ sqlc OK

	@echo.
	@echo 🧪 Running tests...
	@go test ./... || (echo ❌ Tests failing & exit 1)

	@echo.
	@echo ================================
	@echo   ✅ ENVIRONMENT HEALTHY
	@echo ================================


# =========================
# Docker
# =========================

reset-all:
	@echo =========================================
	@echo   ☢️  NUCLEAR RESET STARTING
	@echo =========================================

	@echo.
	@echo 🧹 1/5 Destroying old environment...
	@docker compose down -v --rmi all --remove-orphans || (echo ❌ Teardown failed & exit 1)
	@echo ✅ Environment cleaned

	@echo.
	@echo ⚙️  2/5 Generating database code (sqlc)...
	@sqlc generate || (echo ❌ sqlc generation failed & exit 1)
	@echo ✅ Code generated

	@echo.
	@echo 🏗️  3/5 Building fresh images...
	@docker compose build || (echo ❌ Build failed & exit 1)
	@echo ✅ Images built

	@echo.
	@echo 🐘 4/5 Bootstrapping database and running migrations...
	@docker compose run --rm migrate up || (echo ❌ Migrations failed & exit 1)
	@echo ✅ Database ready and migrated

	@echo.
	@echo 🚀 5/5 Starting all services...
	@docker compose up -d || (echo ❌ Failed to start services & exit 1)
	@echo ✅ All services running

	@echo.
	@echo =========================================
	@echo   ✅ ENVIRONMENT RESET COMPLETE
	@echo =========================================
	

up:
	docker compose up -d

up-build:
	docker compose up --build

down:
	docker compose down

logs:
	docker compose logs -f

build:
	docker compose build

restart:
	docker compose down && docker compose up -d --build

remove:
	docker compose down -v

remove-all:
	docker compose down -v --rmi all --remove-orphans


# =========================
# Migrations
# =========================

migrate-up:
	docker compose run migrate up

migrate-down:
	docker compose run migrate down 1

migrate-reset:
	docker compose run migrate down -all

migrate-version:
	docker compose run migrate version

migrate-force:
	docker compose run migrate force $(v)

migrate-create:
	docker run --rm -v ./migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq $(name)


# =========================
# Development
# =========================

run:
	go run ./cmd/api

test:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

lint:
	golangci-lint run

sqlc:
	sqlc generate

tidy:
	go mod tidy

# =========================
# DB
# =========================

db-shell:
	docker compose exec db psql -U ${DB_USER} -d ${DB_NAME}

# =========================
# Production build
# =========================

build-bin:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api ./cmd/api