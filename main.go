package main

import (
	"log"
	"net/http"
)

func main() {
	servidorDB := &ServidorJogador{CriaArmazenamentoJogadorEmDB("<nome-da-tabela-criada>")} // alterar para nome da tabela usada

	if err := http.ListenAndServe(":5000", servidorDB); err != nil {
		log.Fatalf("não foi possível escutar na porta 5000 %v", err)
	}
}
