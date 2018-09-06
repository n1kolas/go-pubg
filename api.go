package pubg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

const (
	shards  = "shards"  // shards path segment
	matches = "matches" // matches end point
	samples = "samples"
	players = "players" // players end point
	status  = "status"  // status end point
	seasons = "seasons" // seasons end point

	// XboxAsia - Xbox Asia Region
	XboxAsia = "xbox-as"
	// XboxEurope - Xbox Europe Region
	XboxEurope = "xbox-eu"
	// XboxNorthAmerica - Xbox North America Region
	XboxNorthAmerica = "xbox-na"
	// XboxOceania - Xbox Oceana Region
	XboxOceania = "xbox-oc"
	// PCAsia - PC Asia  Region
	PCAsia = "pc-as"
	// PCEurope - PC Europe Region
	PCEurope = "pc-eu"
	// PCNorthAmerica - PC North America Region
	PCNorthAmerica = "pc-na"
	// PCOceania - PC Oceania Region
	PCOceania = "pc-oc"
	// PCKoreaJapan - PC Korea/Japan Region
	PCKoreaJapan = "pc-krjp"
	// PCKorea - PC Korea Region
	PCKorea = "pc-kr"
	// PCJapan - PC Japan Region
	PCJapan = "pc-jp"
	// PCKAKAO - PC KAKAO Region
	PCKAKAO = "pc-kakao"
	// PCSouthEastAsia - PC South East Asia Region
	PCSouthEastAsia = "pc-sea"
	// PCSouthAsia - PC South Asia Region
	PCSouthAsia = "pc-sa"
	// PCTournament - PC Tournament Shard
	PCTournament = "pc-tournament"
)

func (c *Client) newRequest(method, path string, body interface{}, options interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	url := u.String()
	if options != nil {
		v, err := query.Values(options)
		if err != nil {
			return nil, err
		}
		url = fmt.Sprintf("%s?%s", url, v.Encode())
	}
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.apiKey)
	req.Header.Set("Accept", "application/vnd.api+json")
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var buffer bytes.Buffer
	buffer.ReadFrom(resp.Body)

	switch resp.StatusCode {
	case http.StatusOK:
		return json.Unmarshal(buffer.Bytes(), v)
	case http.StatusUnauthorized:
		return NewInvalidKeyError(req.URL.String())
	case http.StatusNotFound:
		return NewNotFoundError(req.URL.String())
	case http.StatusUnsupportedMediaType:
		return NewIncorrectContentTypeError(req.URL.String())
	case http.StatusTooManyRequests:
		return NewTooManyRequestsError(req.URL.String())
	default:
		return NewUnhandledStatusCodeError(req.URL.String(), resp.Status)
	}
}

// GetStatus retrieves status data from the PUBG servers
func (c *Client) GetStatus() (*StatusResponse, error) {
	req, err := c.newRequest("GET", status, nil, nil)
	if err != nil {
		return nil, err
	}
	var response StatusResponse
	err = c.do(req, &response)
	return &response, err
}

// GetPlayer retrieves player data.
func (c *Client) GetPlayer(id, shard string) (*PlayerResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("%s/%s/%s/%s", shards, shard, players, id), nil, nil)
	if err != nil {
		return nil, err
	}

	var response PlayerResponse
	err = c.do(req, &response)
	return &response, err
}

// GetPlayers retrieves data for players from the passed options.
func (c *Client) GetPlayers(options GetPlayersRequestOptions, shard string) (*PlayersResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("%s/%s/%s", shards, shard, players), nil, options)
	if err != nil {
		return nil, err
	}

	var response PlayersResponse
	err = c.do(req, &response)
	return &response, err
}

// GetSeasons retrieves data about seasons in a shard
func (c *Client) GetSeasons(shard string) (*SeasonsResponse, error) {
	req, err := c.newRequest("GET",
		fmt.Sprintf("%s/%s/%s", shards, shard, seasons),
		nil, nil,
	)
	if err != nil {
		return nil, err
	}
	var response SeasonsResponse
	err = c.do(req, &response)
	return &response, err
}

// GetSeasonStats retrieves data about a season
func (c *Client) GetSeasonStats(playerID string, shard string, season string) (*PlayerSeasonResponse, error) {
	req, err := c.newRequest("GET",
		fmt.Sprintf("%s/%s/%s/%s/%s/%s", shards, shard, players, playerID, seasons, season),
		nil, nil,
	)
	if err != nil {
		return nil, err
	}
	var response PlayerSeasonResponse
	err = c.do(req, &response)
	return &response, err
}

// GetSampleMatches retrieves samples matches
func (c *Client) GetSampleMatches(shard string) (*SamplesResponse, error) {
	req, err := c.newRequest("GET",
		fmt.Sprintf("%s/%s/%s", shards, shard, samples),
		nil, nil,
	)
	if err != nil {
		return nil, err
	}
	var response SamplesResponse
	err = c.do(req, &response)
	return &response, err
}

// GetMatch retrieves the data for a match with a given id and shard
func (c *Client) GetMatch(id string, shard string) (*MatchResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("%s/%s/%s/%s", shards, shard, matches, id), nil, nil)
	if err != nil {
		return nil, err
	}

	var response MatchResponse
	err = c.do(req, &response)
	if err != nil {
		return nil, err
	}

	for _, inc := range response.Included {
		var check map[string]string
		json.Unmarshal(inc, &check)
		switch check["type"] {
		case "participant":
			var p MatchParticipant
			json.Unmarshal(inc, &p)
			response.Participants = append(response.Participants, p)
		case "asset":
			var a MatchAsset
			json.Unmarshal(inc, &a)
			response.Assets = append(response.Assets, a)
		case "roster":
			var r MatchRoster
			json.Unmarshal(inc, &r)
			response.Rosters = append(response.Rosters, r)
		}

		if err != nil {
			return nil, err
		}
	}

	return &response, err
}

// GetTelemetry retrieves the telemetry data at a specified url
func (c *Client) GetTelemetry(url string) (*TelemetryResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	buffer.ReadFrom(resp.Body)
	return parseTelemetry(buffer.Bytes())
}

// ReadTelemetryFromFile parses json telemetry data from a given file
// and returns a TelemetryResponse struct. It is more performant to cache
// telemetry data for future use.
func ReadTelemetryFromFile(path string) (tr *TelemetryResponse, err error) {
	var b []byte
	b, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}
	return parseTelemetry(b)
}

// ParseTelemetry reads the telemetry event type from the json
// and passes it to the unmarshaller
func parseTelemetry(b []byte) (tr *TelemetryResponse, err error) {
	var v []json.RawMessage
	json.Unmarshal(b, &v)
	for _, bts := range v {
		var eval map[string]interface{}
		err = json.Unmarshal(bts, &eval)
		if err != nil {
			return
		}
		tr.unmarshalEvent(bts, eval["_T"].(string))
	}
	return
}
