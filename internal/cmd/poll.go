package cmd

import (
	"context"
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"go.albinodrought.com/neptunes-pride/internal/npapi"
	"go.albinodrought.com/neptunes-pride/internal/types"
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

			for _, config := range match.PlayerCreds {
				resp, err := client.State(context.Background(), &npapi.Request{
					GameNumber: gameNumber,
					APIKey:     config.APIKey,
				})

				if err != nil {
					log.Fatalf("failed fetching state for game %v user %v: %v", gameNumber, config.PlayerUID, err)
				}

				log.Printf("retrieved state for game %v user %v \"%v\"", gameNumber, config.PlayerUID, resp.ScanningData.Players[strconv.Itoa(config.PlayerUID)].Alias)
				lastResp = resp
			}

			if lastResp != nil {
				if match.Name == "" {
					match.Name = lastResp.ScanningData.Name
				}
			}

			err = db.SaveMatch(match)
			if err != nil {
				log.Fatalf("failed saving game %v: %v", gameNumber, err)
			}
		}
	},
}
