package analytics

import (
	"fmt"
	"reflect"
	"testing"
)

func TestExtractParties(t *testing.T) {
	t.Parallel()

	log := "13:46 Kill: 4 3 7: %s killed %s by %s"
	scenes := []struct {
		expectedKiller string
		expectedKilled string
		expectedMode   string
	}{
		{expectedKiller: "Player 1", expectedKilled: "God of War", expectedMode: "MOD_ROCKET_SPLASH"},
		{expectedKiller: "Kill", expectedKilled: "killed", expectedMode: "MOD_ROCKET_SPLASH"},
		{expectedKiller: "<world>", expectedKilled: "Zeh", expectedMode: "MOD_ROCKET_SPLASH"},
	}

	for _, s := range scenes {
		line := fmt.Sprintf(log, s.expectedKiller, s.expectedKilled, s.expectedMode)
		parties, err := extractParties(line)

		if err != nil {
			t.Error(err)
		}

		if parties[0] != s.expectedKiller {
			t.Errorf("The killer found %s is different of expected %s", parties[0], s.expectedKiller)
		}

		if parties[1] != s.expectedKilled {
			t.Errorf("The killed found %s is different of expected %s", parties[1], s.expectedKilled)
		}

		if parties[2] != s.expectedMode {
			t.Errorf("The mode found %s is different of expected %s", parties[2], s.expectedMode)
		}
	}

	// error case
	expectedMessage := "impossible extract the parties names from log"
	s := "12:12:22 Command: Some new log pattern"
	_, err := extractParties(s)

	if err == nil {
		t.Error("should be dispatch a error when log pattern changes")
	}

	if err.Error() != expectedMessage {
		t.Error("the error received is different of expected")
	}
}

func TestRegisterKill(t *testing.T) {
	t.Parallel()

	log := `13:46 Kill: 4 3 7: Dono da Bola killed Oootsimo by MOD_ROCKET_SPLASH`
	scenes := []struct {
		totalKills          int
		players             []string
		kills               map[string]int
		killsByMean         map[string]int
		log                 string
		expectedTotalKills  int
		expectedPlayers     []string
		expectedKills       map[string]int
		expectedKillsByMean map[string]int
	}{
		{
			// test first point on the game
			totalKills:          0,
			players:             []string{"Dono da Bola", "Oootsimo"},
			kills:               map[string]int{},
			killsByMean:         map[string]int{},
			log:                 log,
			expectedTotalKills:  1,
			expectedPlayers:     []string{"Dono da Bola", "Oootsimo"},
			expectedKills:       map[string]int{"Dono da Bola": 1},
			expectedKillsByMean: map[string]int{"MOD_ROCKET_SPLASH": 1},
		},
		{
			// test first point of the player
			totalKills:          8,
			players:             []string{"Dono da Bola", "Oootsimo"},
			kills:               map[string]int{"Oootsimo": 8},
			killsByMean:         map[string]int{"MOD_ROCKET_SPLASH": 8},
			log:                 log,
			expectedTotalKills:  9,
			expectedPlayers:     []string{"Dono da Bola", "Oootsimo"},
			expectedKills:       map[string]int{"Oootsimo": 8, "Dono da Bola": 1},
			expectedKillsByMean: map[string]int{"MOD_ROCKET_SPLASH": 9},
		},
		{
			// test when world kill a player
			totalKills:          8,
			players:             []string{"Dono da Bola", "Oootsimo"},
			kills:               map[string]int{"Oootsimo": 8},
			killsByMean:         map[string]int{"MOD_ROCKET_SPLASH": 8},
			log:                 "13:46 Kill: 4 3 7: <world> killed Oootsimo by MOD_TRIGGER_HURT",
			expectedTotalKills:  9,
			expectedPlayers:     []string{"Dono da Bola", "Oootsimo"},
			expectedKills:       map[string]int{"Oootsimo": 7},
			expectedKillsByMean: map[string]int{"MOD_ROCKET_SPLASH": 8, "MOD_TRIGGER_HURT": 1},
		},

		{
			// test when world kill a player
			totalKills:          8,
			players:             []string{"Dono da Bola", "Oootsimo"},
			kills:               map[string]int{"Oootsimo": 8},
			killsByMean:         map[string]int{"MOD_ROCKET_SPLASH": 8},
			log:                 "13:46 Command: some random log",
			expectedTotalKills:  8,
			expectedPlayers:     []string{"Dono da Bola", "Oootsimo"},
			expectedKills:       map[string]int{"Oootsimo": 8},
			expectedKillsByMean: map[string]int{"MOD_ROCKET_SPLASH": 8},
		},
	}

	for _, v := range scenes {
		game := Game{
			TotalKills:    v.totalKills,
			Players:       v.expectedPlayers,
			Kills:         v.kills,
			KillsByMeans:  v.killsByMean,
			PlayerRanking: []string{},
		}

		registerKill(&game, v.log)

		if game.TotalKills != v.expectedTotalKills {
			t.Errorf("The total kills: %d is diffent of expected: %d", game.TotalKills, v.expectedTotalKills)
		}

		if !reflect.DeepEqual(game.Players, v.expectedPlayers) {
			t.Errorf("Players registered after register: %v is diffent of expected: %v", game.Players, v.expectedPlayers)
		}

		if !reflect.DeepEqual(game.Kills, v.expectedKills) {
			t.Errorf("The kills: %v is diffent of expected: %v", game.Kills, v.expectedKills)
		}

		if !reflect.DeepEqual(game.KillsByMeans, v.expectedKillsByMean) {
			t.Errorf("The kills by mean: %v is diffent of expected: %v", game.KillsByMeans, v.expectedKillsByMean)
		}
	}
}

func TestRegisterPlayer(t *testing.T) {
	t.Parallel()

	log := `20:38 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\uriel/zael\hmodel\uriel/zael\g_redteam\\g_blueteam\\c1\5\c2\5\hc\100\w\0\l\0\tt\0\tl\0`
	scenes := []struct {
		players         []string
		kills           map[string]int
		expectedPlayers []string
		expectedKills   map[string]int
	}{
		{
			// test first player connected
			players:         []string{},
			kills:           map[string]int{},
			expectedPlayers: []string{"Isgalamido"},
			expectedKills:   map[string]int{"Isgalamido": 0},
		},
		{
			// test second player connected
			players:         []string{"Dono da Bola"},
			kills:           map[string]int{"Dono da Bola": 0},
			expectedPlayers: []string{"Dono da Bola", "Isgalamido"},
			expectedKills:   map[string]int{"Dono da Bola": 0, "Isgalamido": 0},
		},
		{
			// test when player already registered on the session
			players:         []string{"Dono da Bola", "Isgalamido"},
			kills:           map[string]int{"Dono da Bola": 0, "Isgalamido": 0},
			expectedPlayers: []string{"Dono da Bola", "Isgalamido"},
			expectedKills:   map[string]int{"Dono da Bola": 0, "Isgalamido": 0},
		},
	}

	for _, v := range scenes {
		game := Game{
			TotalKills:    0,
			Players:       v.players,
			Kills:         v.kills,
			KillsByMeans:  map[string]int{},
			PlayerRanking: []string{},
		}
		registerPlayer(&game, log)

		if !reflect.DeepEqual(game.Players, v.expectedPlayers) {
			t.Errorf("Players registered after call: %v is diffent of expected: %v", game.Players, v.expectedPlayers)
		}

		if !reflect.DeepEqual(game.Kills, v.expectedKills) {
			t.Errorf("The kills: %v is diffent of expected: %v", game.Kills, v.expectedKills)
		}
	}
}
