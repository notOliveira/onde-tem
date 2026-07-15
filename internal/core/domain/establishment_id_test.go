// =============================================================================
// ARQUIVO DIDÁTICO — referência para os próximos testes do projeto.
// Comentários densos aqui têm propósito pedagógico. Nos próximos _test.go,
// mantenha apenas comentários curtos explicando o PORQUÊ de decisões não-óbvias.
// =============================================================================

// Pacote EXTERNO de teste (sufixo _test). Contraste com "package domain":
//
//   - "package domain"          → você é interno, vê campos minúsculos e tudo.
//   - "package domain_test"     → você é externo, SÓ vê o que começa em maiúscula.
//
// Por que externo aqui?
//  1. Os Value Objects têm campos não-exportados (ex: `number` no Phone).
//     Se o teste fosse interno, você escreveria `phone.number = ""` direto,
//     bypassando o construtor NewPhone — e aí o teste não validaria nada útil.
//  2. Externo nos força a usar só a API pública (NewPhone, getters, MarshalJSON).
//     Se um dia o campo virar `numberStr`, o teste pega.
//  3. Detecta bugs de EXPOSIÇÃO de API: se uma função necessária não existe
//     (em maiúscula), o teste não compila.
package domain_test

import (
	// "testing" é o pacote da stdlib que define *testing.T, *testing.B, *testing.M.
	// É ele que dá o tipo passado como parâmetro de toda função TestX(t *testing.T).
	// Todo arquivo _test.go importa testing — sempre.
	"testing"

	// Pacote sendo testado. Imports SEMPRE pelo path completo do módulo.
	// Em Go, group imports é convenção: stdlib primeiro, depois third-party,
	// depois internos — separados por linha em branco.
	"github.com/notOliveira/onde-tem/internal/core/domain"

	// testify: assertions/mocks. "stretchr/testify" é uma organização no GitHub.
	// O módulo importável é "github.com/stretchr/testify", e dentro dele temos
	// subpacotes (assert, mock, require, suite). Aqui usamos só "assert".
	"github.com/stretchr/testify/assert"
)

