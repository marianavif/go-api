package poker

type ArmazenamentoJogadorEmMemoria struct {
	armazenamento map[string]int
}

func (a *ArmazenamentoJogadorEmMemoria) ObterPontuacaoJogador(nome string) int {
	return a.armazenamento[nome]
}

func (a *ArmazenamentoJogadorEmMemoria) RegistrarVitoria(nome string) {
	a.armazenamento[nome]++
}

func (a *ArmazenamentoJogadorEmMemoria) ObterLiga() Liga {
	var liga Liga
	for nome, vitorias := range a.armazenamento {
		liga = append(liga, Jogador{nome, vitorias})
	}
	return liga
}
