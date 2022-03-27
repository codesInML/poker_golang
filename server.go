package poker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() League
}

type PlayerServer struct {
	store PlayerStore
	http.Handler
}

type Player struct {
	Name string
	Wins int
}

const jsonContentType = "application/json"

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)

	p.store = store
	router := http.NewServeMux()

	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))

	p.Handler = router
	return p
}

func (p *PlayerServer) leagueHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("content-type", jsonContentType)
	json.NewEncoder(rw).Encode(p.store.GetLeague())
	rw.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) playersHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p.ShowScore(rw, r)
	case http.MethodPost:
		p.ProcessWin(rw, r)
	}
}

func (p *PlayerServer) ShowScore(rw http.ResponseWriter, r *http.Request) {
	player := trimPlayerNamePrefix(r)

	score := p.store.GetPlayerScore(player)

	if score == 0 {
		rw.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(rw, score)
}

func (p *PlayerServer) ProcessWin(rw http.ResponseWriter, r *http.Request) {
	player := trimPlayerNamePrefix(r)
	p.store.RecordWin(player)
	rw.WriteHeader(http.StatusAccepted)
}

func trimPlayerNamePrefix(r *http.Request) string {
	return strings.TrimPrefix(r.URL.Path, "/players/")
}
