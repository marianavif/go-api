package poker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	jogo     Jogo
}

type Jogador struct {
	Nome     string
	Vitorias int
}

type websocketServidorJogador struct {
	*websocket.Conn
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

func (s *ServidorJogador) jogarJogo(w http.ResponseWriter, r *http.Request) {
	s.template.Execute(w, nil)
}

func (s *ServidorJogador) webSocket(w http.ResponseWriter, r *http.Request) {
	ws := NovoWebsocketServidorJogador(w, r)

	mensagemNumeroDeJogadores := ws.EsperarPelaMensagem()
	numeroDeJogadores, _ := strconv.Atoi(string(mensagemNumeroDeJogadores))
	s.jogo.Iniciar(numeroDeJogadores, ws)

	vencedor := ws.EsperarPelaMensagem()
	s.jogo.Terminar(string(vencedor))

}

func (w *websocketServidorJogador) EsperarPelaMensagem() string {
	_, msg, err := w.ReadMessage()
	if err != nil {
		log.Printf("erro ao ler do websocket %v\n", err)
	}

	return string(msg)
}

func (w *websocketServidorJogador) Write(p []byte) (n int, err error) {
	err = w.WriteMessage(1, p)

	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func NovoServidorJogador(armazenamento ArmazenamentoJogador, jogo Jogo) (*ServidorJogador, error) {
	s := new(ServidorJogador)

	tmpl, err := template.ParseFiles(caminhoTemplateHTML)

	if err != nil {
		return nil, fmt.Errorf("problema ao abrir %s %v", caminhoTemplateHTML, err)
	}

	s.jogo = jogo

	s.template = tmpl
	s.armazenamento = armazenamento

	roteador := http.NewServeMux()
	roteador.Handle("/liga", http.HandlerFunc(s.ManipulaLiga))
	roteador.Handle("/jogadores/", http.HandlerFunc(s.ManipulaJogadores))
	roteador.Handle("/jogo", http.HandlerFunc(s.jogarJogo))
	roteador.Handle("/ws", http.HandlerFunc(s.webSocket))

	s.Handler = roteador

	return s, nil
}

func NovoWebsocketServidorJogador(w http.ResponseWriter, r *http.Request) *websocketServidorJogador {
	conexao, err := atualizadorDeWebsocket.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("houve um problema ao atualizar a conexao para WebSockets %v\n", err)
	}

	return &websocketServidorJogador{conexao}
}
