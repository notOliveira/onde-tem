# todo.md — Playbook de Testes Unitários

Documento de trabalho. **Não é documentação do projeto** — é referência pessoal pra conduzir a escrita dos testes.

> Nota: este arquivo NÃO conflita com `docs/todo.md` (que é sobre a refatoração de cache, concluída).
> `docs/roadmap.md` continua sendo o roadmap macro; este é focado no **escopo de testes**.

---

## Status atual

```
Setup
  [x] testify adicionado ao go.mod (// direct)
  [x] Exemplo 100% detalhado escrito (establishment_id_test.go, 4 subtests verdes)

Nível 1 — Value Objects
  [x] establishment_id_test.go          (4/4 subtests verdes)
  [ ] establishment_type_test.go
  [ ] phone_test.go
  [ ] address_test.go
  [ ] location_test.go        ← primeiro table-driven test
  [ ] slug_test.go            ← mais denso (regex + Unicode)

Nível 2 — Entidade Establishment (RehydrateEstablishment, getters, MarshalJSON)
Nível 3 — Use Cases (create, list, health) com mocks das ports
Nível 4 — DTOs (phone_dto, location_dto, address_dto, establishment_cache, establishment_db_row)
Nível 5 — Cached Repository (decorator: cache-aside, invalidação, generateListCacheKey)
Nível 6 — HTTP Handlers (Gin, httptest)
```

---

## 1. Setup já feito (não repetir)

```bash
go get github.com/stretchr/testify     # promove // indirect → // direct
go mod tidy                             # limpa go.sum
```

Convenção de pacote: **sempre `package <nome>_test`** (externo, sufixo `_test`).
Ganha: força usar API pública, pega bugs de exposição e impede bypass de construtores.

---

## 2. Anatomia do arquivo de teste (template-base)

Arquivo de referência: `internal/core/domain/establishment_id_test.go` (já no projeto).
Cole abaixo como esqueleto e adapte:

```go
package <nome>_test

import (
    "testing"

    "github.com/notOliveira/onde-tem/internal/core/<área>"
    "github.com/stretchr/testify/assert"
)

func Test<X>(t *testing.T) {
    t.Run("<cenário em inglês, sem snake_case se possível>", func(t *testing.T) {
        // Act
        got, err := <área>.<X>(inputs)

        // Assert #1: erro esperado (sentinel)
        assert.ErrorIs(t, err, <área>.<ErrX>)

        // Assert #2: zero value do tipo em caso de erro
        assert.Equal(t, "<zero>", string(got))
    })
    // ... mais subtests: caminho feliz, getters, MarshalJSON
}
```

### Padrões de testify (decorar)

| Padrão | Quando usar |
|---|---|
| `assert.ErrorIs(t, err, sentinel)` | Erros sentinela (`ErrInvalid*`). Tolera wrap via `errors.Is`. |
| `assert.NoError(t, err)` | Pré-condição antes de asserts sobre o valor. |
| `assert.Equal(t, esperado, atual)` | Comparação de valor. Testify reporta "esperado X, obtido Y". |
| `assert.NotNil(t, got)` | Quando o retorno é struct/ponteiro e você só quer garantir não-zero. |
| `t.Run("nome", func(t *testing.T) { ... })` | Subtest isolado; um falha, outros rodam; report nomeado. |

### Convenções

- Nome de subtest em **inglês**, em `snake_case` ou frase curta.
- Ordem dos subtests dentro de uma `TestX`: **erro → feliz → variações → MarshalJSON**.
- Espere `esperado` na **esquerda** (`assert.Equal(t, "foo", got)`) — convention p/ report legível.
- Imports agrupados: stdlib / third-party / internos, separados por linha em branco.

---

## 3. Anti-armadilhas

