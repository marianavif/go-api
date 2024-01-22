package poker

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

var DummyEspiaoAlertador = &EspiaoAlertadorDeBlind{}
var DummyArmazenamentoJogador = &EsbocoArmazenamentoJogador{}
var DummyStdIn = &bytes.Buffer{}
var DummyStdOut = &bytes.Buffer{}

type EsbocoArmazenamentoJogador struct {
	pontuacoes        map[string]int
	RegistrosVitorias []string
	liga              Liga
}

type AlertaAgendado struct {
	As         time.Duration
	Quantidade int
}

type EspiaoAlertadorDeBlind struct {
	Alertas []AlertaAgendado
}

func (e *EspiaoAlertadorDeBlind) AgendarAlertaAs(duracao time.Duration, quantidade int) {
	e.Alertas = append(e.Alertas, AlertaAgendado{duracao, quantidade})
}

func (a *AlertaAgendado) String() string {
	return fmt.Sprintf("%d toca às %v", a.Quantidade, a.As)
}

func (e *EsbocoArmazenamentoJogador) ObterPontuacaoJogador(nome string) int {
	pontuacao := e.pontuacoes[nome]
	return pontuacao
}

func (e *EsbocoArmazenamentoJogador) RegistrarVitoria(nome string) {
	e.RegistrosVitorias = append(e.RegistrosVitorias, nome)
}

func (e *EsbocoArmazenamentoJogador) ObterLiga() Liga {
	return e.liga
}

func NovaRequisicaoObterPontuacao(nome string) *http.Request {
	requisicao, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/jogadores/%s", nome), nil)
	return requisicao
}

func NovaRequisicaoRegistrarVitoriaPost(nome string) *http.Request {
	requisicao, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/jogadores/%s", nome), nil)
	return requisicao
}

func VerificarRespostaCodigoStatus(t *testing.T, recebido, esperado int) {
	t.Helper()
	if recebido != esperado {
		t.Errorf("não recebeu código de status HTTP esperado, recebido %d, esperado %d", recebido, esperado)
	}
}

func VerificarCorpoRequisicao(t *testing.T, recebido, esperado string) {
	t.Helper()
	if recebido != esperado {
		t.Errorf("corpo da requisição é inválido, obtive '%s', esperava '%s' ", recebido, esperado)
	}
}

func ObterLigaDaResposta(t *testing.T, body io.Reader) []Jogador {
	t.Helper()
	liga, _ := NovaLiga(body)
	return liga
}

func VerificaLiga(t *testing.T, obtido, esperado []Jogador) {
	t.Helper()
	if !reflect.DeepEqual(obtido, esperado) {
		t.Errorf("obtido %v esperado %v", obtido, esperado)
	}
}

func NovaRequisicaoDeLiga() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/liga", nil)
	return req
}

func VerificaTipoDoConteudo(t *testing.T, resposta *httptest.ResponseRecorder, esperado string) {
	t.Helper()
	if resposta.Result().Header.Get("content-type") != esperado {
		t.Errorf("resposta não obteve content-type de %s, obtido %v", esperado, resposta.Result().Header)
	}
}

func VerificaVitoriaJogador(t *testing.T, armazenamento *EsbocoArmazenamentoJogador, jogador string) {
	t.Helper()

	if len(armazenamento.RegistrosVitorias) != 1 {
		t.Errorf("verifiquei %d chamadas a RegistrarVitoria, esperava %d", len(armazenamento.RegistrosVitorias), 1)
	}

	if armazenamento.RegistrosVitorias[0] != jogador {
		t.Errorf("não registrou o vencedor corretamente, recebi '%s', esperava '%s'", armazenamento.RegistrosVitorias[0], jogador)
	}
}

func VerificaAlertaAgendado(t *testing.T, recebido, esperado AlertaAgendado) {
	t.Helper()

	if recebido.Quantidade != esperado.Quantidade {
		t.Errorf("obtida quantidade %d, esperado %d", recebido.Quantidade, esperado.Quantidade)
	}

	if recebido.As != esperado.As {
		t.Errorf("obteve tempo agendado de %v, esperado %v", recebido.As, esperado.As)
	}
}

func VerificaCasosAgendados(t *testing.T, casos []AlertaAgendado, alertadorDeBlind *EspiaoAlertadorDeBlind) {
	t.Helper()

	for i, esperado := range casos {
		t.Run(fmt.Sprint(esperado), func(t *testing.T) {

			if len(alertadorDeBlind.Alertas) <= i {
				t.Fatalf("alerta %d não foi agendado %v", i, alertadorDeBlind.Alertas)
			}

			alerta := alertadorDeBlind.Alertas[i]

			VerificaAlertaAgendado(t, alerta, esperado)
		})
	}
}