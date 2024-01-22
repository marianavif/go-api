package poker

type Jogo interface {
	Iniciar(numeroDeJogadores int)
	Terminar(vencedor string)
}