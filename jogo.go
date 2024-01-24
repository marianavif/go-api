package poker

import "io"

type Jogo interface {
	Iniciar(numeroDeJogadores int, destinoDosAlertas io.Writer)
	Terminar(vencedor string)
}
