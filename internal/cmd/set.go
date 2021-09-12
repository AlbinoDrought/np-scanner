package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"go.albinodrought.com/neptunes-pride/internal/matches"
	"go.albinodrought.com/neptunes-pride/internal/npapi"
)

var setCmd = &cobra.Command{
	Use:   "set [game number] [key]",
	Short: "Set player credentials for a game",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// validate credentials
		client := openClient()
		resp, err := client.State(context.Background(), &npapi.Request{
			GameNumber: args[0],
			APIKey:     args[1],
		})

		if err != nil {
			log.Fatal("failed to validate credentials", err)
		}

		// credentials are OK, save them
		playerID := resp.ScanningData.PlayerUID

		db, err := openDB()
		if err != nil {
			log.Fatal("failed to open DB", err)
		}

		match, err := db.FindOrCreateMatch(args[0])
		if err != nil {
			log.Fatal("failed to find or create match", err)
		}

		match.PlayerCreds[playerID] = matches.PlayerCreds{
			PlayerUID: playerID,
			APIKey:    args[1],
		}

		err = db.SaveMatch(match)
		if err != nil {
			log.Fatal("failed to save match", err)
		}

		log.Printf("set credentials for player %v in game %v", playerID, match.GameNumber)
	},
}
