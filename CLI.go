package poker

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const ComandoJogador = "Por favor insira o numero de jogadores: "
const MensagemErroVencedorInvalido = "valor errado obtido para vencedo, esperado formato 'NomeDoJogador venceu'"
const MensagemErroJogadorInvalido = "valor errado obtido para numero de jogadores, por favor tente de novo com um numero"

type CLI struct {
	in   *bufio.Scanner
	out  io.Writer
	jogo Jogo
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

func (cli *CLI) JogarPoker() {
	fmt.Fprint(cli.out, ComandoJogador)

	numeroDeJogadores, err := strconv.Atoi(cli.readLine())

	if err != nil {
		fmt.Fprint(cli.out, MensagemErroJogadorInvalido)
		return
	}

	cli.jogo.Iniciar(numeroDeJogadores)

	vencedorInput := cli.readLine()
	vencedor, err := extrairVencedor(vencedorInput)

	if err != nil {
		fmt.Fprint(cli.out, MensagemErroVencedorInvalido)
		return
	}

	cli.jogo.Terminar(vencedor)
}

func extrairVencedor(userInput string) (string, error) {
	if !strings.Contains(userInput, " venceu") {
		return "", errors.New(MensagemErroVencedorInvalido)
	}
	return strings.Replace(userInput, " venceu", "", 1), nil
}

// construtor

func NovoCLI(in io.Reader, out io.Writer, jogo Jogo) *CLI {
	return &CLI{
		in:   bufio.NewScanner(in),
		out:  out,
		jogo: jogo,
	}
}
