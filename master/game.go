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
	BoardID string
	Result  GameResult
}

type GameResult struct {
	MasterScore int
	SlaveScore  int
	GameHistory []GameHistoryEntry
}

type GameHistoryEntry struct {
	MasterSymbol string
	SlaveSymbol  string
}

type GameSymbol struct {
	Symbol string
}

func GetUpdatedResult(result GameResult, slaveSymbol GameSymbol) GameResult {
	rand.Seed(time.Now().Unix())
	random := rand.Intn(3)
	masterSymbol := availableSymbols[random]
	historyEntry := GameHistoryEntry{MasterSymbol: masterSymbol, SlaveSymbol: slaveSymbol.Symbol}
	result.GameHistory = append(result.GameHistory, historyEntry)
	switch slaveSymbol.Symbol {
	case ROCK:
		if masterSymbol == PAPER {
			result.MasterScore++
		}
		if masterSymbol == SCISSORS {
			result.SlaveScore++
		}
	case PAPER:
		if masterSymbol == SCISSORS {
			result.MasterScore++
		}
		if masterSymbol == ROCK {
			result.SlaveScore++
		}
	case SCISSORS:
		if masterSymbol == ROCK {
			result.MasterScore++
		}
		if masterSymbol == PAPER {
			result.SlaveScore++
		}
	}
	return result
}
