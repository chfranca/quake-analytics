package analytics

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
)

type Command string

const (
	initGame      Command = "InitGame"
	kill          Command = "Kill"
	endGame       Command = "ShutdownGame"
	userConnected Command = "ClientUserinfoChanged"
)

func (c Command) String() string {
	return string(c)
}

type World string

const (
	worldName World = "<world>"
)

type Report struct {
	Games []Game
}

type Game struct {
	TotalKills    int            `json:"total_kills"`
	Players       []string       `json:"players"`
	Kills         map[string]int `json:"kills"`
	KillsByMeans  map[string]int `json:"kills_by_means"`
	PlayerRanking []string       `json:"ranking"`
}

func (g *Game) String() string {
	json, err := json.MarshalIndent(g, "", "   ")

	if err != nil {
		log.Fatal(err)
	}

	return string(json)
}

func (g *Game) makeRanking() {

	type pair struct {
		Key   string
		Value int
	}

	var players []pair
	for k, v := range g.Kills {
		players = append(players, pair{Key: k, Value: v})
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].Value > players[j].Value
	})

	g.PlayerRanking = []string{}
	for _, v := range players {
		g.PlayerRanking = append(g.PlayerRanking, fmt.Sprintf("%s: %d", v.Key, v.Value))
	}
}
