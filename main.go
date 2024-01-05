package main

import (
	"log"
	"net/http"
)

func main() {

	servidor := &ServidorJogador{NovoArmazenamentoJogadorEmMemoria()}

	if err := http.ListenAndServe(":5000", servidor); err != nil {
		log.Fatalf("não foi possível escutar na porta 5000 %v", err)
	}

	// main para servidor com base de dados postgres
	/* 	servidorDB := &ServidorJogador{CriaArmazenamentoJogadorEmDB("<nome-da-tabela-criada>")} // alterar para nome da tabela usada

	if err := http.ListenAndServe(":5000", servidorDB); err != nil {
		log.Fatalf("não foi possível escutar na porta 5000 %v", err)
	} */
}
