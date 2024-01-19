package poker

import (
	"bufio"
	"io"
	"strings"
)

type CLI struct {
	armazenamentoJogador ArmazenamentoJogador
	in                   *bufio.Scanner
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

func (cli *CLI) JogarPoker() {
	userInput := cli.readLine()
	cli.armazenamentoJogador.RegistrarVitoria(extrairVencedor(userInput))
}

// construtor

func NovoCLI(armazenamento ArmazenamentoJogador, in io.Reader) *CLI {
	return &CLI{
		armazenamentoJogador: armazenamento,
		in:                   bufio.NewScanner(in),
	}
}

func extrairVencedor(userInput string) string {
	return strings.Replace(userInput, " venceu", "", 1)
}
