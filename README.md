# Onde tem? / Find it

System to find information about establishments, enriched with collaborative informations.

## Build the application (with Docker)
```bash
> docker compose up --build
# Stop
> Ctrl + C
```

## Autoformat code
```bash
> gofmt -w .
```

## Access Postgres Database using Docker
```
docker exec -it onde-tem-db psql -U <USER> -d <DATABASE>
```

## Removing containers
```
docker compose down -v
```