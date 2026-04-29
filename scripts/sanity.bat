@echo off
cd /d "%~dp0"
cd ..

echo ================================
echo   🔍 SANITY CHECK STARTING
echo ================================

echo.
echo 🐳 Checking Docker...
docker ps > NUL 2>&1 || (echo ❌ Docker is NOT running & exit /b 1)
echo ✅ Docker is running

echo.
echo 🐘 Checking Postgres container...
docker compose ps db --format "{{.State}}" | findstr "running" > NUL || (echo ❌ Postgres container is NOT running & exit /b 1)
echo ✅ Postgres container is up

echo.
echo ⚙️  Running sqlc generate...
sqlc generate || (echo ❌ sqlc failed & exit /b 1)
echo ✅ sqlc OK

echo.
echo 🧪 Running tests...
go test -v ./... || (echo ❌ Tests failing & exit /b 1)

echo.
echo ================================
echo   ✅ ENVIRONMENT HEALTHY
echo ================================