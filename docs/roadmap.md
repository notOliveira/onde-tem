# Roadmap — Onde tem?

Plano de ação para evoluir o projeto do estado atual (Create + List com decorator de cache) até uma API pronta pra produção, testada e observável.

Cada fase tem:
- **Goal** — o que vamos entregar.
- **Scope** — arquivos afetados.
- **Out of scope** — o que decidimos deixar pra depois.
- **Validation** — como saber que funciona.

> **Convenção de commits** segue `docs/decisions/decisions.md` (test, feat, refactor, fix, chore, docs, build, perf, db, ci).

---

## Phase 0 — Realignment (já feito em sessões anteriores)

A refatoração do cache com Decorator Pattern está **concluída e wireada**. Histórico em `todo.md`.

---

## Phase 1 — Testar a fundação

**Por que primeiro**: o decorator (`CachedRepository`) tem zero cobertura de testes, mas é dono de toda a lógica de invalidação de cache. Refatorar qualquer coisa sem testes arrisca regressões silenciosas.

### 1.1 — Testes unitários do `CachedRepository`
- **Goal**: cobrir leituras com cache-aside + invalidação nas mutações.
- **Scope**:
  - `internal/adapters/outbound/cache/cached_repository_test.go`
  - Mocks: `MockEstablishmentRepository` e `MockCache` (ou `testify/mock` se ainda não for dep — senão, stubs manuais).
- **Casos**:
  - `GetByID` — miss → chama repo → popula cache.
  - `GetByID` — hit → repo **não** é chamado.
  - `GetByID` — JSON do cache corrompido → trata como miss, chama repo.
  - `GetBySlug` — simétrico ao `GetByID`.
  - `List` — `filter.Types` fora de ordem produzem a mesma chave (determinismo).
  - `List` — hit retorna todos; um DTO inválido → todos descartados → repo chamado.
  - `Create` — invalida `est-{id}`, `est-slug:{slug}` e `est-list:*`.
  - `Update` / `Delete` — mesmo contrato de invalidação.
  - TTL é repassado pro `cache.Set`.
- **Validation**: `go test ./internal/adapters/outbound/cache/... -race -count=1`.

### 1.2 — Testes unitários dos use cases
- **Goal**: travar o comportamento de negócio (regras de validação, transformações).
- **Scope**:
  - `internal/core/usecase/create_establishment_test.go`
  - `internal/core/usecase/list_establishments_test.go`
- **Casos**:
  - `Create` — nome vazio → `ErrInvalidName`.
  - `Create` — sem types → `ErrInvalidTypes`.
  - `Create` — location/address inválidos → erro propagado.
  - `Create` — happy path → `repo.Create` chamado uma vez com campos reidratados.
  - `List` — delega o filter tal qual pro repo.
- **Validation**: `go test ./internal/core/... -race`.

### 1.3 — Testes das entidades de domínio
- **Goal**: proteger invariantes dos value objects.
- **Scope**:
  - `internal/core/domain/establishment_test.go`
  - `slug_test.go`, `phone_test.go`, `address_test.go`, `location_test.go`
- **Casos**: caminhos de falha de cada construtor + alguns happy paths; `RehydrateEstablishment` não valida (preserva estado do DB).
- **Validation**: `go test ./internal/core/domain/...`.

---

## Phase 2 — Completar o CRUD via HTTP

**Por que segundo**: as rotas hoje só expõem Create + List. Update e Delete já estão wireados no domain, repository e decorator — só faltam handler + router.

### 2.1 — `PUT /api/v1/establishments/:id`
- **Scope**:
  - `internal/adapters/inbound/http/establishment_handler.go` — `HandleUpdate`.
  - `internal/core/usecase/update_establishment.go` — novo use case.
  - `internal/infra/api/router.go` — rota.
  - Body: update parcial (no mínimo: name, types, email, website, phones, location, address, timezone).
- **Decisões a tomar**:
  - **PATCH vs PUT**: PATCH é semanticamente correto pra update parcial. Vou com **PATCH** a menos que você queira um PUT de substituição completa.
  - **Semântica parcial**: atualizar só os campos presentes no body (campos ponteiro ou `*string` explícitos).
- **Validation**: `make sanity`.

### 2.2 — `DELETE /api/v1/establishments/:id`
- **Scope**:
  - `HandleDelete` handler.
  - `internal/core/usecase/delete_establishment.go`.
  - Entrada no router.
  - Retorna `204 No Content` no sucesso, `404` se não encontrado (precisa introduzir um erro "not found" no `domain/errors.go`).
- **Validation**: integração via `make sanity` + curl/Postman manual.

### 2.3 — `GET /api/v1/establishments/:id` e `GET /api/v1/establishments/slug/:slug`
- **Scope**: handlers pros métodos já existentes no repo. Útil pra clientes que buscam direto por ID ou slug.
- **Validation**: mesma da 2.2.

---

## Phase 3 — Observabilidade e health

**Por que**: um serviço que não te diz se está saudável é um serviço que não dá pra rodar em produção.

### 3.1 — `/health` de verdade
- **Scope**:
  - `internal/core/usecase/health.go` — ping no Postgres pool + Valkey.
  - Handler + rota (`GET /health`, **fora** de `/api/v1` — convenção comum).
  - Shape da resposta: `{ "status": "ok"|"degraded", "checks": { "db": ..., "cache": ... } }`.
- **Validation**: derrubar um container, bater no `/health`, confirmar que a resposta reflete a falha.

