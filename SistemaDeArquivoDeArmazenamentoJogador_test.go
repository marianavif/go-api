package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestSistemaDeArquivoDeArmazenamentoJogador(t *testing.T) {

	t.Run("/liga ordenada", func(t *testing.T) {
		bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, `[
			{"Nome": "Cleo", "Vitorias": 10},
			{"Nome": "Chris", "Vitorias": 33}
		]`)
		defer limpaBancoDeDados()

		armazenamento, err := NovoSistemaDeArquivoDeArmazenamentoJogador(bancoDeDados)

		verificaSemErro(t, err)

		recebido := armazenamento.ObterLiga()

		esperado := []Jogador{
			{"Chris", 33},
			{"Cleo", 10},
		}

		verificaLiga(t, recebido, esperado)

		// ler novamente
		recebido = armazenamento.ObterLiga()
		verificaLiga(t, recebido, esperado)
	})

	t.Run("pegar pontuação do jogador", func(t *testing.T) {
		bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, `[
			{"Nome": "Cleo", "Vitorias": 10},
			{"Nome": "Chris", "Vitorias": 33}
		]`)
		defer limpaBancoDeDados()

		armazenamento, err := NovoSistemaDeArquivoDeArmazenamentoJogador(bancoDeDados)

		verificaSemErro(t, err)

		recebido := armazenamento.ObterPontuacaoJogador("Chris")

		esperado := 33
		verificaPontuacaoIgual(t, recebido, esperado)
	})

	t.Run("armazena vitórias de um jogador existente", func(t *testing.T) {
		bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, `[
			{"Nome": "Cleo", "Vitorias": 10},
			{"Nome": "Chris", "Vitorias": 33}
		]`)
		defer limpaBancoDeDados()

		armazenamento, err := NovoSistemaDeArquivoDeArmazenamentoJogador(bancoDeDados)

		verificaSemErro(t, err)

		armazenamento.RegistrarVitoria("Chris")

		recebido := armazenamento.ObterPontuacaoJogador("Chris")
		esperado := 34
		verificaPontuacaoIgual(t, recebido, esperado)
	})

	t.Run("armazena vitorias de novos jogadores", func(t *testing.T) {
		bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, `[
			{"Nome": "Cleo", "Vitorias": 10},
			{"Nome": "Chris", "Vitorias": 33}
		]`)
		defer limpaBancoDeDados()

		armazenamento, err := NovoSistemaDeArquivoDeArmazenamentoJogador(bancoDeDados)

		verificaSemErro(t, err)

		armazenamento.RegistrarVitoria("Pepper")

		recebido := armazenamento.ObterPontuacaoJogador("Pepper")
		esperado := 1
		verificaPontuacaoIgual(t, recebido, esperado)
	})

	t.Run("funciona com um arquivo vazio", func(t *testing.T) {
		bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, "")
		defer limpaBancoDeDados()

		_, err := NovoSistemaDeArquivoDeArmazenamentoJogador(bancoDeDados)

		verificaSemErro(t, err)
	})
}

func verificaPontuacaoIgual(t *testing.T, recebido, esperado int) {
	t.Helper()
	if recebido != esperado {
		t.Errorf("recebido %d esperado %d", recebido, esperado)
	}
}

func criaArquivoTemporario(t *testing.T, dadoInicial string) (*os.File, func()) {
	t.Helper()

	arquivotmp, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("não foi possível escrever o arquivo temporário %v", err)
	}

	arquivotmp.Write([]byte(dadoInicial))

	removeArquivo := func() {
		arquivotmp.Close()
		os.Remove(arquivotmp.Name())
	}

	return arquivotmp, removeArquivo
}

func verificaSemErro(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("não esperava um erro mas obteve um, %v", err)
	}
}
