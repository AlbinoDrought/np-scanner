package actions

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.albinodrought.com/neptunes-pride/internal/matchstore"
	"go.albinodrought.com/neptunes-pride/internal/multierror"
	"go.albinodrought.com/neptunes-pride/internal/npapi"
	"go.albinodrought.com/neptunes-pride/internal/types"
)

type PollOptions struct {
	Force         bool
	MinTimePassed time.Duration
}

var DefaultPollOptions = PollOptions{
	Force:         false,
	MinTimePassed: 15 * time.Minute,
}

type PollResult struct {
	Changed bool
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

func PollMatch(ctx context.Context, db matchstore.MatchStore, client npapi.NeptunesPrideClient, gameNumber string, pollOptions *PollOptions) (PollResult, error) {
	if pollOptions == nil {
		pollOptions = &DefaultPollOptions
	}

	pollResult := PollResult{}

	match, err := db.FindMatchOrFail(gameNumber)
	if err != nil {
		return pollResult, PollError{
			Base:       err,
			GameNumber: gameNumber,
			Message:    "failed finding match",
		}
	}

	if match.Finished {
		return pollResult, nil
	}

	if len(match.PlayerCreds) == 0 {
		log.Printf("match %v has no credentials", gameNumber)
		return pollResult, nil
	}

	pollErrors := []error{}

	for i, config := range match.PlayerCreds {
		if !pollOptions.Force && time.Since(config.LastPoll) < pollOptions.MinTimePassed {
			log.Printf("recently polled game %v user %v \"%v\" on %v", gameNumber, config.PlayerUID, config.PlayerAlias, config.LastPoll)
			continue
		}

		if config.PollingDisabled {
			log.Printf("polling is disabled for game %v user %v \"%v\" ", gameNumber, config.PlayerUID, config.PlayerAlias)
			continue
		}

		resp, err := client.State(ctx, &npapi.Request{
			GameNumber: gameNumber,
			APIKey:     config.APIKey,
		})

		if err != nil {
			pollErrors = append(pollErrors, PollError{
				Base:        err,
				GameNumber:  gameNumber,
				PlayerUID:   config.PlayerUID,
				PlayerAlias: config.PlayerAlias,
				Message:     "failed fetching remote state",
			})
			continue
		}

		err = db.SaveSnapshot(gameNumber, resp)
		if err != nil {
			pollErrors = append(pollErrors, PollError{
				Base:        err,
				GameNumber:  gameNumber,
				PlayerUID:   config.PlayerUID,
				PlayerAlias: config.PlayerAlias,
				Message:     "failed saving snapshot",
			})
			continue
		}

		playerData := resp.ScanningData.Players[strconv.Itoa(resp.ScanningData.PlayerUID)]
		if playerData.Conceded == types.ConcededWipedOut {
			// player's scanning data will no longer be updated: they're completely dead.
			// stop scanning it.
			config.PollingDisabled = true
			log.Printf("last poll for completely wiped out game %v user %v \"%v\"", gameNumber, config.PlayerUID, config.PlayerAlias)
		}

		if resp.ScanningData.Started && playerData.TotalStrength == 0 && playerData.TotalStars == 0 {
			// same as above, but API has not given player the "totally wiped out" status
			// stop scanning it.
			// (I think this happens if the player quits or goes AFK before total wipeout)
			// (players can still be active as AI when conceded is == 1 or 2)
			config.PollingDisabled = true
			log.Printf("last poll for assumed-wiped-out game %v user %v \"%v\"", gameNumber, config.PlayerUID, config.PlayerAlias)
		}

		config.LastPoll = time.Now()
		config.LatestSnapshot = resp.ScanningData.Now
		match.PlayerCreds[i] = config

		if resp.ScanningData.GameOver {
			match.Finished = true
			log.Printf("finished game %v user %v \"%v\"", gameNumber, config.PlayerUID, config.PlayerAlias)
		}

		pollResult.Changed = true
		log.Printf("retrieved state for game %v user %v \"%v\"", gameNumber, config.PlayerUID, config.PlayerAlias)
	}

	match.LastPoll = time.Now()
	err = db.SaveMatch(match)
	if err != nil {
		pollErrors = append(pollErrors, PollError{
			Base:       err,
			GameNumber: gameNumber,
			Message:    "failed saving game",
		})
	}

	return pollResult, multierror.Optional(pollErrors)
}

func PollMatches(ctx context.Context, db matchstore.MatchStore, client npapi.NeptunesPrideClient, gameNumbers []string, pollOptions *PollOptions) (map[string]PollResult, error) {
	pollResults := make(map[string]PollResult)
	pollErrors := []error{}

	for _, gameNumber := range gameNumbers {
		pollResult, err := PollMatch(ctx, db, client, gameNumber, pollOptions)
		pollResults[gameNumber] = pollResult

		if err != nil {
			pollErrors = append(pollErrors, err)
		}
	}

	return pollResults, multierror.Optional(pollErrors)
}

func PollAllMatches(ctx context.Context, db matchstore.MatchStore, client npapi.NeptunesPrideClient, pollOptions *PollOptions) (map[string]PollResult, error) {
	matches, err := db.Matches()
	if err != nil {
		return nil, err
	}

	return PollMatches(ctx, db, client, matches, pollOptions)
}
