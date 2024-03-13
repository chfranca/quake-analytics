package analytics

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractParties(t *testing.T) {
	t.Parallel()

	log := "13:46 Kill: 4 3 7: %s killed %s by %s"
	tests := []struct {
		expectedKiller string
		expectedKilled string
		expectedMode   string
	}{
		{
			expectedKiller: "Player 1",
			expectedKilled: "God of War",
			expectedMode:   "MOD_ROCKET_SPLASH",
		},
		{
			expectedKiller: "Kill",
			expectedKilled: "killed",
			expectedMode:   "MOD_ROCKET_SPLASH",
		},
		{
			expectedKiller: "<world>",
			expectedKilled: "Zeh",
			expectedMode:   "MOD_ROCKET_SPLASH",
		},
	}

	for _, s := range tests {
		line := fmt.Sprintf(log, s.expectedKiller, s.expectedKilled, s.expectedMode)
		parties, err := extractParties(line)

		assert.Equal(t, nil, err, "Error not expected")
		assert.Equalf(t, s.expectedKiller, parties[0], "The killer found %s is different of expected %s", parties[0], s.expectedKiller)
		assert.Equalf(t, s.expectedKilled, parties[1], "The killed found %s is different of expected %s", parties[1], s.expectedKilled)
		assert.Equalf(t, s.expectedMode, parties[2], "The mode found %s is different of expected %s", parties[2], s.expectedMode)
	}

	s := "12:12:22 Command: Some new log pattern"
	_, err := extractParties(s)

	assert.EqualError(t, err, "impossible extract the parties names from log")
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

		assert.Equalf(t, v.expectedTotalKills, game.TotalKills, "The total kills: %d is diffent of expected: %d", game.TotalKills, v.expectedTotalKills)
		assert.Truef(t, reflect.DeepEqual(game.Players, v.expectedPlayers), "Players registered after register: %v is diffent of expected: %v", game.Players, v.expectedPlayers)
		assert.Truef(t, reflect.DeepEqual(game.Kills, v.expectedKills), "The kills: %v is diffent of expected: %v", game.Kills, v.expectedKills)
		assert.Truef(t, reflect.DeepEqual(game.KillsByMeans, v.expectedKillsByMean), "The kills by mean: %v is diffent of expected: %v", game.KillsByMeans, v.expectedKillsByMean)
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

		assert.Truef(t, reflect.DeepEqual(game.Players, v.expectedPlayers), "Players registered after call: %v is diffent of expected: %v", game.Players, v.expectedPlayers)
		assert.Truef(t, reflect.DeepEqual(game.Kills, v.expectedKills), "The kills: %v is diffent of expected: %v", game.Kills, v.expectedKills)
	}
}

// func TestProcess(t *testing.T) {
// 	mockReader := reader.NewMock()
// 	mockReader.On("Read", mock.Anything).Return(100, nil)
// 	scanSvc := scan.NewInstance(mockReader)

// 	expectedSize := 100
// 	assert.Equal(t, scan.size(), expectedSize)
// }
