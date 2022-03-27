package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := createTempFile(t, "")
	defer cleanDatabase()
	store := NewFileSystemPlayerStore(database)
	server := NewPlayerServer(store)
	player := "Ifeoluwa"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		assertResponseStatusCode(response.Code, http.StatusOK, t)
		assertResponseBody(response.Body.String(), "4", t)
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		assertResponseStatusCode(response.Code, http.StatusOK, t)

		got := getLeagueFromResponse(t, response.Body)
		want := []Player{
			{"Ifeoluwa", 4},
		}

		assertLeague(t, got, want)

	})
}
