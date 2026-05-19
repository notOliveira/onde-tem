# Roadmap de Evolução - Onde Tem?

## Visão Geral

Este documento contém o roadmap completo de evolução do projeto, organizado em fases de prioridade. O objetivo é transformar a aplicação em um sistema fácil de manter, escalável e moderno, respeitando Clean Architecture + Hexagonal Architecture.

---

## FASE 1 — Crítico (Bugs e Falhas)

### 1.1 Implementar `List()` corretamente
**Arquivo:** `internal/adapters/outbound/persistence/postgres/establishment_repository.go:310`

O método `List()` retorna `nil, nil` - é um bug grave que faz qualquer listagem falhar silenciosamente.

**Ação:**
- Implementar query builder para `EstablishmentFilter`
- Suportar filtros por `Types` (GIN index), `Search` (name/slug), `Limit`, `Offset`
- Retornar `[]*domain.Establishment` ou erro

**Referências:**
- `internal/core/domain/establishment_filter.go`
- `internal/core/ports/establishment_repository.go:14`

---

### 1.2 Corrigir Cache Read/Write Pattern
**Arquivo:** `internal/core/usecase/create_establishment.go`

O cache atual só faz `Delete` (linha 53) mas nunca lê. É inútil.

**Ação:**
- Implementar cache-aside pattern:
  1. `GetByID` → verificar cache primeiro → se miss, buscar no DB → store in cache
  2. `GetBySlug` → mesma lógica
  3. `List` → cachear resultado com TTL curto (1-5 min)
- Usar key pattern: `est-{id}`, `est-slug-{slug}`, `est-list-{filter-hash}`
- Tratar `redis.Nil` corretamente em `ValkeyClient.Get()`

**Referências:**
- `internal/adapters/outbound/cache/valkey_client.go`
- `internal/core/ports/cache.go`

---

### 1.3 Adicionar Transactions em Create/Update
**Arquivo:** `internal/adapters/outbound/persistence/postgres/establishment_repository.go`

`Create` e `Update` não usam transações PostgreSQL.

**Ação:**
- Envolver operações em `pgx.Tx`
- Criar `WithTx(ctx context.Context, tx pgx.Tx)` no repository
- Para `Create`: gerar UUID no código, não deixar o DB fazer (manter control)
- Para `Update`: usar `SELECT ... FOR UPDATE` se necessário

---

### 1.4 Eliminar Código Duplicado em Repository
**Arquivo:** `internal/adapters/outbound/persistence/postgres/establishment_repository.go`

`GetByID` e `GetBySlug` têm lógica idêntica de unmarshal e reidratação.

**Ação:**
- Extrair função privada `rowToEstablishment(row interface{}) (*domain.Establishment, error)`
- Reduzir linhas de ~50 para ~10 em cada método

---

## FASE 2 — Qualidade

### 2.1 Implementar Validation Library
**Arquivos:** DTOs e Handlers

**Ação:**
- Adicionar `go-playground/validator` ou `asaskevich/govalidator`
- Validar DTOs HTTP no handler antes de converter para domain
- Validar no use case como segunda barreira

**Exemplos de validação:**
- Email: formato válido, não vazio se fornecido
- Phone: formato E.164 ou similar
- Location: lat/lon dentro de ranges válidos
- Name: não vazio, tamanho razoável (2-255 chars)
- Slug: regex `[a-z0-9-]+`, sem leading/trailing hyphens

---

### 2.2 Structured Logging
**Arquivos:** `internal/infra/logger/logger.go`

**Ação:**
- Migrar para `zerolog` ou `zap`
- Logs em JSON para integração com ELK/Datadog
- Adicionar campos: `request_id`, `user_id` (se houver), `trace_id`

**Formato esperado:**
```json
{"level":"info","time":"2026-05-06T22:00:00Z","msg":"Establishment created","id":"uuid","slug":"name","duration_ms":45}
```

---

### 2.3 Error Wrapping Consistente
**Arquivos:** Todos os use cases e repositories

