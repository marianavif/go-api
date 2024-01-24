package poker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"
)

const jsonTipoDeConteudo = "application/json"
const caminhoTemplateHTML = "jogo.html"

type ArmazenamentoJogador interface {
	ObterPontuacaoJogador(nome string) int
	RegistrarVitoria(nome string)
	ObterLiga() Liga
}

type ServidorJogador struct {
	armazenamento ArmazenamentoJogador
	http.Handler
	template *template.Template
}

type Jogador struct {
	Nome     string
	Vitorias int
}

var atualizadorDeWebsocket = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *ServidorJogador) ManipulaLiga(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonTipoDeConteudo)
	json.NewEncoder(w).Encode(s.armazenamento.ObterLiga())
	w.WriteHeader(http.StatusOK)
}

func (s *ServidorJogador) ManipulaJogadores(w http.ResponseWriter, r *http.Request) {
	jogador := r.URL.Path[len("/jogadores/"):]

	switch r.Method {
	case http.MethodPost:
		s.registrarVitoria(w, jogador)
	case http.MethodGet:
		s.mostrarPontuacao(w, jogador)
	}
}

func (s *ServidorJogador) mostrarPontuacao(w http.ResponseWriter, jogador string) {
	pontuacao := s.armazenamento.ObterPontuacaoJogador(jogador)

	if pontuacao == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, pontuacao)
}

func (s *ServidorJogador) registrarVitoria(w http.ResponseWriter, jogador string) {
	s.armazenamento.RegistrarVitoria(jogador)
	w.WriteHeader(http.StatusAccepted)
}

func (s *ServidorJogador) jogo(w http.ResponseWriter, r *http.Request) {
	s.template.Execute(w, nil)
}

func (s *ServidorJogador) webSocket(w http.ResponseWriter, r *http.Request) {
	conexao, _ := atualizadorDeWebsocket.Upgrade(w, r, nil)
	_, msgVencedor, _ := conexao.ReadMessage()
	s.armazenamento.RegistrarVitoria(string(msgVencedor))

}

func NovoServidorJogador(armazenamento ArmazenamentoJogador) (*ServidorJogador, error) {
	s := new(ServidorJogador)

	tmpl, err := template.ParseFiles(caminhoTemplateHTML)

	if err != nil {
		return nil, fmt.Errorf("problema ao abrir %s %v", caminhoTemplateHTML, err)
	}

	s.template = tmpl
	s.armazenamento = armazenamento

	roteador := http.NewServeMux()
	roteador.Handle("/liga", http.HandlerFunc(s.ManipulaLiga))
	roteador.Handle("/jogadores/", http.HandlerFunc(s.ManipulaJogadores))
	roteador.Handle("/jogo", http.HandlerFunc(s.jogo))
	roteador.Handle("/ws", http.HandlerFunc(s.webSocket))

	s.Handler = roteador

	return s, nil
}
