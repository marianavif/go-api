package poker_test

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"

	poker "github.com/marianavif/go-api"
)

func TestCLI(t *testing.T) {
	t.Run("inicia jogo com 3 jogadores e termina jogo com vencedor 'Chris'", func(t *testing.T) {
		jogo := &poker.EspiaoJogo{}
		saida := &bytes.Buffer{}

		in := usuarioEnvia("3", "Chris venceu")
		cli := poker.NovoCLI(in, saida, jogo)

		cli.JogarPoker()

		verificaMensagensEnviadasAoUsuario(t, saida, poker.ComandoJogador)
		verificaJogoIniciadoCom(t, jogo, 3)
		verificaTerminoChamadoCom(t, jogo, "Chris")

	})

	t.Run("inicia jogo com 8 jogadores and grava 'Cleo' como vencedora", func(t *testing.T) {
		jogo := &poker.EspiaoJogo{}

		in := usuarioEnvia("8", "Cleo venceu")
		cli := poker.NovoCLI(in, poker.DummySaida, jogo)

		cli.JogarPoker()

		verificaJogoIniciadoCom(t, jogo, 8)
		verificaTerminoChamadoCom(t, jogo, "Cleo")
	})

	t.Run("imprime um erro quando um valor nao numerico é inserido e nao inicia o jogo", func(t *testing.T) {
		jogo := &poker.EspiaoJogo{}

		saida := &bytes.Buffer{}
		in := usuarioEnvia("tortas")

		cli := poker.NovoCLI(in, saida, jogo)
		cli.JogarPoker()

		verificaJogoNaoIniciado(t, jogo)
		verificaMensagensEnviadasAoUsuario(t, saida, poker.ComandoJogador, poker.MensagemErroJogadorInvalido)
	})

	t.Run("imprime um erro quando o vencedor é declarado incorretamente", func(t *testing.T) {
		jogo := &poker.EspiaoJogo{}
		saida := &bytes.Buffer{}

		in := usuarioEnvia("8", "Lloyd é um assassino")
		cli := poker.NovoCLI(in, saida, jogo)

		cli.JogarPoker()

		verificaJogoNaoTerminado(t, jogo)
		verificaMensagensEnviadasAoUsuario(t, saida, poker.ComandoJogador, poker.MensagemErroVencedorInvalido)
	})
}

func verificaJogoNaoIniciado(t *testing.T, jogo *poker.EspiaoJogo) {
	t.Helper()
	if jogo.InicioChamado {
		t.Error("jogo nao deveria ter começado")
	}
}

func verificaJogoNaoTerminado(t *testing.T, jogo *poker.EspiaoJogo) {
	t.Helper()
	if jogo.TerminoChamado {
		t.Error("jogo nao deveria ter terminado")
	}
}

func verificaMensagensEnviadasAoUsuario(t *testing.T, saida *bytes.Buffer, mensagens ...string) {
	t.Helper()
	esperado := strings.Join(mensagens, "")
	recebido := saida.String()
	if recebido != esperado {
		t.Errorf("obteve '%s' enviado para saida mas esperava %+v", esperado, mensagens)
	}
}

func verificaJogoIniciadoCom(t *testing.T, jogo *poker.EspiaoJogo, numeroDeJogadores int) {
	t.Helper()

	passou := tentarNovamenteAte(500*time.Millisecond, func() bool {
		return jogo.InicioChamadoCom == numeroDeJogadores
	})

	if !passou {
		t.Errorf("queria Iniciar chamado com %d mas obteve %d", numeroDeJogadores, jogo.InicioChamadoCom)
	}
}

func verificaTerminoChamadoCom(t *testing.T, jogo *poker.EspiaoJogo, vencedor string) {
	t.Helper()

	passou := tentarNovamenteAte(500*time.Millisecond, func() bool {
		return jogo.TerminoChamadoCom == vencedor
	})

	if !passou {
		t.Errorf("esperava terminar chamando %s mas obteve %q", vencedor, jogo.TerminoChamadoCom)
	}
}

func tentarNovamenteAte(d time.Duration, f func() bool) bool {
	tempoLimite := time.Now().Add(d)
	for time.Now().Before(tempoLimite) {
		if f() {
			return true
		}
	}
	return false
}

func usuarioEnvia(mensagens ...string) io.Reader {
	return strings.NewReader(strings.Join(mensagens, "\n"))
}
