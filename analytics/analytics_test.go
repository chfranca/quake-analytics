package analytics

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractParties(t *testing.T) {
	log := "13:46 Kill: 4 3 7: %s killed %s by %s"
	tests := []struct {
		name           string
		expectedKiller string
		expectedKilled string
		expectedMode   string
	}{
		{
			name:           "should extract parties correctly",
			expectedKiller: "Player 1",
			expectedKilled: "God of War",
			expectedMode:   "MOD_ROCKET_SPLASH",
		},
		{
			name:           "should extract parties correctly when user name common words",
			expectedKiller: "Kill",
			expectedKilled: "killed",
			expectedMode:   "MOD_ROCKET_SPLASH",
		},
		{
			name:           "should extract parties correctly when user is world",
			expectedKiller: "<world>",
			expectedKilled: "Zeh",
			expectedMode:   "MOD_ROCKET_SPLASH",
		},
	}

	for _, s := range tests {
		t.Run(s.name, func(t *testing.T) {
			line := fmt.Sprintf(log, s.expectedKiller, s.expectedKilled, s.expectedMode)
			parties, err := extractParties(line)

			assert.Equal(t, nil, err, "Error not expected")
			assert.Equalf(t, s.expectedKiller, parties[0], "The killer found %s is different of expected %s", parties[0], s.expectedKiller)
			assert.Equalf(t, s.expectedKilled, parties[1], "The killed found %s is different of expected %s", parties[1], s.expectedKilled)
			assert.Equalf(t, s.expectedMode, parties[2], "The mode found %s is different of expected %s", parties[2], s.expectedMode)
		})
	}
}

func TestRegisterKill_addScore(t *testing.T) {

	log := `13:46 Kill: 4 3 7: Dono da Bola killed Oootsimo by MOD_ROCKET_SPLASH`
	scenes := []struct {
		name                string
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
			name:                "should add score of player when is first point of game",
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
			name:                "should add score of player when is first point of player",
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
	}

	for _, v := range scenes {
		t.Run(v.name, func(t *testing.T) {

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
		})
	}
}

func TestRegisterKill_decreaseScore(t *testing.T) {

	log := `13:46 Kill: 4 3 7: <world> killed Oootsimo by MOD_TRIGGER_HURT`
	tests := []struct {
		name                string
		expectedTotalKills  int
		expectedPlayers     []string
		expectedKills       map[string]int
		expectedKillsByMean map[string]int
	}{
		{
			name:                "should decrease score of player when killed by world",
			expectedTotalKills:  9,
			expectedPlayers:     []string{"Oootsimo"},
			expectedKills:       map[string]int{"Oootsimo": 7},
			expectedKillsByMean: map[string]int{"MOD_ROCKET_SPLASH": 8, "MOD_TRIGGER_HURT": 1},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			game := Game{
				TotalKills:    8,
				Players:       []string{"Oootsimo"},
				Kills:         map[string]int{"Oootsimo": 8},
				KillsByMeans:  map[string]int{"MOD_ROCKET_SPLASH": 8},
				PlayerRanking: []string{},
			}

			registerKill(&game, log)

			assert.Equalf(t, v.expectedTotalKills, game.TotalKills, "The total kills: %d is diffent of expected: %d", game.TotalKills, v.expectedTotalKills)
			assert.Truef(t, reflect.DeepEqual(game.Players, v.expectedPlayers), "Players registered after register: %v is diffent of expected: %v", game.Players, v.expectedPlayers)
			assert.Truef(t, reflect.DeepEqual(game.Kills, v.expectedKills), "The kills: %v is diffent of expected: %v", game.Kills, v.expectedKills)
			assert.Truef(t, reflect.DeepEqual(game.KillsByMeans, v.expectedKillsByMean), "The kills by mean: %v is diffent of expected: %v", game.KillsByMeans, v.expectedKillsByMean)
		})
	}
}

func TestRegisterPlayer(t *testing.T) {

	log := `20:38 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\uriel/zael\hmodel\uriel/zael\g_redteam\\g_blueteam\\c1\5\c2\5\hc\100\w\0\l\0\tt\0\tl\0`
	scenes := []struct {
		name            string
		players         []string
		kills           map[string]int
		expectedPlayers []string
		expectedKills   map[string]int
	}{
		{
			name:            "should return one player in the game",
			players:         []string{},
			kills:           map[string]int{},
			expectedPlayers: []string{"Isgalamido"},
			expectedKills:   map[string]int{"Isgalamido": 0},
		},
		{
			name:            "should return two player in the game",
			players:         []string{"Dono da Bola"},
			kills:           map[string]int{"Dono da Bola": 0},
			expectedPlayers: []string{"Dono da Bola", "Isgalamido"},
			expectedKills:   map[string]int{"Dono da Bola": 0, "Isgalamido": 0},
		},
		{
			name:            "should not change the player list",
			players:         []string{"Dono da Bola", "Isgalamido"},
			kills:           map[string]int{"Dono da Bola": 0, "Isgalamido": 0},
			expectedPlayers: []string{"Dono da Bola", "Isgalamido"},
			expectedKills:   map[string]int{"Dono da Bola": 0, "Isgalamido": 0},
		},
	}

	for _, v := range scenes {
		t.Run(v.name, func(t *testing.T) {
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
		})
	}
}

