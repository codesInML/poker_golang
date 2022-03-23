package main

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func createTempFile(t testing.TB, initialData string) (io.ReadWriteSeeker, func()) {
	t.Helper()

	tempfile, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tempfile.Write([]byte(initialData))

	removeFile := func() {
		tempfile.Close()
		os.Remove(tempfile.Name())
	}

	return tempfile, removeFile
}

func TestFileSystemStore(t *testing.T) {

	t.Run("league from a reader", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Ifeoluwa", "Wins": 20},
			{"Name": "Ifeanyi", "Wins": 5}
		]`)

		defer cleanDatabase()

		store := FileSystemPlayerStore{database}

		got := store.GetLeague()
		want := []Player{
			{"Ifeoluwa", 20},
			{"Ifeanyi", 5},
		}

		assertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		assertLeague(t, got, want)
	})

	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Ifeoluwa", "Wins": 20},
			{"Name": "Ifeanyi", "Wins": 5}
		]`)

		defer cleanDatabase()

		store := FileSystemPlayerStore{database}

		got := store.GetPlayerScore("Ifeoluwa")
		want := 20
		assertScoreEquals(t, got, want)
	})

	t.Run("store wins on existing players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Ifeoluwa", "Wins": 20},
			{"Name": "Ifeanyi", "Wins": 5}
		]`)

		defer cleanDatabase()
		store := FileSystemPlayerStore{database}
		store.RecordWin("Ifeoluwa")

		got := store.GetPlayerScore("Ifeoluwa")
		want := 21

		assertScoreEquals(t, got, want)
	})
}

func assertScoreEquals(t *testing.T, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
