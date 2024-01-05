package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegistrarVitoriasEBuscarEstasVitorias(t *testing.T) {
	t.Run("Em memória", func(t *testing.T) {
		jogador := "Maria"
		armazenamento := NovoArmazenamentoJogadorEmMemoria()
		servidor := ServidorJogador{armazenamento}

		servidor.ServeHTTP(httptest.NewRecorder(), novaRequisicaoRegistrarVitoriaPost(jogador))
		servidor.ServeHTTP(httptest.NewRecorder(), novaRequisicaoRegistrarVitoriaPost(jogador))
		servidor.ServeHTTP(httptest.NewRecorder(), novaRequisicaoRegistrarVitoriaPost(jogador))

		resposta := httptest.NewRecorder()
		servidor.ServeHTTP(resposta, novaRequisicaoObterPontuacao(jogador))
		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)

		verificarCorpoRequisicao(t, resposta.Body.String(), "3")
	})

}
