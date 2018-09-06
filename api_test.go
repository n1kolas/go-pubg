package pubg

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type Conf struct {
	Key string `json:"key"`
}

var client *Client
var testPlayers = []string{"chocoTaco"}
var testPlayerIds = []string{"account.eca3ff5448b6418092f0e3e7e4693066"}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestGetStatus(t *testing.T) {
	_, err := client.GetStatus()
	if err != nil {
		t.Errorf("Error getting status from API: %s", err.Error())
	}
}
func TestGetPlayer(t *testing.T) {
	player, err := client.GetPlayer(testPlayerIds[0], PCNorthAmerica)
	if err != nil {
		t.Errorf("Error getting player from API: %s", err.Error())
	}

	if player.Data.ID != testPlayerIds[0] {
		t.Errorf("expected a player name of %s but received %s", testPlayerIds[0], player.Data.ID)
	}
}
func TestGetPlayersByName(t *testing.T) {
	options := GetPlayersRequestOptions{
		PlayerNamesFilter: testPlayers,
	}

	players, err := client.GetPlayers(options, PCNorthAmerica)
	if err != nil {
		t.Errorf("Error getting players from API: %s", err.Error())
	}

	if len(testPlayers) != len(players.Data) {
		t.Errorf("expected %d players in PlayerResponse but received %d", len(players.Data), len(players.Data))
	}
}
func TestGetMatch(t *testing.T) {
	var id, player string
	options := GetPlayersRequestOptions{
		PlayerNamesFilter: testPlayers,
	}

	players, err := client.GetPlayers(options, PCNorthAmerica)
	if err != nil {
		t.Errorf("Error getting players from API: %s", err.Error())
	}

	if len(testPlayers) != len(players.Data) {
		t.Errorf("expected %d players in PlayerResponse but received %d", len(players.Data), len(players.Data))
	}

	for _, prd := range players.Data {
		if len(prd.GetMatchIDs()) > 0 {
			id = prd.GetMatchIDs()[0]
			player = prd.Attributes.Name
			return
		}
	}

	if id == "" {
		t.Errorf("no match ids for players %s", strings.Join(testPlayers, ","))
		return
	}

	match, err := client.GetMatch(id, PCNorthAmerica)

	if err != nil {
		t.Errorf("Error getting match from API: %s", err.Error())
	}

	if match.GetStatsByName()[player] == nil {
		t.Errorf("expected player %s to be in match but was not found", player)
	}
}

func TestGetSamples(t *testing.T) {
	samples, err := client.GetSampleMatches(PCNorthAmerica)

	if err != nil {
		t.Errorf("Error getting samples from API: %s", err.Error())
	}

	if len(samples.GetMatches()) == 0 {
		t.Errorf("Expected samples in SampleResponse but received 0")
	}
}

func TestGetSeasons(t *testing.T)            {}
func TestGetTelemetry(t *testing.T)          {}
func TestReadTelemetryFromFile(t *testing.T) {}

func setup() {
	b, _ := ioutil.ReadFile("conf.json")
	var conf Conf
	json.Unmarshal(b, &conf)
	var err error
	client, err = New(conf.Key, nil)
	if err != nil {
		panic("Error instantiating client")
	}
}
func teardown() {}
