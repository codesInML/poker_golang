package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)
	player := "Pepper"

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
			{"Ifeoluwa", 20},
		}

		assertLeague(t, got, want)

		// github personal token
		// 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
	})
}
