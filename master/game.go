package master

import (
	"math/rand"
	"time"
)

const ROCK = "Rock"
const PAPER = "Paper"
const SCISSORS = "Scissors"

var availableSymbols = []string{ROCK, PAPER, SCISSORS}

type Game struct {
	BoardID     string
	MasterScore int
	SlaveScore  int
	GameHistory []GameHistoryEntry
}

type PostResult struct {
	MasterScore  int
	SlaveScore   int
	MasterSymbol string
	SlaveSymbol  string
}

type GameHistoryEntry struct {
	MasterSymbol string
	SlaveSymbol  string
}

type GameSymbol struct {
	Symbol string
}

func GetUpdatedResult(game Game, slaveSymbol GameSymbol) Game {
	rand.Seed(time.Now().Unix())
	random := rand.Intn(3)
	masterSymbol := availableSymbols[random]
	historyEntry := GameHistoryEntry{MasterSymbol: masterSymbol, SlaveSymbol: slaveSymbol.Symbol}
	game.GameHistory = append(game.GameHistory, historyEntry)
	switch slaveSymbol.Symbol {
	case ROCK:
		if masterSymbol == PAPER {
			game.MasterScore++
		}
		if masterSymbol == SCISSORS {
			game.SlaveScore++
		}
	case PAPER:
		if masterSymbol == SCISSORS {
			game.MasterScore++
		}
		if masterSymbol == ROCK {
			game.SlaveScore++
		}
	case SCISSORS:
		if masterSymbol == ROCK {
			game.MasterScore++
		}
		if masterSymbol == PAPER {
			game.SlaveScore++
		}
	}
	return game
}
