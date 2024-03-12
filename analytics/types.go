package analytics

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
)

const (
	initGame      string = "InitGame"
	kill          string = "Kill"
	endGame       string = "ShutdownGame"
	userConnected string = "ClientUserinfoChanged"
)

const (
	worldUser string = "<world>"
)

type Report struct {
	Games []Game
}

type Game struct {
	TotalKills   int            `json:"total_kills"`
	Players      []string       `json:"players"`
	Kills        map[string]int `json:"kills"`
	KillsByMeans map[string]int `json:"kills_by_means"`
}

func (g Game) String() string {
	json, err := json.MarshalIndent(g, "", "   ")

	if err != nil {
		log.Fatal(err)
	}

	return string(json)
}

func (g *Game) PlayerRanking() string {

	type pair struct {
		Key   string
		Value int
	}

	var ranking []pair
	for k, v := range g.Kills {
		ranking = append(ranking, pair{k, v})
	}

	sort.Slice(ranking, func(i, j int) bool {
		return ranking[i].Value > ranking[j].Value
	})

	var rankingString []string
	for _, pair := range ranking {
		rankingString = append(rankingString, fmt.Sprintf("%s:%d", pair.Key, pair.Value))
	}
	return strings.Join(rankingString, ",\n")
}
