package matches

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AccessProfile struct {
	Code               []byte
	CanViewEveryPlayer bool
	AllowedPlayerIDs   []int
}

func (ap *AccessProfile) CheckCode(plaintext []byte) error {
	return bcrypt.CompareHashAndPassword(ap.Code, plaintext)
}

func (ap *AccessProfile) CanViewPlayerID(id int) bool {
	if ap.CanViewEveryPlayer {
		return true
	}

	if ap.AllowedPlayerIDs != nil {
		for _, allowedPlayerID := range ap.AllowedPlayerIDs {
			if id == allowedPlayerID {
				return true
			}
		}
	}

	return false
}

func NewAccessProfile(plaintext []byte) (AccessProfile, error) {
	hash, err := bcrypt.GenerateFromPassword(plaintext, bcrypt.DefaultCost)
	return AccessProfile{
		Code:             hash,
		AllowedPlayerIDs: []int{},
	}, err
}

func PermissiveAccessProfile() AccessProfile {
	return AccessProfile{
		CanViewEveryPlayer: true,
		AllowedPlayerIDs:   []int{},
	}
}

type Match struct {
	GameNumber     string              `json:"game_number"`
	Finished       bool                `json:"finished"`
	Name           string              `json:"name"`
	LastPoll       time.Time           `json:"last_poll"`
	PlayerCreds    map[int]PlayerCreds `json:"player_creds,omitempty"`
	OldAccessCode  []byte              `json:"access_code,omitempty"`
	AccessProfiles []AccessProfile     `json:"access_profiles,omitempty"`
}

func (match *Match) HasAccessCode() bool {
	// old access code
	if match.OldAccessCode != nil && len(match.OldAccessCode) > 0 {
		return true
	}

	// new access profiles
	if match.AccessProfiles != nil && len(match.AccessProfiles) > 0 {
		return true
	}

	return false
}

var ErrNoAccessCodes = errors.New("no access codes to check")

func (match *Match) CheckAccessCode(plaintext []byte) (AccessProfile, error) {
	err := ErrNoAccessCodes

	if match.OldAccessCode != nil {
		accessProfile := PermissiveAccessProfile()
		accessProfile.Code = match.OldAccessCode
		err = accessProfile.CheckCode(plaintext)
		if err == nil {
			// old access code verified, return permissive access profile
			return accessProfile, nil
		}
	}

	if match.AccessProfiles != nil {
		for _, accessProfile := range match.AccessProfiles {
			err = accessProfile.CheckCode(plaintext)
			if err == nil {
				return accessProfile, nil
			}
		}
	}

	return AccessProfile{}, err
}

func (match *Match) WipeAccessCodes() {
	match.OldAccessCode = nil
	match.AccessProfiles = nil
}

var ErrAccessCodeAlreadyUsed = errors.New("access code is already used for another access profile")

func (match *Match) AddAccessProfile(ap AccessProfile, plaintext []byte) error {
	// plaintext must work for this profile
	err := ap.CheckCode(plaintext)
	if err != nil {
		return err
	}

	if match.AccessProfiles == nil {
		match.AccessProfiles = []AccessProfile{}
	}

	// ensure no other profiles already use this plaintext
	_, err = match.CheckAccessCode(plaintext)
	if err == nil {
		return ErrAccessCodeAlreadyUsed
	}

	match.AccessProfiles = append(match.AccessProfiles, ap)

	return nil
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