| Erro comum | Por que quebra | Solução |
|---|---|---|
| `assert.Equal(t, err, sentinel)` | Quebra com `fmt.Errorf("...: %w", err)` | Use `assert.ErrorIs` |
| `assert.EqualError(t, err, "msg")` | Acopla a frase de log | Use `assert.ErrorIs` |
| `assert.Equal(t, "foo", string(data))` em MarshalJSON | Compara string decodificada | Use `[]byte(\`"foo"\`)` ou `assert.JSONEq` |
| `assert.Same` achando que é `assert.Equal` | Compara ponteiros, não valores | Quase nunca é o que você quer em VO |
| `package domain` (interno) | Permite `phone.number = ""`, bypass do construtor | Use `package domain_test` |
| Esquecer `import "testing"` | Falha de build no `*testing.T` | Sempre importar |
| Esquecer de rodar `go mod tidy` após `go get` | `// indirect` em vez de `// direct` | Roda tidy sempre |

---

## 4. Roteiro dos 6 níveis (ordem pedagógica)

Cada nível introduz UM conceito novo. Não pular.

### Nível 1 — Value Objects 🟢 (rodada atual)
- `establishment_id`, `establishment_type`, `phone`, `address`, `location`, `slug`
- Conceito: nada novo; absorver testify + subtests
- Ponto alto: `slug.go` (regex + Unicode), `location.go` (primeiro table-driven)

### Nível 2 — Entidade 🟡
- `NewEstablishment`, `UpdateName`, `UpdateContact`, `AddPhone`, `RemovePhone`,
  `UpdateLocation`, `UpdateAddress`, `UpdateTimezone`, `SetTypes`, `MarshalJSON`, `RehydrateEstablishment`
- Conceito novo: **table-driven tests** consolidado

### Nível 3 — Use Cases 🟠
- `create_establishment`, `list_establishments`, `health`
- Conceito novo: **mocks** (manuais ou `testify/mock`)

### Nível 4 — DTOs 🟧
- `phone_dto`, `location_dto`, `address_dto`, `establishment_cache`, `establishment_db_row`
- Conceitos: `t.Run` subtests, fixtures, `t.Cleanup()`

### Nível 5 — Cached Repository 🔴
- Comportamento cache-aside, invalidação, `generateListCacheKey`
- Conceito novo: **spies** (verificar argumentos passados ao mock)

### Nível 6 — HTTP Handlers 🟣
- Testes com `httptest`, `gin.SetMode(gin.TestMode)`, mock de use cases

### ❌ Fora de escopo
- Testes de banco (precisariam `testcontainers-go` ou `pgxmock`)
- Testes do Redis/Valkey (precisariam `miniredis`)
- Logger / Config (código trivial)

---

## 5. Detalhamento Nível 1 (checklist por arquivo)

### 5.1 `establishment_id_test.go` ✅ FEITO
Fonte: `internal/core/domain/establishment_id.go`
- [x] `ParseEstablishmentID("")` → erro `ErrInvalidEstablishmentID`
- [x] `ParseEstablishmentID("qualquer-coisa")` → retorna ID com mesmo valor, `nil`
- [x] `id.String()` retorna valor original
- [x] `MarshalJSON` produz JSON com aspas (`[]byte("\"foo\"")`)

### 5.2 `establishment_type_test.go`
Fonte: `internal/core/domain/establishment_type.go`
- Funções expostas: `String`
- Constante: `EstablishmentTypeRestaurant == "restaurant"`
- [ ] `EstablishmentType("restaurant").String()` → `"restaurant"`
- [ ] Constante `EstablishmentTypeRestaurant` tem valor `"restaurant"`
- Dica: arquivo curto, ótimo pra se acostumar com `t.Run`.

### 5.3 `phone_test.go`
Fonte: `internal/core/domain/phone.go` (`NewPhone`, `CountryCode`, `Number`, `Label`, `MarshalJSON`)
Erro: `ErrInvalidPhone`
- [ ] `NewPhone("", "", "")` → erro
- [ ] `NewPhone("", "  ", "")` → erro (trim → vazio)
- [ ] `NewPhone("+55", "11999998888", "WhatsApp")` → sucesso, getters ok
- [ ] `NewPhone("  +55  ", "  11999998888  ", "  Zap  ")` → sucesso, getters retornam **trimados**
- [ ] `MarshalJSON` produz `{"countryCode":"...","number":"...","label":"..."}`

