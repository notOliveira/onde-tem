# Refatoração de Cache — Decorator Pattern

## Problema Atual

O cache está no Use Case (`list_establishments.go`, `create_establishment.go`), violando Hexagonal:
- Use cases conhecem detalhes de infraestrutura (cache)
- Serialização de domain entities falha (campos privados não deserializam)
- Cache logic duplicada em múltiplos use cases

## Solução: Cache Decorator

```
[Handler] → [UseCase] → [CacheDecorator] → [RepositoryImpl] → [DB]
                         ↓
                    implements
                 RepositoryInterface
```

### Princípio
- Use cases são pure business logic — não sabem que cache existe
- Cache é um detalhe de infraestrutura que envolve o repository
- Repository decorator implementa a mesma interface que o repository real
- Genérico e reutilizável para qualquer repository

---

## Estado: Concluído

### Etapa 1-3: Completo

- `dto/establishment_cache.go` — `EstablishmentCacheDTO` + `ToEstablishmentCacheDTO()`
- `dto/helpers.go` — `typeSliceToStrings()`
- `internal/adapters/outbound/cache/cached_repository.go` — `CachedRepository` (sem generics) com cache-aside
- `internal/core/usecase/list_establishments.go` — refatorado, só delega ao repo
- `internal/core/usecase/create_establishment.go` — `cache.Delete()` removido
- `domain/*.go` — MarshalJSON em todos os types
- `queries/establishments.sql` — query `ListEstablishments`
- `establishment_repository.go` — método `List()` implementado
- Handler e router — `GET /establishments` funcional

### Etapa 4: Completo

- `dto/establishment_db_row.go` — DTO comum com `*RowToDBRow` (retornam erro) + `RowToEstablishment` usando `domain.RehydrateEstablishment`
- `establishment_repository.go` — usa DTO helpers; arquivo caiu de 434 → 165 linhas
- `ports/cache.go` — adicionada `DeletePattern(ctx, pattern) error`
- `cache/valkey_client.go` — implementa `DeletePattern` via `SCAN` + `DEL` em batch de 100
- `cache/cached_repository.go` — implementa `GetBySlug`, `Update`, `Delete`; assertion de interface; TTL configurável; `filter.Types` ordenados alfabeticamente; invalida List keys em todas as mutações
- `app.go` — wire do `cachedRepo` com `cfg.CacheTTL`
- Use Cases — parâmetro `cache` removido

---

## Resumo dos Arquivos Modificados (Etapa 4)

| Arquivo | Ação |
|---------|------|
| `internal/adapters/dto/establishment_db_row.go` | **REESCRITO** — remove órfão, corrige rehidratação, propaga erros |
| `internal/adapters/outbound/persistence/postgres/establishment_repository.go` | **REESCRITO** — usa DTO helpers (434 → 165 linhas) |
| `internal/core/ports/cache.go` | **MODIFICADO** — adiciona `DeletePattern` |
| `internal/adapters/outbound/cache/valkey_client.go` | **MODIFICADO** — implementa `DeletePattern` (SCAN+DEL batch) |
| `internal/adapters/outbound/cache/cached_repository.go` | **REESCRITO** — completa interface, remove generics, usa `ttl` configurável |
| `internal/infra/app/app.go` | **MODIFICADO** — wire do `cachedRepo` |
| `internal/core/usecase/create_establishment.go` | **MODIFICADO** — remove `cache` |
| `internal/core/usecase/list_establishments.go` | **MODIFICADO** — remove `cache` da assinatura |

---

## Próximos Passos

- **Testes unitários** para `CachedRepository` (mock de `ports.EstablishmentRepository` e `ports.Cache`) — prioridade alta, decorator ainda não tem cobertura
- **Transactions em Create/Update** (Fase 1.3 do `.roadmap`)
- **Error wrapping consistente** (`fmt.Errorf("...: %w", err)`)
- **Structured logging** (zerolog/zap)
- **Health endpoint real** (`/health` com ping no DB e Valkey)
- **Slug factory com verificação de duplicidade**

---

## Conceitos-Chave

- **Decorator Pattern**: envolve um objeto existente sem alterar sua interface
- **Dependency Inversion**: use case depende de abstração (interface), não de implementação concreta
- **Single Responsibility**: use case faz lógica de negócio, decorator faz cache
- **Open/Closed**: adiciona cache sem modificar use cases existentes
- **Cache-Aside Pattern**: verificar cache → miss → buscar DB → guardar no cache
- **SCAN + DEL pattern**: invalidação de chaves por prefixo sem bloquear o servidor (vs `KEYS *` que é O(N) e bloqueante)

---

## Session Log

### 09/06/2026 — Conclusão da Etapa 4 (Decorator Pattern)

#### Objetivo

Concluir os Passos 4a-4d da refatoração de cache iniciada em 25-26/05/2026. O decorator já existia mas não estava wireado, o que significava que o cache não estava sendo usado em runtime — o bug original (cache hit retornando dados vazios) só foi resolvido parcialmente.

#### Estado herdado das sessões anteriores

- ✅ Etapas 1-3 completas (DTOs, decorator, use cases sem cache, query SQL, handler GET)
- ⚠️ 4a parcialmente feito: `dto/establishment_db_row.go` existia com 2 bugs
- ⚠️ 4b parcialmente feito: `List` usava `RehydrateEstablishments`, mas `GetByID`/`GetBySlug` não
- ❌ 4c-4d pendentes
- 🐛 Bug adicional: `CachedRepository` não implementava `Update`/`Delete`/`GetBySlug`, o que quebraria a build se fosse wireado

