package poker_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	poker "github.com/marianavif/go-api"
)

type EspiaoJogo struct {
	InicioChamadoCom  int
	TerminoChamadoCom string
	InicioChamado     bool
	TerminoChamado    bool
}

func (e *EspiaoJogo) Iniciar(numeroDeJogadores int) {
	e.InicioChamado = true
	e.InicioChamadoCom = numeroDeJogadores
}

func (e *EspiaoJogo) Terminar(vencedor string) {
	e.TerminoChamado = true
	e.TerminoChamadoCom = vencedor
}

func TestCLI(t *testing.T) {
	t.Run("inicia jogo com 3 jogadores e termina jogo com vencedor 'Chris'", func(t *testing.T) {
		jogo := &EspiaoJogo{}
		stdout := &bytes.Buffer{}

		in := usuarioEnvia("3", "Chris venceu")
		cli := poker.NovoCLI(in, stdout, jogo)

		cli.JogarPoker()

		verificaMensagensEnviadasAoUsuario(t, stdout, poker.ComandoJogador)
		verificaJogoIniciadoCom(t, jogo, 3)
		verificaTerminoChamadoCom(t, jogo, "Chris")

	})

	t.Run("inicia jogo com 8 jogadores and grava 'Cleo' como vencedora", func(t *testing.T) {
		jogo := &EspiaoJogo{}

		in := usuarioEnvia("8", "Cleo venceu")
		cli := poker.NovoCLI(in, poker.DummyStdOut, jogo)

		cli.JogarPoker()

		verificaJogoIniciadoCom(t, jogo, 8)
		verificaTerminoChamadoCom(t, jogo, "Cleo")
	})

	t.Run("imprime um erro quando um valor nao numerico é inserido e nao inicia o jogo", func(t *testing.T) {
		jogo := &EspiaoJogo{}

		stdout := &bytes.Buffer{}
		in := usuarioEnvia("tortas")

		cli := poker.NovoCLI(in, stdout, jogo)
		cli.JogarPoker()

		verificaJogoNaoIniciado(t, jogo)
		verificaMensagensEnviadasAoUsuario(t, stdout, poker.ComandoJogador, poker.MensagemErroJogadorInvalido)
	})

	t.Run("imprime um erro quando o vencedor é declarado incorretamente", func(t *testing.T) {
		jogo := &EspiaoJogo{}
		stdout := &bytes.Buffer{}

		in := usuarioEnvia("8", "Lloyd é um assassino")
		cli := poker.NovoCLI(in, stdout, jogo)

		cli.JogarPoker()

		verificaJogoNaoTerminado(t, jogo)
		verificaMensagensEnviadasAoUsuario(t, stdout, poker.ComandoJogador, poker.MensagemErroVencedorInvalido)
	})
}

func verificaJogoNaoIniciado(t *testing.T, jogo *EspiaoJogo) {
	t.Helper()
	if jogo.InicioChamado {
		t.Error("jogo nao deveria ter começado")
	}
}

func verificaJogoNaoTerminado(t *testing.T, jogo *EspiaoJogo) {
	t.Helper()
	if jogo.TerminoChamado {
		t.Error("jogo nao deveria ter terminado")
	}
}

func verificaMensagensEnviadasAoUsuario(t *testing.T, stdout *bytes.Buffer, mensagens ...string) {
	t.Helper()
	esperado := strings.Join(mensagens, "")
	recebido := stdout.String()
	if recebido != esperado {
		t.Errorf("obteve '%s' enviado para stdout mas esperava %+v", esperado, mensagens)
	}
}

func verificaJogoIniciadoCom(t *testing.T, jogo *EspiaoJogo, numeroDeJogadores int) {
	t.Helper()

	if jogo.InicioChamadoCom != numeroDeJogadores {
		t.Errorf("queria Iniciar chamado com %d mas obteve %d", numeroDeJogadores, jogo.InicioChamadoCom)
	}
}

func verificaTerminoChamadoCom(t *testing.T, jogo *EspiaoJogo, vencedor string) {
	t.Helper()

	if jogo.TerminoChamadoCom != vencedor {
		t.Errorf("esperava terminar chamando %s mas obteve %q", vencedor, jogo.TerminoChamadoCom)
	}
}

func usuarioEnvia(mensagens ...string) io.Reader {
	return strings.NewReader(strings.Join(mensagens, "\n"))
}
