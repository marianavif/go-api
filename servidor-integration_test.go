package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegistrarVitoriasEBuscarEstasVitorias(t *testing.T) {

	jogador := "Maria"
	armazenamento := NovoArmazenamentoJogadorEmMemoria()
	servidor := NovoServidorJogador(armazenamento)

	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisicaoRegistrarVitoriaPost(jogador))
	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisicaoRegistrarVitoriaPost(jogador))
	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisicaoRegistrarVitoriaPost(jogador))

	t.Run("obter pontuação", func(t *testing.T) {
		resposta := httptest.NewRecorder()
		servidor.ServeHTTP(resposta, novaRequisicaoObterPontuacao(jogador))
		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)

		verificarCorpoRequisicao(t, resposta.Body.String(), "3")
	})

	t.Run("obter liga", func(t *testing.T) {
		resposta := httptest.NewRecorder()
		servidor.ServeHTTP(resposta, novaRequisicaoDeLiga())
		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)

		obtido := obterLigaDaResposta(t, resposta.Body)
		esperado := []Jogador{
			{"Maria", 3},
		}
		verificaLiga(t, obtido, esperado)
	})
}
