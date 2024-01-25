package poker_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	poker "github.com/marianavif/go-api"
)

const tipoDoConteudoJSON = "application/json"

func TestObterJogadores(t *testing.T) {
	armazenamento := &poker.EsbocoArmazenamentoJogador{
		map[string]int{
			"Maria": 20,
			"Pedro": 10,
		},
		nil,
		nil,
	}

	servidor := poker.DeveFazerServidorJogador(t, armazenamento, poker.DummyJogo)
	t.Run("retornar resultado de Maria", func(t *testing.T) {
		requisicao := poker.NovaRequisicaoObterPontuacao("Maria")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		poker.VerificarRespostaCodigoStatus(t, resposta, http.StatusOK)
		poker.VerificarCorpoRequisicao(t, resposta.Body.String(), "20")
	})

	t.Run("retornar resultado de Pedro", func(t *testing.T) {
		requisicao := poker.NovaRequisicaoObterPontuacao("Pedro")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		poker.VerificarRespostaCodigoStatus(t, resposta, http.StatusOK)
		poker.VerificarCorpoRequisicao(t, resposta.Body.String(), "10")
	})

	t.Run("retorna 404 para jogador não encontrado", func(t *testing.T) {
		requisicao := poker.NovaRequisicaoObterPontuacao("Jorge")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		poker.VerificarRespostaCodigoStatus(t, resposta, http.StatusNotFound)
	})

}

func TestArmazenamentoVitorias(t *testing.T) {
	armazenamento := &poker.EsbocoArmazenamentoJogador{
		map[string]int{},
		nil,
		nil,
	}

	servidor := poker.DeveFazerServidorJogador(t, armazenamento, poker.DummyJogo)

	t.Run("registra vitorias na chamada ao método HTTP POST", func(t *testing.T) {
		jogador := "Maria"

		requisicao := poker.NovaRequisicaoRegistrarVitoriaPost(jogador)
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		poker.VerificarRespostaCodigoStatus(t, resposta, http.StatusAccepted)

		poker.VerificaVitoriaJogador(t, armazenamento, jogador)
	})
}

func TestLiga(t *testing.T) {
	armazenamento := &poker.EsbocoArmazenamentoJogador{}
	servidor := poker.DeveFazerServidorJogador(t, armazenamento, poker.DummyJogo)

	t.Run("retorna 200 em /liga", func(t *testing.T) {
		requisicao, _ := http.NewRequest(http.MethodGet, "/liga", nil)
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		var obtido []poker.Jogador

		err := json.NewDecoder(resposta.Body).Decode(&obtido)

		if err != nil {
			t.Fatalf("Não foi possível fazer parse da resposta do servidor '%s'no slice de Jogador, '%v' ", resposta.Body, err)
		}

		poker.VerificarRespostaCodigoStatus(t, resposta, http.StatusOK)
	})
	t.Run("retorna a tabela da Liga como JSON", func(t *testing.T) {
		ligaEsperada := []poker.Jogador{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		armazenamento := &poker.EsbocoArmazenamentoJogador{nil, nil, ligaEsperada}
		servidor := poker.DeveFazerServidorJogador(t, armazenamento, poker.DummyJogo)

		requisicao := poker.NovaRequisicaoDeLiga()
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		obtido := poker.ObterLigaDaResposta(t, resposta.Body)
		poker.VerificarRespostaCodigoStatus(t, resposta, http.StatusOK)
		poker.VerificaLiga(t, obtido, ligaEsperada)
		poker.VerificaTipoDoConteudo(t, resposta, tipoDoConteudoJSON)
	})

}

func TestJogo(t *testing.T) {
	t.Run("GET /jogo retorna 200", func(t *testing.T) {
		servidor := poker.DeveFazerServidorJogador(t, &poker.EsbocoArmazenamentoJogador{}, poker.DummyJogo)

		requisicao := poker.NovaRequisicaoDeJogo()
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		poker.VerificarRespostaCodigoStatus(t, resposta, http.StatusOK)
	})

	t.Run("começa um jogo com 3 jogadores, envia alguns alertas de blind no websocket e declara Ruth vencedora", func(t *testing.T) {
		alertaDeBlindEsperado := "Blind é 100"
		jogo := &poker.EspiaoJogo{AlertaDeBlind: []byte(alertaDeBlindEsperado)}
		vencedor := "Ruth"
		servidor := httptest.NewServer(poker.DeveFazerServidorJogador(t, poker.DummyArmazenamentoJogador, jogo))
		ws := poker.DeveConectarAoWebSocket(t, "ws"+strings.TrimPrefix(servidor.URL, "http")+"/ws")

		defer servidor.Close()
		defer ws.Close()

		poker.EscreverMensagemNoWebsocket(t, ws, "3")
		poker.EscreverMensagemNoWebsocket(t, ws, vencedor)

		verificaJogoIniciadoCom(t, jogo, 3)
		verificaTerminoChamadoCom(t, jogo, vencedor)

		poker.Within(t, 10*time.Millisecond, func() { poker.VerificaSeWebSocketObteveMensagem(t, ws, alertaDeBlindEsperado) })
	})
}
