package quakeanalytics

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func ProcessLog(path string, report *Report) {

	stream := make(chan string, 100)

	go process(stream, report)
	readFile(path, stream)
}

func readFile(path string, stream chan<- string) {
	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		stream <- scanner.Text()
	}
}

func startGame() *Game {
	return &Game{
		TotalKills:   0,
		Players:      []string{},
		Kills:        map[string]int{},
		KillsByMeans: map[string]int{},
	}
}

func shutdownGame(game *Game, report *Report) {
	report.Games = append(report.Games, *game)
}

func registerPlayer(game *Game, line string) {
	player := strings.Split(line, "\\")[1]
	if _, exists := game.Kills[player]; !exists {
		game.Players = append(game.Players, player)
		game.Kills[player] = 0
	}
}

func registerKill(game *Game, line string) {
	log := strings.Split(line, ":")
	killer, killed, mode := extractParties(log[len(log)-1])
	game.TotalKills++

	if killer != worldUser {
		game.Kills[killer]++
	} else if game.Kills[killed] > 0 {
		game.Kills[killed]--
	}

	game.KillsByMeans[mode]++
}

func extractParties(log string) (killer string, killed string, mode string) {
	parties := strings.Split(log, "killed")
	killer = strings.TrimSpace(parties[0])
	killed = strings.TrimSpace(strings.Split(parties[1], "by")[0])
	mode = strings.TrimSpace(strings.Split(parties[1], "by")[1])
	return
}

func process(stream <-chan string, report *Report) {

	var game *Game
	for line := range stream {
		switch {
		case strings.Contains(line, initGame):
			if game != nil && len(game.Players) > 0 {
				shutdownGame(game, report)
			}
			game = startGame()
		case strings.Contains(line, userConnected):
			registerPlayer(game, line)
		case strings.Contains(line, kill):
			registerKill(game, line)
		default:
		}
	}
}
