package main

type ArmazenamentoJogadorEmMemoria struct {
	armazenamento map[string]int
}

func (a *ArmazenamentoJogadorEmMemoria) ObterPontuacaoJogador(nome string) int {
	return a.armazenamento[nome]
}

func (a *ArmazenamentoJogadorEmMemoria) RegistrarVitoria(nome string) {
	a.armazenamento[nome]++
}

func NovoArmazenamentoJogadorEmMemoria() *ArmazenamentoJogadorEmMemoria {
	return &ArmazenamentoJogadorEmMemoria{map[string]int{}}
}
