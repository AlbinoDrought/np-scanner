package matches

import "time"

type Match struct {
	GameNumber  string              `json:"game_number"`
	Name        string              `json:"name"`
	LastPoll    time.Time           `json:"last_poll"`
	PlayerCreds map[int]PlayerCreds `json:"player_creds"`
}

type PlayerCreds struct {
	PlayerUID int       `json:"player_uid"`
	APIKey    string    `json:"api_key"`
	LastPoll  time.Time `json:"last_poll"`
}

func NewMatch(gameNumber string) *Match {
	return &Match{
		GameNumber:  gameNumber,
		Name:        "",
		PlayerCreds: map[int]PlayerCreds{},
	}
}