**Ação:**
- Todos os erros devem ser wrapped com contexto: `fmt.Errorf("CreateEstablishmentUseCase.Execute: %w", err)`
- Criar erros sentinel para cada camada:
  - `domain.Err*` — erros de domínio (já existem)
  - `repository.Err*` — erros de persistência
  - `usecase.Err*` — erros de lógica de negócio
- HTTP handler mapeia erros para status codes

**Mapeamento sugerido:**
| Erro | HTTP Status |
|------|-------------|
| `ErrEstablishmentNotFound` | 404 |
| `ErrInvalidName`, `ErrInvalidTypes`, etc | 422 |
| `ErrDuplicateSlug` | 409 |
| `repository.Err*` | 500 |
| `ErrInternal` | 500 |

---

### 2.4 Health Endpoint Real
**Arquivo:** `internal/core/usecase/health.go`

O health check atual é apenas um teste de cache.

**Ação:**
- Criar `/health` endpoint (GET)
- Verificar:
  - DB connectivity (`pool.Ping`)
  - Cache connectivity (`valkeyClient.Ping`)
- Retornar JSON estruturado:
```json
{
  "status": "healthy",
  "checks": {
    "database": "ok",
    "cache": "ok"
  }
}
```
- Criar `/ready` para Kubernetes readiness probe

---

### 2.5 Slug Factory com Verificação de Duplicidade
**Arquivo:** `internal/core/domain/slug.go`

**Ação:**
- `NewSlug` deve receber `EstablishmentRepository` (ou interface `SlugUniquenessChecker`)
- Verificar se slug já existe antes de retornar
- Criar método `IsUnique(ctx context.Context, slug string) (bool, error)` no repository

---

## FASE 3 — Escalabilidade

### 3.1 Paginação Cursor-Based
**Arquivos:** `internal/core/domain/establishment_filter.go`, `establishment_repository.go`

OFFSET é lento para tabelas grandes. Usar keyset pagination.

**Ação:**
- Adicionar `Cursor string` no filtro (base64 do último ID visto)
- Query: `WHERE id > :cursor ORDER BY id LIMIT :limit`
- Retornar `next_cursor` na resposta

**Resposta padrão:**
```json
{
  "data": [...],
  "pagination": {
    "next_cursor": "base64-encoded-id",
    "has_more": true
  }
}
```

---

### 3.2 CQRS Básico
**Arquivos:** `internal/core/ports/`, `adapters/outbound/persistence/`

Separar read (queries) de write (commands).

**Ação:**
- Manter `EstablishmentRepository` para commands (Create, Update, Delete)
- Criar `EstablishmentQueryRepository` para queries (GetByID, GetBySlug, List)
- Queries podem apontar para read replicas no futuro
- Use cases usam a interface appropriate

---

### 3.3 Read Replica Support
**Arquivo:** `docker-compose.yml`, `internal/infra/database/database.go`

**Ação:**
- Adicionar segundo container PostgreSQL como replica
- Configurar `DB_READ_HOST` / `DB_WRITE_HOST`
- Queries readonly usam `DB_READ_HOST`
- Commands usam `DB_WRITE_HOST`

---

### 3.4 OpenTelemetry Tracing
**Ação:**
- Adicionar `go.opentelemetry.io/otel`
- Middleware no Gin para tracing de requests
- Spans em:
  - HTTP handlers
  - Use cases
  - Repository methods
- Exportar para Jaeger/Zipkin

---

### 3.5 Value Objects Properly Implemented
**Arquivos:** `internal/core/domain/address.go`, `location.go`, `phone.go`

Address, Location e Phone devem ter semantics de Value Object.

**Ação:**
- `Address`, `Location`, `Phone` devem ser imutáveis após criação
- Métodos de comparação: `Equals(other Address) bool`
- `Phone` deve validar formato (E.164: `+55-11-99999-9999`)
- Considerar usar `ValueObject[T]` genérico do pacote `domain/valueobject`

---

## FASE 4 — Modernização

### 4.1 Graceful Shutdown
**Arquivo:** `cmd/api/main.go`

**Ação:**
- Capturar sinais `SIGTERM`, `SIGINT`
- Parar de aceitar novos requests
- Finalizar requests em andamento (com timeout)
- Fechar conexões: DB pool, cache, listener

