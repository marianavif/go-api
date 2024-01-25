package main

import (
	"log"
	"net/http"

	poker "github.com/marianavif/go-api"
)

const nomeArquivoBD = "jogo.db.json"

func main() {
	armazenamento, close, err := poker.SistemaDeArquivoDeArmazenamentoJogadorAPartirDeArquivo(nomeArquivoBD)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	jogo := poker.NovoTexasHoldem(poker.AlertadorDeBlindFunc(poker.Alertador), armazenamento)
	servidor, err := poker.NovoServidorJogador(armazenamento, jogo)

	if err != nil {
		log.Fatalf("problema ao criar o servidor do jogador %v", err)
	}

	if err := http.ListenAndServe(":5000", servidor); err != nil {
		log.Fatalf("Não foi possível escutar na porta 5000 %v", err)
	}
}
