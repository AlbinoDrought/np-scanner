package cmd

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"go.albinodrought.com/neptunes-pride/internal/npapi"
	"go.albinodrought.com/neptunes-pride/internal/types"
)

var (
	pollCmdForce         bool
	pollCmdMinTimePassed time.Duration
)

var pollCmd = &cobra.Command{
	Use:   "poll [...game numbers]",
	Short: "Poll games",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		db, err := openDB()
		if err != nil {
			log.Fatal("failed to open DB", err)
		}

		if len(args) == 1 && args[0] == "all" {
			args, err = db.Matches()
			if err != nil {
				log.Fatal("failed getting all matches", err)
			}
		}

		client := openClient()

		for _, gameNumber := range args {
			match, err := db.FindMatchOrFail(gameNumber)
			if err != nil {
				log.Fatal("failed finding match", gameNumber, err)
			}

			if len(match.PlayerCreds) == 0 {
				log.Printf("match %v has no credentials", gameNumber)
				continue
			}

			var lastResp *types.APIResponse

			for i, config := range match.PlayerCreds {
				if !pollCmdForce && time.Since(config.LastPoll) < pollCmdMinTimePassed {
					log.Printf("recently polled game %v user %v on %v", gameNumber, config.PlayerUID, config.LastPoll)
					continue
				}

				resp, err := client.State(context.Background(), &npapi.Request{
					GameNumber: gameNumber,
					APIKey:     config.APIKey,
				})

				if err != nil {
					log.Fatalf("failed fetching state for game %v user %v: %v", gameNumber, config.PlayerUID, err)
				}

				config.LastPoll = time.Now()
				match.PlayerCreds[i] = config

				log.Printf("retrieved state for game %v user %v \"%v\"", gameNumber, config.PlayerUID, resp.ScanningData.Players[strconv.Itoa(config.PlayerUID)].Alias)
				lastResp = resp
			}

			if lastResp != nil {
				if match.Name == "" {
					match.Name = lastResp.ScanningData.Name
				}
			}

			match.LastPoll = time.Now()
			err = db.SaveMatch(match)
			if err != nil {
				log.Fatalf("failed saving game %v: %v", gameNumber, err)
			}
		}
	},
}

func init() {
	pollCmd.Flags().BoolVar(&pollCmdForce, "force", false, "Poll matches without regard for last-polled time")
	pollCmd.Flags().DurationVar(&pollCmdMinTimePassed, "min-time-passed", 15*time.Minute, "This much time must pass before we poll again")
}