func TestProcessLog(t *testing.T) {

	tests := []struct {
		name           string
		mockReadFile   func(path string) <-chan string
		expectedReport Report
	}{
		{
			name:         "should return two games",
			mockReadFile: mockCompleteLog,
			expectedReport: Report{
				Games: []Game{
					{
						TotalKills:    11,
						Players:       []string{"Isgalamido", "Dono da Bola", "Mocinha"},
						Kills:         map[string]int{"Isgalamido": 0, "Dono da Bola": 0, "Mocinha": 0},
						KillsByMeans:  map[string]int{"MOD_TRIGGER_HURT": 7, "MOD_ROCKET_SPLASH": 3, "MOD_FALLING": 1},
						PlayerRanking: []string{"Isgalamido: 0", "Dono da Bola: 0", "Mocinha: 0"},
					},
					{
						TotalKills:    0,
						Players:       []string{"Isgalamido"},
						Kills:         map[string]int{"Isgalamido": 0},
						KillsByMeans:  map[string]int{},
						PlayerRanking: []string{"Isgalamido: 0"},
					}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := ProcessLog(tt.mockReadFile, "logfile.log")
			assert.Equal(t, tt.expectedReport, *report)
		})
	}
}

func mockCompleteLog(path string) <-chan string {
	stream := make(chan string, 100)

	go func() {
		stream <- `20:37 ------------------------------------------------------------`
		stream <- `20:37 InitGame: \sv_floodProtect\1\sv_maxPing\0\sv_minPing\0\sv_maxRate\10000\sv_minRate\0\sv_hostname\Code Miner Server\g_gametype\0\sv_privateClients\2\sv_maxclients\16\sv_allowDownload\0\bot_minplayers\0\dmflags\0\fraglimit\20\timelimit\15\g_maxGameClients\0\capturelimit\8\version\ioq3 1.36 linux-x86_64 Apr 12 2009\protocol\68\mapname\q3dm17\gamename\baseq3\g_needpass\0`
		stream <- `20:38 ClientConnect: 2`
		stream <- `20:38 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\uriel/zael\hmodel\uriel/zael\g_redteam\\g_blueteam\\c1\5\c2\5\hc\100\w\0\l\0\tt\0\tl\0`
		stream <- `20:38 ClientBegin: 2`
		stream <- `20:40 Item: 2 weapon_rocketlauncher`
		stream <- `20:40 Item: 2 ammo_rockets`
		stream <- `20:42 Item: 2 item_armor_body`
		stream <- `20:54 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT`
		stream <- `20:59 Item: 2 weapon_rocketlauncher`
		stream <- `21:04 Item: 2 ammo_shells`
		stream <- `21:07 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT`
		stream <- `21:10 ClientDisconnect: 2`
		stream <- `21:15 ClientConnect: 2`
		stream <- `21:15 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\uriel/zael\hmodel\uriel/zael\g_redteam\\g_blueteam\\c1\5\c2\5\hc\100\w\0\l\0\tt\0\tl\0`
		stream <- `21:17 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\uriel/zael\hmodel\uriel/zael\g_redteam\\g_blueteam\\c1\5\c2\5\hc\100\w\0\l\0\tt\0\tl\0`
		stream <- `21:17 ClientBegin: 2`
		stream <- `21:18 Item: 2 weapon_rocketlauncher`
		stream <- `21:21 Item: 2 item_armor_body`
		stream <- `21:32 Item: 2 item_health_large`
		stream <- `21:33 Item: 2 weapon_rocketlauncher`
		stream <- `21:34 Item: 2 ammo_rockets`
		stream <- `21:42 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT`
		stream <- `21:49 Item: 2 weapon_rocketlauncher`
		stream <- `21:51 ClientConnect: 3`
		stream <- `21:51 ClientUserinfoChanged: 3 n\Dono da Bola\t\0\model\sarge/krusade\hmodel\sarge/krusade\g_redteam\\g_blueteam\\c1\5\c2\5\hc\95\w\0\l\0\tt\0\tl\0`
		stream <- `21:53 ClientUserinfoChanged: 3 n\Mocinha\t\0\model\sarge\hmodel\sarge\g_redteam\\g_blueteam\\c1\4\c2\5\hc\95\w\0\l\0\tt\0\tl\0`
		stream <- `21:53 ClientBegin: 3`
		stream <- `22:04 Item: 2 weapon_rocketlauncher`
		stream <- `22:04 Item: 2 ammo_rockets`
		stream <- `22:06 Kill: 2 3 7: Isgalamido killed Mocinha by MOD_ROCKET_SPLASH`
		stream <- `22:11 Item: 2 item_quad`
		stream <- `22:11 ClientDisconnect: 3`
		stream <- `22:18 Kill: 2 2 7: Isgalamido killed Isgalamido by MOD_ROCKET_SPLASH`
		stream <- `22:26 Item: 2 weapon_rocketlauncher`
		stream <- `22:27 Item: 2 ammo_rockets`
		stream <- `22:40 Kill: 2 2 7: Isgalamido killed Isgalamido by MOD_ROCKET_SPLASH`
		stream <- `22:43 Item: 2 weapon_rocketlauncher`
		stream <- `22:45 Item: 2 item_armor_body`
		stream <- `23:06 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT`
		stream <- `23:09 Item: 2 weapon_rocketlauncher`
		stream <- `23:10 Item: 2 ammo_rockets`
		stream <- `23:25 Item: 2 item_health_large`
		stream <- `23:30 Item: 2 item_health_large`
		stream <- `23:32 Item: 2 weapon_rocketlauncher`
		stream <- `23:35 Item: 2 item_armor_body`
		stream <- `23:36 Item: 2 ammo_rockets`
		stream <- `23:37 Item: 2 weapon_rocketlauncher`
		stream <- `23:40 Item: 2 item_armor_shard`
		stream <- `23:40 Item: 2 item_armor_shard`
		stream <- `23:40 Item: 2 item_armor_shard`
		stream <- `23:40 Item: 2 item_armor_combat`
		stream <- `23:43 Item: 2 weapon_rocketlauncher`
		stream <- `23:57 Item: 2 weapon_shotgun`
		stream <- `23:58 Item: 2 ammo_shells`
		stream <- `24:13 Item: 2 item_armor_shard`
		stream <- `24:13 Item: 2 item_armor_shard`
		stream <- `24:13 Item: 2 item_armor_shard`
		stream <- `24:13 Item: 2 item_armor_combat`
		stream <- `24:16 Item: 2 item_health_large`
		stream <- `24:18 Item: 2 ammo_rockets`
		stream <- `24:19 Item: 2 weapon_rocketlauncher`
		stream <- `24:22 Item: 2 item_armor_body`
		stream <- `24:24 Item: 2 ammo_rockets`
		stream <- `24:24 Item: 2 weapon_rocketlauncher`
		stream <- `24:36 Item: 2 item_health_large`
		stream <- `24:43 Item: 2 item_health_mega`
		stream <- `25:05 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT`
		stream <- `25:09 Item: 2 weapon_rocketlauncher`
		stream <- `25:09 Item: 2 ammo_rockets`
		stream <- `25:11 Item: 2 item_armor_body`
		stream <- `25:18 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT`
		stream <- `25:21 Item: 2 weapon_rocketlauncher`
		stream <- `25:22 Item: 2 ammo_rockets`
		stream <- `25:34 Item: 2 weapon_rocketlauncher`
		stream <- `25:41 Kill: 1022 2 19: <world> killed Isgalamido by MOD_FALLING`
		stream <- `25:50 Item: 2 item_armor_combat`
		stream <- `25:52 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT`
		stream <- `25:54 Item: 2 ammo_rockets`
		stream <- `25:55 Item: 2 weapon_rocketlauncher`
		stream <- `25:55 Item: 2 weapon_rocketlauncher`
		stream <- `25:59 Item: 2 item_armor_shard`
		stream <- `25:59 Item: 2 item_armor_shard`
		stream <- `26:05 Item: 2 item_armor_shard`
		stream <- `26:05 Item: 2 item_armor_shard`
		stream <- `26:05 Item: 2 item_armor_shard`
		stream <- `26:09 Item: 2 weapon_rocketlauncher`
		stream <- `0:00 ------------------------------------------------------------`
		stream <- `0:00 ------------------------------------------------------------`
		stream <- `0:00 InitGame: \sv_floodProtect\1\sv_maxPing\0\sv_minPing\0\sv_maxRate\10000\sv_minRate\0\sv_hostname\Code Miner Server\g_gametype\0\sv_privateClients\2\sv_maxclients\16\sv_allowDownload\0\dmflags\0\fraglimit\20\timelimit\15\g_maxGameClients\0\capturelimit\8\version\ioq3 1.36 linux-x86_64 Apr 12 2009\protocol\68\mapname\q3dm17\gamename\baseq3\g_needpass\0`
		stream <- `15:00 Exit: Timelimit hit.`
		stream <- `20:34 ClientConnect: 2`
		stream <- `20:34 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\xian/default\hmodel\xian/default\g_redteam\\g_blueteam\\c1\4\c2\5\hc\100\w\0\l\0\tt\0\tl\0`
		stream <- `20:37 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\uriel/zael\hmodel\uriel/zael\g_redteam\\g_blueteam\\c1\5\c2\5\hc\100\w\0\l\0\tt\0\tl\0`
		stream <- `20:37 ClientBegin: 2`
		stream <- `20:37 ShutdownGame:`

		close(stream)
	}()

	return stream
}
