package master

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var games = make(map[string]Game)
var azureUrl = "https://techweek2018functionapp.azurewebsites.net/api/scores"

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

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
		w.Header().Set("content-location", "http://"+r.Host+"/games/"+game.BoardID)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Game created"))

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method not allowed"))
	}
}

func getGames(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
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
		w.Header().Set("content-type", "application/json")
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

		w.Header().Set("content-type", "application/json")
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

		game.Result = GetUpdatedResult(game.Result, slaveSymbol)
		games[boardId] = game
		resultRef := &game.Result
		encodedResult, err := json.Marshal(resultRef)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}

		w.Header().Set("content-type", "application/json")
		w.Write(encodedResult)

		go updateAzure(game)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method not allowed"))
	}
}

func updateAzure(game Game) {
	resultJson := []byte(fmt.Sprintf(`{"masterScore":%d, "slaveScore":%d}`, game.Result.MasterScore, game.Result.SlaveScore))
	r, err := http.NewRequest(http.MethodPut, azureUrl, bytes.NewBuffer(resultJson))
	if err != nil {
		fmt.Println(err)
		return
	}
	q := r.URL.Query()
	q.Add("id", game.BoardID)
	r.URL.RawQuery = q.Encode()
	r.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
}

func Start() {
	fmt.Println("Starting master...")
	http.HandleFunc("/registry", createGame)
	http.HandleFunc("/games", getGames)
	http.HandleFunc("/games/", getGame)
	http.ListenAndServe("localhost:8080", nil)
}
