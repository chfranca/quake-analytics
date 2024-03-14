package analytics

import (
	"bufio"
	"errors"
	"log"
	"os"
	"regexp"
	"strings"
)

type Reader func(path string) <-chan string

var regex = regexp.MustCompile(`(?m)(.*) killed (.*) by (.*)`)

func ProcessLog(reader Reader, path string) *Report {

	stream := reader(path)
	results := process(stream)

	report := Report{Games: []Game{}}

	for g := range results {
		g.makeRanking()
		report.Games = append(report.Games, g)
	}

	return &report
}

func ReadFile(path string) <-chan string {

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
			case strings.Contains(line, string(initGame)):
				if len(game.Players) != 0 {
					results <- game
				}
				startGame(&game)
			case strings.Contains(line, string(endGame)):
				results <- game
				game = Game{}
			case strings.Contains(line, string(userConnected)):
				registerPlayer(&game, line)
			case strings.Contains(line, string(kill)):
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
	parties, err := extractParties(line)
	if err != nil {
		log.Println("Log not corresponding with pattern", line)
		return
	}

	killer, killed, mode := parties[0], parties[1], parties[2]
	game.TotalKills++

	// just ensure that user self killed not add score
	if killer != string(worldName) {
		game.Kills[killer]++
	} else if game.Kills[killed] > 0 {
		game.Kills[killed]--
	}

	game.KillsByMeans[mode]++
}

func extractParties(log string) ([]string, error) {
	line := strings.Split(log, ":")
	log = strings.TrimSpace(line[len(line)-1])

	parties := regex.FindStringSubmatch(log) // the first element is the full string

	if len(parties) != 4 {
		return nil, errors.New("impossible extract the parties names from log")
	}

	return parties[1:], nil
}
