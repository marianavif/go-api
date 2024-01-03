package main

import (
	"log"
	"net/http"
)

type ArmazenamentoJogadorEmMemoria struct{}

func (a *ArmazenamentoJogadorEmMemoria) ObterPontuacaoJogador(nome string) int {
	return 123
}

func (a *ArmazenamentoJogadorEmMemoria) RegistrarVitoria(nome string) {}

func main() {
	servidor := &ServidorJogador{&ArmazenamentoJogadorEmMemoria{}}

	if err := http.ListenAndServe(":5000", servidor); err != nil {
		log.Fatalf("não foi possível escutar na porta 5000 %v", err)
	}
}
