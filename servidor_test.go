package poker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
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

	servidor := NovoServidorJogador(armazenamento)

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

		VerificarRespostaCodigoStatus(t, resposta, http.StatusOK)
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
		VerificarRespostaCodigoStatus(t, resposta, http.StatusOK)
		VerificaLiga(t, obtido, ligaEsperada)
		VerificaTipoDoConteudo(t, resposta, tipoDoConteudoJSON)
	})

}

func TestJogo(t *testing.T) {
	t.Run("GET /jogo retorna 200", func(t *testing.T) {
		servidor := NovoServidorJogador(&EsbocoArmazenamentoJogador{})

		requisicao := NovaRequisicaoDeJogo()
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		VerificarRespostaCodigoStatus(t, resposta, http.StatusOK)
	})

	t.Run("quando recebemos uma mensagem de um websocket que é vencedor do jogo", func(t *testing.T) {
		armazenamento := &EsbocoArmazenamentoJogador{}
		vencedor := "Ruth"
		servidor := httptest.NewServer(NovoServidorJogador(armazenamento))
		defer servidor.Close()

		wsURL := "ws" + strings.TrimPrefix(servidor.URL, "http") + "/ws"

		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("não foi possível abrir uma conexão de websocket em %s %v", wsURL, err)
		}
		defer ws.Close()

		if err := ws.WriteMessage(websocket.TextMessage, []byte(vencedor)); err != nil {
			t.Fatalf("não foi possível enviar mensagem na conexão websocket %v", err)
		}

		time.Sleep(10 * time.Millisecond)
		VerificaVitoriaJogador(t, armazenamento, vencedor)
	})
}
