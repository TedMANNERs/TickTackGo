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

		game.GameHistory = []GameHistoryEntry{}
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

		game = GetUpdatedResult(game, slaveSymbol)
		games[boardId] = game

		postResult := createPostResult(game)
		resultRef := &postResult
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

func createPostResult(game Game) PostResult {
	lastHistoryEntry := game.GameHistory[len(game.GameHistory)-1]
	postResult := PostResult{
		MasterScore:  game.MasterScore,
		SlaveScore:   game.SlaveScore,
		MasterSymbol: lastHistoryEntry.MasterSymbol,
		SlaveSymbol:  lastHistoryEntry.SlaveSymbol}
	return postResult
}

func updateAzure(game Game) {
	resultString := fmt.Sprintf(`{"masterScore":%d, "slaveScore":%d}`, game.MasterScore, game.SlaveScore)
	fmt.Println("Update Azure: " + resultString)
	resultJson := []byte(resultString)
	req, err := http.NewRequest(http.MethodPut, azureUrl, bytes.NewBuffer(resultJson))
	if err != nil {
		fmt.Println(err)
		return
	}
	query := req.URL.Query()
	query.Add("id", game.BoardID)
	req.URL.RawQuery = query.Encode()
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(body))
	defer resp.Body.Close()
}

func Start() {
	fmt.Println("Starting master...")
	http.HandleFunc("/registry", createGame)
	http.HandleFunc("/games", getGames)
	http.HandleFunc("/games/", getGame)
	http.ListenAndServe(":8080", nil)
}
