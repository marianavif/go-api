package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type EsbocoArmazenamentoJogadorDB struct {
	pontuacoes        ArmazenamentoJogadorEmDB
	registrosVitorias []string
}

func (e *EsbocoArmazenamentoJogadorDB) ObterPontuacaoJogador(nome string) int {
	db := ConectaArmazenamentoJogadorEmDB()
	pontuacao := pegaPontuacao(db, e.pontuacoes.tabela, nome)
	return pontuacao
}

func (e *EsbocoArmazenamentoJogadorDB) RegistrarVitoria(nome string) {
	e.registrosVitorias = append(e.registrosVitorias, nome)
}

func TestObterJogadoresDB(t *testing.T) {
	armazenamentoDB := EsbocoArmazenamentoJogadorDB{
		*CriaArmazenamentoJogadorEmDB("tabelateste"),
		nil,
	}
	servidorDB := &ServidorJogador{&armazenamentoDB}

	t.Run("retornar resultado de Maria da base de dados", func(t *testing.T) {
		requisicao := novaRequisicaoObterPontuacao("Maria")
		resposta := httptest.NewRecorder()

		servidorDB.ServeHTTP(resposta, requisicao)

		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)
		verificarCorpoRequisicao(t, resposta.Body.String(), "20")
	})

	t.Run("retornar resultado de Pedro da base de dados", func(t *testing.T) {
		requisicao := novaRequisicaoObterPontuacao("Pedro")
		resposta := httptest.NewRecorder()

		servidorDB.ServeHTTP(resposta, requisicao)

		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)
		verificarCorpoRequisicao(t, resposta.Body.String(), "10")
	})

	t.Run("retorna 404 para jogador não encontrado na base de dados", func(t *testing.T) {
		requisicao := novaRequisicaoObterPontuacao("Jorge")
		resposta := httptest.NewRecorder()

		servidorDB.ServeHTTP(resposta, requisicao)

		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusNotFound)
	})
}

func TestArmazenamentoVitoriasDB(t *testing.T) {
	armazenamentoDB := EsbocoArmazenamentoJogadorDB{
		*CriaArmazenamentoJogadorEmDB(zeraTabelaNaDB("tabelateste2")),
		nil,
	}
	servidorDB := &ServidorJogador{&armazenamentoDB}

	t.Run("registra vitorias na chamada ao método HTTP POST via base de dados", func(t *testing.T) {
		jogador := "Maria"

		requisicao := novaRequisicaoRegistrarVitoriaPost(jogador)
		resposta := httptest.NewRecorder()

		servidorDB.ServeHTTP(resposta, requisicao)

		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusAccepted)

		if len(armazenamentoDB.registrosVitorias) != 1 {
			t.Errorf("verifiquei %d chamadas a RegistrarVitoria, esperava %d", len(armazenamentoDB.registrosVitorias), 1)
		}

		if armazenamentoDB.registrosVitorias[0] != jogador {
			t.Errorf("não registrou o vencedor corretamente, recebi '%s', esperava '%s'", armazenamentoDB.registrosVitorias[0], jogador)
		}
	})
}

func zeraTabelaNaDB(tabela string) string {
	db := ConectaArmazenamentoJogadorEmDB()
	zeraTabela := fmt.Sprintf("delete from %s", tabela)
	_, err := db.Exec(zeraTabela)
	if err != nil {
		log.Println("Não foi possível criar uma tabela nova:", err)
	}
	return tabela
}
