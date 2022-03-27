package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
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

		store, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)

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

		store, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)

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
		store, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)
		store.RecordWin("Ifeoluwa")

		got := store.GetPlayerScore("Ifeoluwa")
		want := 21

		assertScoreEquals(t, got, want)
	})

	t.Run("store wins of new players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Ifeoluwa", "Wins": 20},
			{"Name": "Ifeanyi", "Wins": 5}
		]`)

		defer cleanDatabase()
		store, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)
		store.RecordWin("Oluwole")

		got := store.GetPlayerScore("Oluwole")
		want := 1
		assertScoreEquals(t, got, want)
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")

		defer cleanDatabase()
		_, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)
	})

	t.Run("league sorted", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Ifeoluwa", "Wins": 20},
			{"Name": "Ifeanyi", "Wins": 5},
			{"Name": "Oluwole", "Wins": 25}]`)

		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)

		got := store.GetLeague()

		want := League{
			{"Oluwole", 25},
			{"Ifeoluwa", 20},
			{"Ifeanyi", 5},
		}
		assertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		assertLeague(t, got, want)
	})

}

func TestTape_Write(t *testing.T) {
	file, clean := createTempFile(t, "12345")
	defer clean()

	tape := &tape{file}

	tape.Write([]byte("abc"))

	file.Seek(0, 0)
	newFileContents, _ := ioutil.ReadAll(file)

	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertScoreEquals(t *testing.T, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func assertNoError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("got an error but did not expect one, %v", err)
	}
}
