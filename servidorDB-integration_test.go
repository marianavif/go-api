package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegistrarVitoriasEBuscarEstasVitoriasDB(t *testing.T) {
	t.Run("Em base de dados", func(t *testing.T) {
		jogador := "Maria"
		armazenamento := CriaArmazenamentoJogadorEmDB(zeraTabelaNaDB("tabelateste2"))
		servidor := &ServidorJogador{armazenamento}

		servidor.ServeHTTP(httptest.NewRecorder(), novaRequisicaoRegistrarVitoriaPost(jogador))
		servidor.ServeHTTP(httptest.NewRecorder(), novaRequisicaoRegistrarVitoriaPost(jogador))
		servidor.ServeHTTP(httptest.NewRecorder(), novaRequisicaoRegistrarVitoriaPost(jogador))

		resposta := httptest.NewRecorder()
		servidor.ServeHTTP(resposta, novaRequisicaoObterPontuacao(jogador))
		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)

		verificarCorpoRequisicao(t, resposta.Body.String(), "3")
	})
}