### 5.4 `address_test.go`
Fonte: `internal/core/domain/address.go`
Funções: `NewAddress`, `Street`, `Number`, `District`, `City`, `State`, `Country`, `ZipCode`, `MarshalJSON`
Erro: `ErrInvalidAddress`
Campos obrigatórios: `street`, `city`, `state`, `country`, `zipCode` (vide `address.go:19`).
Opcionais (podem ser vazios): `number`, `district`.
- [ ] Todos obrigatórios vazios → erro
- [ ] Só `street` preenchido → erro (faltam os outros)
- [ ] `street=" "` → erro (trim)
- [ ] Todos obrigatórios preenchidos → sucesso
- [ ] Getters retornam valor (não struct zero)
- [ ] `MarshalJSON` produz `{"street":..,"number":..,"district":..,"city":..,"state":..,"country":..,"zipCode":..}`

### 5.5 `location_test.go` ← primeiro table-driven
Fonte: `internal/core/domain/location.go`
Funções: `NewLocation`, `IsValid`, `Lat`, `Lon`, `MarshalJSON`
Erro: `ErrInvalidLocation`
- [ ] `IsValid()` nos limites: `-90/90`, `-180/180` → `true` (table)
- [ ] Lat fora: `91`, `-91` → `false` (table)
- [ ] Lon fora: `181`, `-181` → `false` (table)
- [ ] `NewLocation(0, 0)` → sucesso, getters retornam 0
- [ ] `NewLocation(91, 0)` → erro `ErrInvalidLocation`
- [ ] `MarshalJSON` produz `{"lat":...,"lon":...}`

Estrutura table-driven alvo:
```go
tests := []struct{
    name string
    lat, lon float64
    wantErr bool
}{
    {"limit lat upper",    90,   0,   false},
    {"limit lat lower",   -90,   0,   false},
    {"lat above range",    91,   0,   true},
    ...
}
for _, tc := range tests {
    t.Run(tc.name, func(t *testing.T) { ... })
}
```

### 5.6 `slug_test.go` (mais denso)
Fonte: `internal/core/domain/slug.go`
Funções: `NewSlug`, `GenerateSlugFromName`, `String`, `MarshalJSON`
Erro: `ErrInvalidSlugFormat`
Regex: `^[a-z0-9-]+$` (sem acentos, sem underscores, sem maiúsculas, sem espaços)

`NewSlug`:
- [ ] vazio / só espaços → erro
- [ ] com maiúsculas → erro (`"Foo"`)
- [ ] com underscore → erro (`"foo_bar"`)
- [ ] com caractere especial → erro (`"café!"`)
- [ ] válido simples → sucesso
- [ ] válido com hífen no meio → sucesso
- [ ] com espaços nas pontas → sucesso e trimmed

`GenerateSlugFromName` (NÃO retorna erro; pode dar `""`):
- [ ] `"Restaurante Bom"` → `"restaurante-bom"`
- [ ] `"Café Açaí"` → `"cafe-acai"` (testa remoção via `unicode.Mn`)
- [ ] `"  Pizza  "` → `"pizza"` (trim)
- [ ] `"Loja 123"` → `"loja-123"`
- [ ] `"!!!@#"` → `""` (vazio sem erro)
- [ ] `" "` → `""`

`removeAccents` (não-exportada): testada indiretamente via `GenerateSlugFromName`.

`String()` e `MarshalJSON`:
- [ ] `Slug("ola").String()` → `"ola"`
- [ ] `MarshalJSON` de slug válido → `"ola"` (string JSON com aspas)

---

## 6. Detalhamento Níveis 2-6 (resumo — preencher conforme avança)

### Nível 2 — `establishment_test.go`
- Cobrir cada método de mutação com 1 caso feliz + 1 caso de erro
- `RehydrateEstablishment` **não valida** (preserva estado do DB) — documentar com comentário no teste
- `MarshalJSON` produz objeto com todos os campos agregados

### Nível 3 — `*_usecase_test.go`
- Mock manual via struct ou `testify/mock.Mock`
- Injetar mocks via construtor dos use cases
- Casos: validação rejeita entrada, happy path chama repo com payload correto, erro do repo propaga

