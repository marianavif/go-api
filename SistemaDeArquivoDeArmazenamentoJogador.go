package poker

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type SistemaDeArquivoDeArmazenamentoJogador struct {
	bancoDeDados *json.Encoder
	liga         Liga
}

func (s *SistemaDeArquivoDeArmazenamentoJogador) ObterLiga() Liga {
	sort.Slice(s.liga, func(i, j int) bool {
		return s.liga[i].Vitorias > s.liga[j].Vitorias
	})
	return s.liga
}

func (s *SistemaDeArquivoDeArmazenamentoJogador) ObterPontuacaoJogador(nome string) int {
	jogador := s.liga.EncontrarJogador(nome)

	if jogador != nil {
		return jogador.Vitorias
	}

	return 0
}

func (s *SistemaDeArquivoDeArmazenamentoJogador) RegistrarVitoria(nome string) {
	jogador := s.liga.EncontrarJogador(nome)

	if jogador != nil {
		jogador.Vitorias++
	} else {
		s.liga = append(s.liga, Jogador{nome, 1})
	}

	s.bancoDeDados.Encode(s.liga)
}

func NovoSistemaDeArquivoDeArmazenamentoJogador(arquivo *os.File) (*SistemaDeArquivoDeArmazenamentoJogador, error) {

	err := iniciaArquivoBDDeJogador(arquivo)

	if err != nil {
		return nil, fmt.Errorf("problema inicializando arquivo do jogador, %v", err)
	}

	liga, err := NovaLiga(arquivo)

	if err != nil {
		return nil, fmt.Errorf("problema carregando o armazenamento do jogador de arquivo %s, %v", arquivo.Name(), err)
	}

	return &SistemaDeArquivoDeArmazenamentoJogador{
		bancoDeDados: json.NewEncoder(&fita{arquivo}),
		liga:         liga,
	}, nil
}

func SistemaDeArquivoDeArmazenamentoJogadorAPartirDeArquivo(path string) (*SistemaDeArquivoDeArmazenamentoJogador, func(), error) {
	db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, nil, fmt.Errorf("falha ao abrir %s %v", path, err)
	}

	closeFunc := func() {
		db.Close()
	}

	armazenamento, err := NovoSistemaDeArquivoDeArmazenamentoJogador(db)

	if err != nil {
		return nil, nil, fmt.Errorf("falha ao criar sistema de arquivos para armazenar jogadores, %v", err)
	}

	return armazenamento, closeFunc, nil
}

func iniciaArquivoBDDeJogador(arquivo *os.File) error {
	arquivo.Seek(0, 0)

	info, err := arquivo.Stat()

	if err != nil {
		return fmt.Errorf("problema ao usar o arquivo %s, %v", arquivo.Name(), err)
	}

	if info.Size() == 0 {
		arquivo.Write([]byte("[]"))
		arquivo.Seek(0, 0)
	}

	return nil
}