#### Bugs identificados no código herdado

1. **`establishment_db_row.go`** — função órfã `EstablishmentRowToDBRow` (nunca usada)
2. **`RowToEstablishment`** chamava `domain.NewEstablishment` que zerava ID, gerava `time.Now()` e ignorava email/website/timezone — quebrava timestamps do DB
3. **`unmarshalPhones`/`unmarshalAddress`** engoliam erros silenciosamente
4. **`generateListCacheKey`** não ordenava `filter.Types`, gerando cache miss para mesma query com ordem diferente
5. **`CachedRepository[T any]`** — generics eram overengineering (interface fixa, nunca usou `T`)
6. **Decorator incompleto** — faltavam 3 métodos da interface, bloquearia o wire

#### Decisões arquiteturais desta sessão

- **`DeletePattern` na interface `Cache`**: para o SCAN funcionar sem vazar Redis para o decorator (alternativas rejeitadas: type-assert feio, expor `Scan` na interface)
- **Ordenação alfabética de `filter.Types`**: evita cache miss por ordem, custo mínimo
- **SCAN + DEL em batch de 100**: invalidação de List keys não-bloqueante (vs `KEYS *`)
- **TTL configurável via construtor**: `time.Duration`, lido de `config.CacheTTL` no `app.go`
- **Remoção dos generics `[T any]`**: decorator não é reutilizável para outras interfaces no momento; YAGNI

#### Mudanças aplicadas

| Passo | Arquivo | Mudança |
|-------|---------|---------|
| 1 | `dto/establishment_db_row.go` | Remove órfão, corrige `RowToEstablishment` (usa `RehydrateEstablishment`), `*RowToDBRow` retornam erro, `unmarshal*` propagam erro |
| 2 | `establishment_repository.go` | 434→165 linhas. Remove `establishmentDBRow` local e `RehydrateEstablishments` exportado. `GetByID`/`GetBySlug`/`List` usam helpers do dto/ |
| 3 | `ports/cache.go` | Adiciona `DeletePattern(ctx, pattern) error` |
| 4 | `cache/valkey_client.go` | Implementa `DeletePattern` com SCAN iterator + DEL em batch de 100 |
| 5 | `cache/cached_repository.go` | Remove `[T any]`, adiciona assertion de interface, implementa `GetBySlug`/`Update`/`Delete`, ordena `filter.Types`, TTL via construtor, `invalidateListCache` em todas as mutações |
| 6 | `app.go` | Wire do `cachedRepo := cache.NewCachedRepository(repo, valkeyClient, time.Duration(cfg.CacheTTL)*time.Second)` |
| 7 | `create_establishment.go` | Remove `cache` da struct e do construtor |
| 8 | `list_establishments.go` | Remove `cache` da assinatura do construtor |

#### Comportamento final do cache

- **Leituras**: cache-aside em `GetByID` (chave `est-{id}`), `GetBySlug` (`est-slug:{slug}`), `List` (`est-list:{md5}`). Filtros com mesma semântica mas ordem diferente em `Types` geram a mesma chave.
- **Mutações** (`Create`/`Update`/`Delete`): invalidam chave do ID, chave do slug (quando conhecida) e **todas** as chaves `est-list:*` via SCAN.
- **TTL**: controlado por `config.CacheTTL` (em segundos), passado como `time.Duration` ao decorator.

#### Verificação

```
go build ./...     ✅
go vet ./...       ✅
gofmt -l .         ✅ (limpo)
go test ./...      ✅
sqlc generate      ✅ (idempotente)
```

#### Pendências identificadas (não resolvidas nesta sessão)

- `HealthUseCase` continua como stub, sem rota `/health`
- Repository não usa transações em Create/Update (Fase 1.3 do `.roadmap`)
- Sem error wrapping consistente
- Sem testes unitários para o `CachedRepository` (próxima prioridade sugerida)
- Hardcoded `cacheTTL` em `config.go` ainda é `int` (em segundos); poderia ser `time.Duration`

#### Lições

- **Decorator Pattern exigiu completude**: um decorator parcial não pode ser wireado, ou quebra a build, ou pior — compila mas não cacheia o que deveria
- **Hexagonal força disciplina**: `ports.Cache` receber `DeletePattern` foi a saída para manter o decorator agnóstico de Redis
- **Bugs latentes em código "morto"**: `RowToEstablishment` tinha bug grave (quebrava timestamps) que só apareceria se fosse usado — o tipo de bug que testes detectariam

#### Commits sugeridos

Ao final desta sessão, separaria em commits atômicos:

1. `refactor(dto): fix EstablishmentDBRow helpers, propagate unmarshal errors`
2. `refactor(repo): use DTO helpers in establishment_repository, remove duplication`
3. `feat(cache): add DeletePattern to Cache port and ValkeyClient`
4. `refactor(cache): complete CachedRepository to satisfy EstablishmentRepository`
5. `refactor(app): wire CachedRepository decorator with configurable TTL`
6. `refactor(usecase): remove unused cache parameter from use cases`
7. `docs(todo): mark etapa 4 as complete`

