package main

import (
	"fmt"
	"log"
	"os"

	poker "github.com/marianavif/go-api"
)

const nomeArquivoBD = "jogo.db.json"

func main() {
	armazenamento, close, err := poker.SistemaDeArquivoDeArmazenamentoJogadorAPartirDeArquivo(nomeArquivoBD)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	jogo := poker.NovoTexasHoldem(poker.AlertadorDeBlindFunc(poker.StdOutAlertador), armazenamento)
	cli := poker.NovoCLI(os.Stdin, os.Stdout, jogo)

	fmt.Println("Vamos jogar poker")
	fmt.Println("Digite {Nome} venceu para registrar uma vitoria")

	cli.JogarPoker()
}
