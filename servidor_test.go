package poker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

	servidor := DeveFazerServidorJogador(t, armazenamento)
	t.Run("retornar resultado de Maria", func(t *testing.T) {
		requisicao := NovaRequisicaoObterPontuacao("Maria")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		VerificarRespostaCodigoStatus(t, resposta, http.StatusOK)
		VerificarCorpoRequisicao(t, resposta.Body.String(), "20")
	})

	t.Run("retornar resultado de Pedro", func(t *testing.T) {
		requisicao := NovaRequisicaoObterPontuacao("Pedro")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		VerificarRespostaCodigoStatus(t, resposta, http.StatusOK)
		VerificarCorpoRequisicao(t, resposta.Body.String(), "10")
	})

	t.Run("retorna 404 para jogador não encontrado", func(t *testing.T) {
		requisicao := NovaRequisicaoObterPontuacao("Jorge")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		VerificarRespostaCodigoStatus(t, resposta, http.StatusNotFound)
	})

}

func TestArmazenamentoVitorias(t *testing.T) {
	armazenamento := &EsbocoArmazenamentoJogador{
		map[string]int{},
		nil,
		nil,
	}

	servidor := DeveFazerServidorJogador(t, armazenamento)

	t.Run("registra vitorias na chamada ao método HTTP POST", func(t *testing.T) {
		jogador := "Maria"

		requisicao := NovaRequisicaoRegistrarVitoriaPost(jogador)
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		VerificarRespostaCodigoStatus(t, resposta, http.StatusAccepted)

		VerificaVitoriaJogador(t, armazenamento, jogador)
	})
}

func TestLiga(t *testing.T) {
	armazenamento := &EsbocoArmazenamentoJogador{}
	servidor := DeveFazerServidorJogador(t, armazenamento)

	t.Run("retorna 200 em /liga", func(t *testing.T) {
		requisicao, _ := http.NewRequest(http.MethodGet, "/liga", nil)
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		var obtido []Jogador

		err := json.NewDecoder(resposta.Body).Decode(&obtido)

		if err != nil {
			t.Fatalf("Não foi possível fazer parse da resposta do servidor '%s'no slice de Jogador, '%v' ", resposta.Body, err)
		}

		VerificarRespostaCodigoStatus(t, resposta, http.StatusOK)
	})
	t.Run("retorna a tabela da Liga como JSON", func(t *testing.T) {
		ligaEsperada := []Jogador{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		armazenamento := &EsbocoArmazenamentoJogador{nil, nil, ligaEsperada}
		servidor := DeveFazerServidorJogador(t, armazenamento)

		requisicao := NovaRequisicaoDeLiga()
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		obtido := ObterLigaDaResposta(t, resposta.Body)
		VerificarRespostaCodigoStatus(t, resposta, http.StatusOK)
		VerificaLiga(t, obtido, ligaEsperada)
		VerificaTipoDoConteudo(t, resposta, tipoDoConteudoJSON)
	})

}

func TestJogo(t *testing.T) {
	t.Run("GET /jogo retorna 200", func(t *testing.T) {
		servidor := DeveFazerServidorJogador(t, &EsbocoArmazenamentoJogador{})

		requisicao := NovaRequisicaoDeJogo()
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		VerificarRespostaCodigoStatus(t, resposta, http.StatusOK)
	})

	t.Run("quando recebemos uma mensagem de um websocket que é vencedor do jogo", func(t *testing.T) {
		armazenamento := &EsbocoArmazenamentoJogador{}
		vencedor := "Ruth"
		servidor := httptest.NewServer(DeveFazerServidorJogador(t, armazenamento))
		defer servidor.Close()

		wsURL := "ws" + strings.TrimPrefix(servidor.URL, "http") + "/ws"

		ws := DeveConectarAoWebSocket(t, wsURL)
		defer ws.Close()

		EscreverMensagemNoWebsocket(t, ws, vencedor)

		time.Sleep(10 * time.Millisecond)
		VerificaVitoriaJogador(t, armazenamento, vencedor)
	})
}
