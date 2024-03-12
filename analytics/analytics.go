package analytics

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func ProcessLog(path string) *Report {

	stream := readFile(path)
	results := process(stream)

	report := Report{Games: []Game{}}

	for r := range results {
		report.Games = append(report.Games, r)
	}

	return &report
}

func readFile(path string) <-chan string {

	stream := make(chan string, 100)

	go func() {
		file, err := os.Open(path)

		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			stream <- scanner.Text()
		}
		close(stream)
	}()

	return stream
}

func process(stream <-chan string) <-chan Game {

	var game Game
	results := make(chan Game, 100)

	go func() {
		for line := range stream {
			switch {
			case strings.Contains(line, initGame):
				if len(game.Players) != 0 {
					results <- game
				}
				startGame(&game)
			case strings.Contains(line, endGame):
				results <- game
				game = Game{}
			case strings.Contains(line, userConnected):
				registerPlayer(&game, line)
			case strings.Contains(line, kill):
				registerKill(&game, line)
			default:
			}
		}
		close(results)
	}()

	return results

}

func startGame(game *Game) {
	*game = Game{
		TotalKills:   0,
		Players:      []string{},
		Kills:        map[string]int{},
		KillsByMeans: map[string]int{},
	}
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
