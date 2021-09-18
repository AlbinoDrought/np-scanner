package matches

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Match struct {
	GameNumber  string              `json:"game_number"`
	Name        string              `json:"name"`
	LastPoll    time.Time           `json:"last_poll"`
	PlayerCreds map[int]PlayerCreds `json:"player_creds,omitempty"`
	AccessCode  []byte              `json:"access_code,omitempty"`
}

func (match *Match) HasAccessCode() bool {
	return match.AccessCode != nil && len(match.AccessCode) > 0
}

func (match *Match) SetAccessCode(plaintext []byte) error {
	hash, err := bcrypt.GenerateFromPassword(plaintext, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	match.AccessCode = hash
	return nil
}

func (match *Match) CheckAccessCode(plaintext []byte) error {
	return bcrypt.CompareHashAndPassword(match.AccessCode, plaintext)
}

type PlayerCreds struct {
	PlayerUID       int       `json:"player_uid"`
	PlayerAlias     string    `json:"player_alias"`
	APIKey          string    `json:"api_key,omitempty"`
	LastPoll        time.Time `json:"last_poll"`
	LatestSnapshot  int64     `json:"latest_snapshot"`
	PollingDisabled bool      `json:"polling_disabled"`
}

func NewMatch(gameNumber string) *Match {
	return &Match{
		GameNumber:  gameNumber,
		Name:        "",
		PlayerCreds: map[int]PlayerCreds{},
	}
}
