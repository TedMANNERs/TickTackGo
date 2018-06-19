package master

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var games = make(map[string]Game)

const ROCK = "Rock"
const PAPER = "Paper"
const SCISSORS = "Scissors"

var availableSymbols = []string{ROCK, PAPER, SCISSORS}

func createGame(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}

		var game Game
		err = json.Unmarshal(body, &game)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}

		game.Result.GameHistory = []GameHistoryEntry{}
		games[game.BoardID] = game
		fmt.Println(game.BoardID + " registered")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(r.Host + "/games/" + game.BoardID))

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method not allowed"))
	}
}

func getGames(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		gamesRef := &games
		encodedGames, err := json.Marshal(gamesRef)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}
		w.Write(encodedGames)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method not allowed"))
	}
}

func getGame(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")
	boardId := p[len(p)-1]
	game, ok := games[boardId]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Not Found, game does not exist"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		gameRef := &game
		encodedGame, err := json.Marshal(gameRef)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}
		w.Write(encodedGame)

	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}

		var slaveSymbol GameSymbol
		err = json.Unmarshal(body, &slaveSymbol)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}

		rand.Seed(time.Now().Unix())
		random := rand.Intn(2)
		masterSymbol := availableSymbols[random]
		historyEntry := GameHistoryEntry{MasterSymbol: masterSymbol, SlaveSymbol: slaveSymbol.Symbol}
		game.Result.GameHistory = append(game.Result.GameHistory, historyEntry)
		switch slaveSymbol.Symbol {
		case ROCK:
			if masterSymbol == PAPER {
				game.Result.MasterScore++
			}
			if masterSymbol == SCISSORS {
				game.Result.SlaveScore++
			}
		case PAPER:
			if masterSymbol == SCISSORS {
				game.Result.MasterScore++
			}
			if masterSymbol == ROCK {
				game.Result.SlaveScore++
			}
		case SCISSORS:
			if masterSymbol == ROCK {
				game.Result.MasterScore++
			}
			if masterSymbol == PAPER {
				game.Result.SlaveScore++
			}
		}
		games[boardId] = game
		resultRef := &game.Result
		encodedResult, err := json.Marshal(resultRef)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}
		w.Write(encodedResult)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method not allowed"))
	}
}

func Start() {
	fmt.Println("Starting master...")
	http.HandleFunc("/registry", createGame)
	http.HandleFunc("/games", getGames)
	http.HandleFunc("/games/", getGame)
	http.ListenAndServe("localhost:8080", nil)
}
