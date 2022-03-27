package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const dbFileName = "game.db.json"
const PORT = "3030"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s, %v", dbFileName, err)
	}

	store, err := NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatalf("problem creating file system store, %v", err)
	}
	server := NewPlayerServer(store)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), server); err != nil {
		log.Fatalf("could not listen on port %s, %v", PORT, err)
	}
}
