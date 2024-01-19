package poker_test

import (
	"strings"
	"testing"

	poker "github.com/marianavif/go-api"
)

func TestCLI(t *testing.T) {

	t.Run("grava vencedor chris digitado pelo usuario", func(t *testing.T) {
		in := strings.NewReader("Chris venceu\n")
		armazenamentoJogador := &poker.EsbocoArmazenamentoJogador{}

		cli := poker.NovoCLI(armazenamentoJogador, in)
		cli.JogarPoker()

		poker.VerificaVitoriaJogador(t, armazenamentoJogador, "Chris")
	})

	t.Run("grava vencedor cleo digitado pelo usuario", func(t *testing.T) {
		in := strings.NewReader("Cleo venceu\n")
		armazenamentoJogador := &poker.EsbocoArmazenamentoJogador{}

		cli := poker.NovoCLI(armazenamentoJogador, in)
		cli.JogarPoker()

		poker.VerificaVitoriaJogador(t, armazenamentoJogador, "Cleo")
	})
}
