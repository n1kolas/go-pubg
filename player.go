package pubg

// GetMatches is a helper function to extract a slice of MatchData
// from PlayerResponseData. This is useful if the caller is only
// interested in MatchData
func (prd *PlayerResponseData) GetMatches() (m []MatchData) {
	m = prd.Relationships.Matches.Data
	return
}

// GetMatchIDs is a helper function to exctract a slice of MatchId's
// from PlayerResponseData. This is useful if the caller is only
// interested in the MatchId's for a player. This helper function
// becomes very handy when trying to work with Telemetry data
func (prd *PlayerResponseData) GetMatchIDs() (ids []string) {
	for _, d := range prd.Relationships.Matches.Data {
		ids = append(ids, d.ID)
	}
	return ids
}
