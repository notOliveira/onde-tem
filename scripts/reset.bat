@echo off
cd /d "%~dp0"
cd ..

echo =========================================
echo  ☢️  NUCLEAR RESET STARTING
echo =========================================

echo.
echo 🧹 1/5 Destroying old environment...
docker compose down -v --rmi all --remove-orphans || (echo ❌ Teardown failed & exit /b 1)
echo ✅ Environment cleaned

echo.
echo 📦 2/5 Downloading dependencies...
go mod tidy || (echo ❌ Dependency download failed & exit /b 1)
echo ✅ Dependencies downloaded

echo.
echo 🏗️  3/5 Building fresh images...
docker compose build || (echo ❌ Build failed & exit /b 1)
echo ✅ Images built

echo.
echo 🐘 4/5 Starting database...
docker compose up -d db --wait || (echo ❌ Failed to start database & exit /b 1)
echo ✅ Database is up and healthy

echo.
echo 🛠️  5/5 Running migrations...
docker compose run --rm migrate up || (echo ❌ Migrations failed & exit /b 1)
echo ✅ Database migrated

echo.
echo =========================================
echo  ✅ ENVIRONMENT RESET COMPLETE
echo =========================================