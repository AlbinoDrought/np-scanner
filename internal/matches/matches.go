package matches

import "time"

type Match struct {
	GameNumber  string        `json:"game_number"`
	Name        string        `json:"name"`
	LastPoll    time.Time     `json:"last_poll"`
	PlayerCreds []PlayerCreds `json:"player_creds"`
}

type PlayerCreds struct {
	PlayerUID int    `json:"player_uid"`
	APIKey    string `json:"api_key"`
}
