package main

import (
	"awesomeProject/db"
	"awesomeProject/service"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
	"net/http"
)

var dbInstance *gorm.DB
var deckService service.DeckService
var err error

func main() {
	db.InitializeDb()
	db.MigrateTables()
	dbInstance = db.GetDbInstance()
	deckService = service.NewDeckService(dbInstance)
	handleRequests()
}

func handleRequests() {
	router := mux.NewRouter()
	router.HandleFunc("/api/decks", deckService.GetDecks).Methods("GET")
	router.HandleFunc("/api/decks", deckService.CreateCards).Methods("POST")
	router.HandleFunc("/api/decks", deckService.CreateCards).Queries("cards", "{cards}", "shuffled", "{shuffled}").Methods("POST")
	router.HandleFunc("/api/decks/{deck_id}", deckService.GetDeckById).Methods("GET")
	router.HandleFunc("/api/decks/{deck_id}/cards/{card_count}", deckService.DrawCards).Methods("GET")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
	}
}
