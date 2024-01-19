package poker

import (
	"encoding/json"
	"fmt"
	"io"
)

type Liga []Jogador

func (l Liga) EncontrarJogador(nome string) *Jogador {
	for i, j := range l {
		if j.Nome == nome {
			return &l[i]
		}
	}
	return nil
}

func NovaLiga(rdr io.Reader) (Liga, error) {
	var liga Liga
	err := json.NewDecoder(rdr).Decode(&liga)
	if err != nil {
		err = fmt.Errorf("problema parseando a liga, %v", err)
	}

	return liga, err
}
