package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type ArmazenamentoJogadorEmDB struct {
	tabela string
}

func ConectaArmazenamentoJogadorEmDB() *sql.DB {
	conexao := "user=postgres dbname=<nome-do-banco-de-dados> password=<senha-do-postgres> host=localhost sslmode=disable" // alterar nome do banco de dados e senha
	db, err := sql.Open("postgres", conexao)
	if err != nil {
		panic(err.Error())
	}
	return db

}

func CriaArmazenamentoJogadorEmDB(tabela string) *ArmazenamentoJogadorEmDB {
	return &ArmazenamentoJogadorEmDB{tabela}
}

func (a *ArmazenamentoJogadorEmDB) ObterPontuacaoJogador(nome string) int {
	db := ConectaArmazenamentoJogadorEmDB()
	pontuacao := pegaPontuacao(db, a.tabela, nome)
	defer db.Close()
	return pontuacao
}

func (a *ArmazenamentoJogadorEmDB) RegistrarVitoria(nome string) {
	db := ConectaArmazenamentoJogadorEmDB()
	pontuacao := pegaPontuacao(db, a.tabela, nome)
	pontuacao += 1
	updateJogador := fmt.Sprintf("update %s set Pontuacao = %d where Nome = '%s' ", a.tabela, pontuacao, nome)
	db.Exec(updateJogador)
	defer db.Close()
}

func pegaPontuacao(db *sql.DB, tabela, nome string) int {
	selectDeJogador := db.QueryRow(fmt.Sprintf("select * from %s where Nome = '%s'", tabela, nome))

	var nomeNaoUtilizado string
	var pontuacao int
	err := selectDeJogador.Scan(&nomeNaoUtilizado, &pontuacao)
	if err != nil {
		db.Exec(fmt.Sprintf("insert into %s (Nome, Pontuacao) values ('%s', 0)", tabela, nome))
		return 0
	}
	return pontuacao
}
