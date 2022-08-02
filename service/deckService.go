package service

import (
	"awesomeProject/domain"
	"awesomeProject/dto"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type DeckService interface {
	GetDecks(w http.ResponseWriter, r *http.Request)
	CreateCards(w http.ResponseWriter, r *http.Request)
	GetDeckById(w http.ResponseWriter, r *http.Request)
	DrawCards(w http.ResponseWriter, r *http.Request)
}

type deckService struct {
	dbInstance *gorm.DB
}

func NewDeckService(dbInstance *gorm.DB) DeckService {
	return &deckService{dbInstance: dbInstance}
}

func (d *deckService) GetDecks(w http.ResponseWriter, r *http.Request) {
	var decks []domain.Deck
	var decksResponse []dto.Deck
	d.dbInstance.Preload("Cards", "revealed IS FALSE").
		Find(&decks)
	for _, deck := range decks {
		if len(deck.Cards) > 0 {
			decksResponse = append(decksResponse, mapToDeckResponseModel(deck))
		}
	}

	jsonResponse, jsonError := json.Marshal(decksResponse)
	if jsonError != nil {
		fmt.Println("Unable to encode JSON")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (d *deckService) CreateCards(w http.ResponseWriter, r *http.Request) {
	var cards []domain.Card
	specificCards := r.FormValue("cards")
	if specificCards != "" {
		cardsSlice := strings.Split(specificCards, ",")
		cards = createSpecificCardsObject(cardsSlice)
	} else {
		cards = createCardsObject()
	}

	var deck = domain.Deck{
		Shuffled: false,
		Cards:    cards,
	}
	d.dbInstance.Save(&deck)

	deckResponse := mapToDeckResponseModel(deck)
	jsonResponse, _ := json.Marshal(deckResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (d *deckService) GetDeckById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deckId := vars["deck_id"]
	var deck, success = d.getDeckFromDb(deckId)
	if !success {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		return
	}
	deckResponse := mapToDeckResponseModel(deck)
	jsonResponse, _ := json.Marshal(deckResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (d *deckService) DrawCards(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deckId := vars["deck_id"]
	cardsCount, err := strconv.Atoi(vars["card_count"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		return
	}

	var deck, success = d.getDeckFromDb(deckId)
	if !success {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		return
	}
	var cards []domain.Card
	for i := 0; i < cardsCount; i++ {
		cards = append(cards, deck.Cards[i])
		deck.Cards[i].Revealed = true
		d.dbInstance.Save(&deck.Cards[i])
	}
	cardsResponse := mapToCardsResponseModel(cards)
	jsonResponse, _ := json.Marshal(cardsResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func createCardsObject() []domain.Card {
	var cards []domain.Card
	colors := []string{"Spades", "Clubs", "Diamonds", "Hearts"}
	faces := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "Jack", "Queen", "King"}

	for _, color := range colors {
		for _, face := range faces {
			value := face
			suit := color
			code := getCode(face, color)
			var card = domain.Card{
				Value: value,
				Suit:  suit,
				Code:  code,
			}
			cards = append(cards, card)
		}
	}

	println(len(cards))
	for _, card := range cards {
		print(card.Value + " ")
		print(card.Suit + " ")
		println(card.Code)
	}
	return cards
}

func createSpecificCardsObject(codes []string) []domain.Card {
	var cards []domain.Card

	for _, code := range codes {
		validateCode(code)
		value := getValue(code)
		suit := getSuit(code)
		code := code
		var card = domain.Card{
			Value: value,
			Suit:  suit,
			Code:  code,
		}
		cards = append(cards, card)
	}

	println(len(cards))
	for _, card := range cards {
		print(card.Value + " ")
		print(card.Suit + " ")
		println(card.Code)
	}
	return cards
}

func validateCode(code string) {
	faces := []string{"A", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	colors := []string{"S", "C", "D", "H"}
	var firstLetter string
	var secondLetter string
	if len(code) == 2 {
		firstLetter = code[0:1]
		secondLetter = code[1:2]
	}
	if len(code) == 3 {
		firstLetter = code[0:2]
		secondLetter = code[2:3]
	}
	faceCorrect := contains(firstLetter, faces)
	colorCorrect := contains(secondLetter, colors)
	if !(faceCorrect && colorCorrect) {
		panic(errors.New("malformed code"))
	}
}

func getValue(code string) string {
	if len(code) == 2 {
		firstLetter := code[0:1]
		switch firstLetter {
		case "K":
			return "King"
		case "Q":
			return "Queen"
		case "J":
			return "Jack"
		case "A":
			return "1"
		default:
			return firstLetter
		}
	} else {
		return code[0:2]
	}
}

func getSuit(code string) string {
	var lastLetter string
	if len(code) == 2 {
		lastLetter = code[1:2]
	} else if len(code) == 3 {
		lastLetter = code[2:3]
	} else {
		panic(errors.New("malformed code"))
	}
	switch lastLetter {
	case "H":
		return "Hearts"
	case "S":
		return "Spades"
	case "D":
		return "Diamonds"
	case "C":
		return "Clubs"
	}
	panic(errors.New("malformed code"))
}

func getCode(face string, color string) string {
	if face == "10" {
		return "10" + color[0:1]
	} else {
		return face[0:1] + color[0:1]
	}
}

func contains(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (d *deckService) getDeckFromDb(deckId string) (domain.Deck, bool) {
	var deck domain.Deck
	d.dbInstance.Preload("Cards", "revealed IS FALSE").First(&deck, "id = ?", deckId)
	if deck.ID.String() != deckId || len(deck.Cards) == 0 {
		return deck, false
	}
	return deck, true
}

func mapToDeckResponseModel(deck domain.Deck) dto.Deck {
	var cards = mapToCardsResponseModel(deck.Cards)
	return dto.Deck{Id: deck.ID.String(), Shuffled: deck.Shuffled, Remaining: len(deck.Cards), Cards: cards}
}

func mapToCardsResponseModel(cards []domain.Card) []dto.Card {
	var cardsResponse []dto.Card
	for _, card := range cards {
		cardsResponse = append(cardsResponse, mapToCardResponseModel(card))
	}
	return cardsResponse
}

func mapToCardResponseModel(card domain.Card) dto.Card {
	return dto.Card{Value: card.Value, Suit: card.Suit, Code: card.Code}
}
