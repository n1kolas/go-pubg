package pubg

// GetPlayersRequestOptions Filter options for the get players endpoint
type GetPlayersRequestOptions struct {
	// PlayerNamesFilter Filters by player name. Usage: filter[playerNames]=player1,player2,…
	PlayerNamesFilter []string `url:"filter[playerNames],comma,omitempty"`
	// PlayerIDsFilter Filters by player Id. Usage:filter[playerIds]=playerId,playerId,…
	PlayerIDsFilter []string `url:"filter[playerIds],comma,omitempty"`
}
