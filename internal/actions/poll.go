package actions

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.albinodrought.com/neptunes-pride/internal/matchstore"
	"go.albinodrought.com/neptunes-pride/internal/npapi"
)

type PollOptions struct {
	Force         bool
	MinTimePassed time.Duration
}

var DefaultPollOptions = PollOptions{
	Force:         false,
	MinTimePassed: 15 * time.Minute,
}

type PollError struct {
	Base        error
	GameNumber  string
	PlayerUID   int
	PlayerAlias string
	Message     string
}

func (err PollError) Error() string {
	return fmt.Sprintf("%v: %+v | [gameNumber=%v] [playerUID=%v] [playerAlias=%v]", err.Message, err.Base, err.GameNumber, err.PlayerUID, err.PlayerAlias)
}

func PollMatch(ctx context.Context, db matchstore.MatchStore, client npapi.NeptunesPrideClient, gameNumber string, pollOptions *PollOptions) error {
	if pollOptions == nil {
		pollOptions = &DefaultPollOptions
	}

	match, err := db.FindMatchOrFail(gameNumber)
	if err != nil {
		return PollError{
			Base:       err,
			GameNumber: gameNumber,
			Message:    "failed finding match",
		}
	}

	if len(match.PlayerCreds) == 0 {
		log.Printf("match %v has no credentials", gameNumber)
		return nil
	}

	for i, config := range match.PlayerCreds {
		if !pollOptions.Force && time.Since(config.LastPoll) < pollOptions.MinTimePassed {
			log.Printf("recently polled game %v user %v \"%v\" on %v", gameNumber, config.PlayerUID, config.PlayerAlias, config.LastPoll)
			continue
		}

		resp, err := client.State(ctx, &npapi.Request{
			GameNumber: gameNumber,
			APIKey:     config.APIKey,
		})

		if err != nil {
			return PollError{
				Base:        err,
				GameNumber:  gameNumber,
				PlayerUID:   config.PlayerUID,
				PlayerAlias: config.PlayerAlias,
				Message:     "failed fetching remote state",
			}
		}

		err = db.SaveSnapshot(gameNumber, resp)
		if err != nil {
			return PollError{
				Base:        err,
				GameNumber:  gameNumber,
				PlayerUID:   config.PlayerUID,
				PlayerAlias: config.PlayerAlias,
				Message:     "failed saving snapshot",
			}
		}

		config.LastPoll = time.Now()
		config.LatestSnapshot = resp.ScanningData.Now
		match.PlayerCreds[i] = config

		log.Printf("retrieved state for game %v user %v \"%v\"", gameNumber, config.PlayerUID, config.PlayerAlias)
	}

	match.LastPoll = time.Now()
	err = db.SaveMatch(match)
	if err != nil {
		return PollError{
			Base:       err,
			GameNumber: gameNumber,
			Message:    "failed saving game",
		}
	}

	return nil
}

func PollMatches(ctx context.Context, db matchstore.MatchStore, client npapi.NeptunesPrideClient, gameNumbers []string, pollOptions *PollOptions) error {
	for _, gameNumber := range gameNumbers {
		if err := PollMatch(ctx, db, client, gameNumber, pollOptions); err != nil {
			return err
		}
	}

	return nil
}

func PollAllMatches(ctx context.Context, db matchstore.MatchStore, client npapi.NeptunesPrideClient, pollOptions *PollOptions) error {
	matches, err := db.Matches()
	if err != nil {
		return err
	}

	return PollMatches(ctx, db, client, matches, pollOptions)
}