**Implementação:**
```go
signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
<-ch
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
server.Shutdown(ctx)
app.Close()
```

---

### 4.2 Migrate Tool
**Arquivo:** `docker-compose.yml`

Migrar de `migrate/migrate` para `golang-migrate/migrate` ou `sql-migrate`.

**Ação:**
- Usar `golang-migrate/migrate` que é mais idiático em Go
- Manter compatibilidade com queries SQL existentes

---

### 4.3 Dependency Injection com Wire
**Arquivo:** `internal/infra/app/app.go`

**Ação:**
- Instalar `google/wire`
- Criar `wire.go` comwire.Inject
- Gerar código de DI automaticamente
- Eliminar inicialização manual em `NewApp()`

---

### 4.4 Event Sourcing Básico
**Arquivos:** `internal/core/domain/`, `internal/core/events/`

Emitir eventos de domínio para audit trail.

**Ação:**
- Criar pacote `events`
- Definir eventos: `EstablishmentCreated`, `EstablishmentUpdated`, `EstablishmentDeleted`
- Interface `EventPublisher` no ports
- Implementar com `github.com/segmentio/kafka-go` ou `nsq`

---

### 4.5 Input DTOs e Output DTOs
**Ação:**
- `CreateEstablishmentInput` → validado no handler
- `CreateEstablishmentOutput` → expõe apenas campos relevantes para o cliente
- Não expor entidades de domínio diretamente via API

---

### 4.6 API Versioning e Pagination Global
**Arquivo:** `internal/adapters/inbound/http/`

**Ação:**
- Padronizar resposta:
```json
{
  "data": [...],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  }
}
```
- Adicionar headers: `X-Pagination-Cursor`, `X-Total-Count`

---

## CHECKLIST DE IMPLEMENTAÇÃO

### Fase 1 — Crítico
- [ ] Implementar `List()` com filtros
- [ ] Corrigir cache-aside pattern (read + write)
- [ ] Adicionar transactions em Create/Update
- [ ] Extrair `rowToEstablishment()` privado

### Fase 2 — Qualidade
- [ ] Adicionar validation library
- [ ] Migrar para structured logging (zerolog/zap)
- [ ] Implementar error wrapping consistente
- [ ] Criar health endpoint real
- [ ] Slug factory com verificação de duplicidade

### Fase 3 — Escalabilidade
- [ ] Implementar cursor-based pagination
- [ ] Separar query e command repositories
- [ ] Configurar read replicas
- [ ] Adicionar OpenTelemetry tracing
- [ ] Implementar Value Objects corretamente

### Fase 4 — Modernização
- [ ] Implementar graceful shutdown
- [ ] Migrar para golang-migrate
- [ ] Configurar Wire para DI
- [ ] Adicionar event sourcing básico
- [ ] Separar Input/Output DTOs
- [ ] Padronizar API response format

---

## NOTAS

- **FASE 1 é obrigatória** antes de qualquer produção
- **FASE 2 pode ser feita em paralelo** com desenvolvimento de features
- **FASE 3 é para escalar** quando necessário
- **FASE 4 é técnico** e pode ser implementado gradualmente

---

## ERROS A SEREM CRIADOS

### Domain Errors (já existem em `errors.go`)
```
ErrInvalidName
ErrInvalidTypes
ErrInvalidLocation
ErrInvalidPhone
ErrInvalidAddress
ErrInvalidTimezone
ErrBlankEmailAndWebsite
ErrEstablishmentNotFound
```

### Novos erros a adicionar
```
ErrInvalidSlugFormat
ErrDuplicateSlug
ErrInvalidEmailFormat
ErrInvalidPhoneFormat
ErrInvalidCoordinate
```

---

## FONTES

- Clean Architecture: https://blog.cleancoder.com/uncle-bob/2012/08/08/the-clean-architecture.html
- Go Project Layout: https://github.com/golang-standards/project-layout
- Hexagonal Architecture: https://alistair.cockburn.us/hexagonal-architecture/
- Cache-Aside Pattern: https://docs.microsoft.com/en-us/azure/architecture/patterns/cache-aside
- Cursor Pagination: https://dev.to/przpiw/cursor-based-pagination-fgf
