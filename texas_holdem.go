package poker

import "time"

type TexasHoldem struct {
	alertador     AlertadorDeBlind
	armazenamento ArmazenamentoJogador
}

func (j *TexasHoldem) Iniciar(numeroDeJogadores int) {
	incrementoDeBlind := time.Duration(5+numeroDeJogadores) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	tempoDeBlind := 0 * time.Second

	for _, blind := range blinds {
		j.alertador.AgendarAlertaAs(tempoDeBlind, blind)
		tempoDeBlind = tempoDeBlind + incrementoDeBlind
	}
}

func (j *TexasHoldem) Terminar(vencedor string) {
	j.armazenamento.RegistrarVitoria(vencedor)
}

// construtor

func NovoTexasHoldem(alertador AlertadorDeBlind, armazenamento ArmazenamentoJogador) *TexasHoldem {
	return &TexasHoldem{
		alertador:     alertador,
		armazenamento: armazenamento,
	}
}