### 3.2 — Logging estruturado
- **Goal**: trocar logs ad-hoc por logs estruturados (zerolog ou zap).
- **Scope**:
  - Escolher uma lib (zerolog é leve; zap é mais rápido). Adicionar no `go.mod`.
  - Substituir o conteúdo de `internal/infra/logger/logger.go`.
  - Plugar no middleware de request (logger do Gin) com request_id.
- **Validation**: logs são JSON; um log por request com status, latency e rota.

### 3.3 — Error wrapping consistente
- **Goal**: erros carregam contexto subindo a stack (`fmt.Errorf("op: %w", err)`).
- **Scope**: varrer `repository.go`, `usecase/*.go`, `cached_repository.go`.
- **Validation**: error handler top-level no Gin mapeando sentinels do `domain/errors.go` pra status HTTP.

---

## Phase 4 — Transações e integridade dos dados

**Por que**: hoje, uma falha parcial no meio de um write pode deixar o DB inconsistente. O roadmap já tinha sinalizado isso como Fase 1.3.

### 4.1 — Create transacional
- **Scope**:
  - Repository aceita `pgx.Tx` (ou uma interface tipo `DBTX`) pra que múltiplos writes sejam atômicos.
  - `usecase.CreateEstablishment` abre tx, grava establishment + entidades relacionadas (phones etc.), commit ou rollback.
- **Sub-decisão**: vamos splittar `Establishment.Phones` em tabela própria agora, ou manter como JSONB? Se JSONB, transação fica simples; se normalizar, essa fase cresce.
- **Validation**: teste que força falha no meio do write e garante que nenhuma linha foi commitada.

### 4.2 — Update transacional
- Mesmo formato da 4.1 pros fluxos de update.

---

## Phase 5 — Segurança na geração de slug

**Por que**: hoje, dois estabelecimentos com nomes que slugificam igual colidem na UNIQUE constraint, e o erro sobe como erro genérico do Postgres pro cliente.

### 5.1 — Checagem de unicidade de slug no create
- **Scope**:
  - `domain.NewEstablishment` continua puro (sem dep de repo).
  - `usecase.CreateEstablishment`: se `rawSlug` é vazio, gera do nome e checa unicidade via repo; se tomado, sufixa `-2`, `-3`, etc.
  - Ou: manter o gerador de slug puro e tratar colisão dentro do repository / num `SlugService` dedicado.
- **Validation**: tentar criar dois establishments com o mesmo nome → o segundo leva slug sufixado, ambos passam.

---

## Phase 6 — CI/CD

### 6.1 — Workflow no GitHub Actions
- **Scope**: `.github/workflows/ci.yml`
- **Steps**: `go mod download`, `gofmt -l .`, `go vet ./...`, `golangci-lint run`, `go test -race ./...`, `sqlc generate` (verificar diff vazio).
- **Validation**: PR dispara CI; status check fica verde.

### 6.2 — Build & push de imagem
- **Scope**: workflow que builda a imagem Docker e sobe pra um registry (GHCR é grátis pra repos públicos).
- **Out of scope por ora**: step de deploy.

---

## Phase 7 — Paginação e semântica de listagem

**Por que**: o filter já tem `Limit`/`Offset`, mas não tem count, cursor nem total.

### 7.1 — Paginação por cursor
- **Scope**:
  - Trocar `Offset` por `Cursor` (opaco, base64 de `(created_at, id)`).
  - Retornar `next_cursor` e `has_more` na resposta.
- **Por que em vez de offset**: estável sob writes, mais barato pra tabelas grandes.

### 7.2 — Queries espaciais
- **PostGIS já está habilitado.** Expor filtro `?near=lat,lon&radius_m=...` no `GET /establishments`.
- **Scope**: filter struct ganha `Near *Location`, `RadiusM int`. Repository monta `ST_DWithin`.

---

## Phase 8 — Documentação e DX

### 8.1 — Spec OpenAPI
- **Scope**: `docs/openapi.yaml` descrevendo todas as rotas. Escrever à mão ou gerar a partir dos handlers Gin (ex.: `swaggo/swag`).
- **Validation**: spec valida; requests de exemplo fazem round-trip.

### 8.2 — Script de seed
- **Scope**: um programinha em Go ou arquivo SQL em `scripts/` que insere ~10 estabelecimentos de exemplo espalhados por tipos/locations pra dev local.
- **Validation**: `make seed` popula o DB; `GET /establishments?types=restaurant` retorna linhas.

---

## Out of scope (não-objetivos explícitos)

- **Autenticação / autorização** — ainda não tem spec.
- **Horários de funcionamento (business hours)** — mencionado em `decisions.md` mas sem schema/entidade ainda. Fica pra um doc futuro.
- **Enriquecimento colaborativo ("reviews", "tags")** — o README cita mas não tem modelo.
- **i18n** — `decisions.md` menciona `Accept-Language` mas não tem implementação.

Cada um desses merece um doc próprio quando virar prioridade.

---

## Ordem de execução sugerida

Se quiser o **caminho mais curto até uma v0.2 shippável**:

1. Phase 1.1 (testes do decorator) — ~meio dia
2. Phase 2.1 (Update) — meio dia
3. Phase 2.2 (Delete) — ¼ de dia
4. Phase 3.1 (`/health` de verdade) — ¼ de dia
5. Phase 6.1 (CI) — meio dia

Isso te dá um serviço CRUD-completo, com testes e CI em ~2 dias de trabalho focado.

As Phases 3.2, 3.3, 4, 5, 7 e 8 são qualidade-de-vida e podem ir entrando incrementalmente depois.