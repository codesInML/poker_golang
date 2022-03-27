package poker

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
	return s.league
}

func TestGETPlayers(t *testing.T) {
	store := &StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		[]string{},
		[]Player{},
	}

	server := NewPlayerServer(store)
	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseBody(response.Body.String(), "20", t)
		assertResponseStatusCode(response.Code, http.StatusOK, t)
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseBody(response.Body.String(), "10", t)
		assertResponseStatusCode(response.Code, http.StatusOK, t)
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatusCode(response.Code, http.StatusNotFound, t)
	})
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		[]string{},
		[]Player{},
	}

	server := NewPlayerServer(&store)

	t.Run("it records win on POST", func(t *testing.T) {
		player := "Pepper"
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertResponseStatusCode(response.Code, http.StatusAccepted, t)

		if len(store.winCalls) != 1 {
			t.Fatalf("got %d calls to RecordWin, want %d", len(store.winCalls), 1)
		}

		if store.winCalls[0] != player {
			t.Errorf("did not store the correct winner, got %q, want %q", store.winCalls[0], player)
		}
	})
}

func TestLeague(t *testing.T) {
	wantedLeague := []Player{
		{"Ifeoluwa", 20},
		{"Ifeany", 19},
		{"Samuel", 18},
	}

	store := StubPlayerStore{nil, nil, wantedLeague}
	server := NewPlayerServer(&store)

	t.Run("it returns the league table as JSON", func(t *testing.T) {
		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)
		assertResponseStatusCode(response.Code, http.StatusOK, t)

		assertLeague(t, got, wantedLeague)
		assertContentType(t, response, jsonContentType)
	})
}

func newLeagueRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return request
}

func getLeagueFromResponse(t testing.TB, body io.Reader) (league []Player) {
	t.Helper()

	err := json.NewDecoder(body).Decode(&league)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}
	return
}

func assertLeague(t testing.TB, got, want []Player) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertContentType(t testing.TB, got *httptest.ResponseRecorder, want string) {
	t.Helper()

	if got.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have the content type of %s, got %v", want, got.Result().Header.Get("content-type"))
	}
}

func newGetScoreRequest(player string) *http.Request {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", player), nil)

	if err != nil {
		log.Fatal(err)
	}

	return request
}

func newPostWinRequest(player string) *http.Request {
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", player), nil)

	if err != nil {
		log.Fatal(err)
	}

	return request
}

func assertResponseBody(got, want string, t testing.TB) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertResponseStatusCode(got, want int, t testing.TB) {
	t.Helper()
	if got != want {
		t.Errorf("did not get the correct status code, got %d, want %d", got, want)
	}
}
