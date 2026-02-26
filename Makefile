APP_NAME=onde-tem
DB_URL=postgres://postgres:admin@db:5432/onde_tem?sslmode=disable

# =========================
# Docker
# =========================

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
	docker compose exec db psql -U $(user) -d $(db)

# =========================
# Production build
# =========================

build-bin:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api ./cmd/api