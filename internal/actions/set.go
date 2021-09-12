package actions

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"go.albinodrought.com/neptunes-pride/internal/matches"
	"go.albinodrought.com/neptunes-pride/internal/matchstore"
	"go.albinodrought.com/neptunes-pride/internal/npapi"
)

type SetCredentialsError struct {
	Base       error
	GameNumber string
	PlayerUID  int
	Message    string
}

func (err SetCredentialsError) Error() string {
	return fmt.Sprintf("%v: %+v | [gameNumber=%v] [playerUID=%v]", err.Message, err.Base, err.GameNumber, err.PlayerUID)
}

func SetCredentials(ctx context.Context, db matchstore.MatchStore, client npapi.NeptunesPrideClient, gameNumber string, key string) error {
	// validate credentials
	resp, err := client.State(ctx, &npapi.Request{
		GameNumber: gameNumber,
		APIKey:     key,
	})

	if err != nil {
		return SetCredentialsError{
			Base:       err,
			GameNumber: gameNumber,
			Message:    "failed to validate credentials",
		}
	}

	// credentials are OK, save them
	playerID := resp.ScanningData.PlayerUID

	match, err := db.FindOrCreateMatch(gameNumber)
	if err != nil {
		return SetCredentialsError{
			Base:       err,
			GameNumber: gameNumber,
			PlayerUID:  playerID,
			Message:    "failed to find or create match",
		}
	}

	if match.Name == "" {
		match.Name = resp.ScanningData.Name
	}

	match.PlayerCreds[playerID] = matches.PlayerCreds{
		PlayerUID:   playerID,
		PlayerAlias: resp.ScanningData.Players[strconv.Itoa(playerID)].Alias,
		APIKey:      key,
	}

	err = db.SaveMatch(match)
	if err != nil {
		return SetCredentialsError{
			Base:       err,
			GameNumber: gameNumber,
			PlayerUID:  playerID,
			Message:    "failed to save match",
		}
	}

	log.Printf("set credentials for player %v in game %v", playerID, match.GameNumber)
	return nil
}
