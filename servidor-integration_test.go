package poker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegistrarVitoriasEBuscarEstasVitorias(t *testing.T) {

	jogador := "Maria"
	bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, `[]`)
	defer limpaBancoDeDados()
	armazenamento, err := NovoSistemaDeArquivoDeArmazenamentoJogador(bancoDeDados)
	verificaSemErro(t, err)
	servidor := DeveFazerServidorJogador(t, armazenamento, DummyJogo)

	servidor.ServeHTTP(httptest.NewRecorder(), NovaRequisicaoRegistrarVitoriaPost(jogador))
	servidor.ServeHTTP(httptest.NewRecorder(), NovaRequisicaoRegistrarVitoriaPost(jogador))
	servidor.ServeHTTP(httptest.NewRecorder(), NovaRequisicaoRegistrarVitoriaPost(jogador))

	t.Run("obter pontuação", func(t *testing.T) {
		resposta := httptest.NewRecorder()
		servidor.ServeHTTP(resposta, NovaRequisicaoObterPontuacao(jogador))
		VerificarRespostaCodigoStatus(t, resposta, http.StatusOK)

		VerificarCorpoRequisicao(t, resposta.Body.String(), "3")
	})

	t.Run("obter liga", func(t *testing.T) {
		resposta := httptest.NewRecorder()
		servidor.ServeHTTP(resposta, NovaRequisicaoDeLiga())
		VerificarRespostaCodigoStatus(t, resposta, http.StatusOK)

		obtido := ObterLigaDaResposta(t, resposta.Body)
		esperado := []Jogador{
			{"Maria", 3},
		}
		VerificaLiga(t, obtido, esperado)
	})
}
