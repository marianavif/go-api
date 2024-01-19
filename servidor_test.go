package poker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const tipoDoConteudoJSON = "application/json"

func TestObterJogadores(t *testing.T) {
	armazenamento := &EsbocoArmazenamentoJogador{
		map[string]int{
			"Maria": 20,
			"Pedro": 10,
		},
		nil,
		nil,
	}

	servidor := NovoServidorJogador(armazenamento)
	t.Run("retornar resultado de Maria", func(t *testing.T) {
		requisicao := NovaRequisicaoObterPontuacao("Maria")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		VerificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)
		VerificarCorpoRequisicao(t, resposta.Body.String(), "20")
	})

	t.Run("retornar resultado de Pedro", func(t *testing.T) {
		requisicao := NovaRequisicaoObterPontuacao("Pedro")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		VerificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)
		VerificarCorpoRequisicao(t, resposta.Body.String(), "10")
	})

	t.Run("retorna 404 para jogador não encontrado", func(t *testing.T) {
		requisicao := NovaRequisicaoObterPontuacao("Jorge")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		VerificarRespostaCodigoStatus(t, resposta.Code, http.StatusNotFound)
	})

}

func TestArmazenamentoVitorias(t *testing.T) {
	armazenamento := &EsbocoArmazenamentoJogador{
		map[string]int{},
		nil,
		nil,
	}

	servidor := NovoServidorJogador(armazenamento)

	t.Run("registra vitorias na chamada ao método HTTP POST", func(t *testing.T) {
		jogador := "Maria"

		requisicao := NovaRequisicaoRegistrarVitoriaPost(jogador)
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		VerificarRespostaCodigoStatus(t, resposta.Code, http.StatusAccepted)

		VerificaVitoriaJogador(t, armazenamento, jogador)
	})
}

func TestLiga(t *testing.T) {
	armazenamento := &EsbocoArmazenamentoJogador{}
	servidor := NovoServidorJogador(armazenamento)

	t.Run("retorna 200 em /liga", func(t *testing.T) {
		requisicao, _ := http.NewRequest(http.MethodGet, "/liga", nil)
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		var obtido []Jogador

		err := json.NewDecoder(resposta.Body).Decode(&obtido)

		if err != nil {
			t.Fatalf("Não foi possível fazer parse da resposta do servidor '%s'no slice de Jogador, '%v' ", resposta.Body, err)
		}

		VerificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)
	})
	t.Run("retorna a tabela da Liga como JSON", func(t *testing.T) {
		ligaEsperada := []Jogador{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		armazenamento := &EsbocoArmazenamentoJogador{nil, nil, ligaEsperada}
		servidor := NovoServidorJogador(armazenamento)

		requisicao := NovaRequisicaoDeLiga()
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		obtido := ObterLigaDaResposta(t, resposta.Body)
		VerificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)
		VerificaLiga(t, obtido, ligaEsperada)
		VerificaTipoDoConteudo(t, resposta, tipoDoConteudoJSON)
	})

}