// Convenção Go para nomes de funções de teste:
//
//	func TestX(t *testing.T) { ... }
//
// Onde:
//   - "Test" — prefixo obrigatório. O runner do Go procura por isso.
//   - "X"     — em CamelCase, descrevendo O QUE está sendo testado.
//     Aqui: TestParseEstablishmentID cobre o construtor + getters.
//   - t *testing.T — parâmetro obrigatório. É seu canal pro runner: reportar
//     falhas (t.Fail), rodar subtests (t.Run), pular (t.Skip), etc.
//
// Subtests via t.Run isolam cada cenário dentro da mesma função.
// Vantagem: se um falhar, os outros rodam, e o report mostra QUAL cenário quebrou.
func TestParseEstablishmentID(t *testing.T) {

	// ----------------------------------------------------------------------
	// SUBTEST 1: caso de erro — entrada vazia deve falhar com sentinel error
	// ----------------------------------------------------------------------
	//
	// Nome do subtest descreve o CENÁRIO de negócio, não o que o código faz.
	// O report mostra: `--- FAIL: TestParseEstablishmentID/returns_error_when_value_is_empty`
	// Lendo só o nome, você sabe qual regra foi violada.
	t.Run("returns error when value is empty", func(t *testing.T) {
		// Act: a ação sob teste. Mantemos em uma linha para legibilidade.
		id, err := domain.ParseEstablishmentID("")

		// Assert #1 — o erro deve ser o SENTINEL (público, em domain.ErrInvalidEstablishmentID).
		//
		// Por que ErrorIs e não outras opções:
		//   assert.Equal(t, err, domain.ErrInvalidEstablishmentID)
		//     → comparação por identidade. QUEBRA se um dia o código fizer
		//       fmt.Errorf("...: %w", ErrInvalidEstablishmentID) (wrap).
		//
		//   assert.EqualError(t, err, "invalid establishment id")
		//     → compara a MENSAGEM literal. QUEBRA se alguém mudar a frase.
		//       Ruim: testes acoplados a texto de UI/log.
		//
		//   assert.ErrorIs(t, err, domain.ErrInvalidEstablishmentID)
		//     → navega a chain via errors.Is(). Tolera wrap. É o certo.
		assert.ErrorIs(t, err, domain.ErrInvalidEstablishmentID)

		// Assert #2 — quando há erro, o valor retornado deve ser o ZERO VALUE.
		// EstablishmentID tem underlying string, então zero value é "".
		// "string(id)" converte pra comparar com literal.
		assert.Equal(t, "", string(id))
	})

	// ----------------------------------------------------------------------
	// SUBTEST 2: caso feliz — entrada válida retorna ID sem erro
	// ----------------------------------------------------------------------
	t.Run("returns id without error for valid value", func(t *testing.T) {
		id, err := domain.ParseEstablishmentID("abc-123")

		// NoError é açúcar pra "assert.Nil(t, err)" com mensagem padrão melhor.
		// É boa prática chamar ANTES dos asserts sobre o valor — se err != nil,
		// os asserts seguintes ficam sem sentido.
		assert.NoError(t, err)

		// Convertendo pra string para comparar com literal.
		// Idiomaticamente: assert.Equal(t, domain.EstablishmentID("abc-123"), id)
		// também funciona — ambas as formas estão corretas, escolha uma e seja
		// consistente. Aqui preferimos string(id) por legibilidade para quem
		// não conhece o type subjacente.
		assert.Equal(t, "abc-123", string(id))
	})

	// ----------------------------------------------------------------------
	// SUBTEST 3: getter String() preserva o valor
	// ----------------------------------------------------------------------
	// Pode parecer exagero testar um getter de 1 linha, mas String() é
	// "contrato de exibição" — se alguém amanhã trocar `return string(id)`
	// por outra lógica, esse teste pega. Custo: 3 linhas. Vale.
	t.Run("String returns the underlying value", func(t *testing.T) {
		id, err := domain.ParseEstablishmentID("foo")

		// Precondição: se parse falhar aqui, o teste abaixo não faz sentido.
		// Em testify você poderia usar "require.NoError" — esse chama t.FailNow()
		// e PARA o subtest. Mas assert.NoError só marca falha e continua.
		// Como a próxima linha não depende de err, qualquer um serve.
		assert.NoError(t, err)

		assert.Equal(t, "foo", id.String())
	})

	// ----------------------------------------------------------------------
	// SUBTEST 4: MarshalJSON produz string JSON com aspas
	// ----------------------------------------------------------------------
	// EstablishmentID.MarshalJSON chama json.Marshal(id.String()).
	// json.Marshal de uma string retorna o literal ENTRE ASPAS.
	// Então id.MarshalJSON() deve retornar exatamente: []byte(`"foo"`)
	//
	// Por que não testar o JSON inteiro (ex: assert.JSONEq)?
	//   Porque esse tipo serializa COMO string JSON, não como objeto.
	//   O contrato a testar é "é uma string com aspas", não "é igual a X".
	t.Run("MarshalJSON wraps value in JSON quotes", func(t *testing.T) {
		id, err := domain.ParseEstablishmentID("foo")
		assert.NoError(t, err)

		// MarshalJSON retorna ([]byte, error). Atribuímos os dois.
		data, err := id.MarshalJSON()
		assert.NoError(t, err)

		// Comparação de slices de bytes.
		// []byte(`"foo"`) é uma forma compacta e correta — backticks em Go
		// criam raw string literals, então as aspas duplas não precisam de escape.
		//
		// Forma equivalente com string normal (escape):
		//   []byte("\"foo\"")
		//
		// assert.Equal sobre slices compara elemento a elemento — seguro.
		assert.Equal(t, []byte(`"foo"`), data)
	})
}