### Nível 4 — DTOs
- `t.Run` para cada campo do round-trip DB → domain → cache DTO
- `t.Cleanup()` se necessário

### Nível 5 — `cached_repository_test.go`
- Spies via `testify/mock`: verificar que `cache.Set` foi chamado com a chave certa e TTL certo
- Casos: cache miss → busca no repo → popula cache; cache hit → não chama repo; mutações invalidam chaves

### Nível 6 — Handlers
- `httptest.NewRecorder()` + `gin.SetMode(gin.TestMode)`
- Mock dos use cases
- Cobrir 200/400/404/500 por rota

---

## 7. Comandos úteis

```bash
# Rodar testes do pacote inteiro (verbose, com lista de subtests)
go test -v ./internal/core/domain/

# Rodar só um teste específico
go test -v -run TestParseEstablishmentID ./internal/core/domain/

# Rodar só um subtest (note a barra)
go test -v -run TestParseEstablishmentID/returns_error_when_value_is_empty ./internal/core/domain/

# Rodar tudo no projeto
go test ./...

# Com detecção de race conditions (use a partir do Nível 3)
go test -race ./internal/core/usecase/

# Formatar após editar
gofmt -w <arquivo>

# Verificar formatação sem modificar
gofmt -l .

# Vet básico
go vet ./...
```

---

## 8. Workflow por arquivo

1. Você cria o `*_test.go` (copy do esqueleto da seção 2 + casos da seção 5/6).
2. Roda `go test -v -run Test<Nome> ./...`. Se vermelho, corrijo contigo.
3. Quando os subtests do arquivo passam, segue pro próximo.
4. Ao final de cada nível, roda `go test ./...` pra garantir que nada quebrou em outros pacotes.

---

## 9. Referência rápida dos arquivos do projeto

| Arquivo de produção | O que tem | Teste correspondente |
|---|---|---|
| `internal/core/domain/establishment_id.go` | `ParseEstablishmentID`, `MarshalJSON`, `String` | ✅ `establishment_id_test.go` |
| `internal/core/domain/establishment_type.go` | `String` + constante `EstablishmentTypeRestaurant` | `establishment_type_test.go` |
| `internal/core/domain/phone.go` | `NewPhone` (com trim) + getters + `MarshalJSON` | `phone_test.go` |
| `internal/core/domain/address.go` | `NewAddress` (5 obrigatórios) + getters + `MarshalJSON` | `address_test.go` |
| `internal/core/domain/location.go` | `NewLocation` (faixa lat/lon) + getters + `MarshalJSON` | `location_test.go` |
| `internal/core/domain/slug.go` | `NewSlug` (regex), `GenerateSlugFromName` (Unicode), `removeAccents` | `slug_test.go` |
| `internal/core/domain/establishment.go` | `NewEstablishment`, mutadores, `MarshalJSON`, `RehydrateEstablishment` | `establishment_test.go` (Nível 2) |
| `internal/core/domain/establishment_filter.go` | Filtros de listagem | `establishment_filter_test.go` (Nível 2) |
| `internal/core/usecase/*.go` | Use cases | `*_usecase_test.go` (Nível 3) |
| `internal/adapters/dto/*.go` | Mapeamentos DB/cache ↔ domain | `*_dto_test.go` (Nível 4) |
| `internal/adapters/outbound/cache/cached_repository.go` | Decorator cache-aside | `cached_repository_test.go` (Nível 5) |
| `internal/adapters/inbound/http/*.go` | Handlers Gin | `*_handler_test.go` (Nível 6) |

---

## 10. Quando travar / dúvida

- **Conceitual** (ex: por que ErrorIs e não EqualError): olha a Seção 3 (Anti-armadilhas).
- **Qual arquivo cobrir**: olha a Seção 9 (tabela de referência).
- **Caso de teste faltando**: olha a Seção 5 (checklist por arquivo).
- **Esquecer o que cada função faz**: ler o próprio arquivo `.go` antes de testar.
- **Build quebrou**: rodar `go mod tidy` primeiro; senão conferir imports.
