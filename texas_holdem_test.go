package poker_test

import (
	"os"
	"testing"
	"time"

	poker "github.com/marianavif/go-api"
)

func TestInicioDoJogo(t *testing.T) {
	t.Run("agenda alertas no início do jogo para 5 jogadores", func(t *testing.T) {
		alertadorDeBlind := &poker.EspiaoAlertadorDeBlind{}
		jogo := poker.NovoTexasHoldem(alertadorDeBlind, poker.DummyArmazenamentoJogador)

		jogo.Iniciar(5, os.Stdout)

		casos := []poker.AlertaAgendado{
			{As: 0 * time.Second, Quantidade: 100},
			{As: 10 * time.Minute, Quantidade: 200},
			{As: 20 * time.Minute, Quantidade: 300},
			{As: 30 * time.Minute, Quantidade: 400},
			{As: 40 * time.Minute, Quantidade: 500},
			{As: 50 * time.Minute, Quantidade: 600},
			{As: 60 * time.Minute, Quantidade: 800},
			{As: 70 * time.Minute, Quantidade: 1000},
			{As: 80 * time.Minute, Quantidade: 2000},
			{As: 90 * time.Minute, Quantidade: 4000},
			{As: 100 * time.Minute, Quantidade: 8000},
		}

		poker.VerificaCasosAgendados(t, casos, alertadorDeBlind)

	})

	t.Run("agenda alertas no início do jogo para 7 jogadores", func(t *testing.T) {
		alertadorDeBlind := &poker.EspiaoAlertadorDeBlind{}
		jogo := poker.NovoTexasHoldem(alertadorDeBlind, poker.DummyArmazenamentoJogador)

		jogo.Iniciar(7, os.Stdout)

		casos := []poker.AlertaAgendado{
			{As: 0 * time.Second, Quantidade: 100},
			{As: 12 * time.Minute, Quantidade: 200},
			{As: 24 * time.Minute, Quantidade: 300},
			{As: 36 * time.Minute, Quantidade: 400},
		}

		poker.VerificaCasosAgendados(t, casos, alertadorDeBlind)
	})
}

func TestTerminoDoJogo(t *testing.T) {
	armazenamento := &poker.EsbocoArmazenamentoJogador{}
	jogo := poker.NovoTexasHoldem(poker.DummyEspiaoAlertador, armazenamento)
	vencedor := "Ruth"

	jogo.Terminar(vencedor)
	poker.VerificaVitoriaJogador(t, armazenamento, vencedor)
}
