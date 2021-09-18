package cmd

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

var disablePlayerCmd = &cobra.Command{
	Use:   "disable-player [game number] [player-id]",
	Short: "Stop polling a player",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		db, err := openDB()
		if err != nil {
			log.Fatal("failed to open DB", err)
		}

		match, err := db.FindMatchOrFail(args[0])
		if err != nil {
			log.Fatal("failed to find match", err)
		}

		playerID, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatal("malformed player ID", err)
		}

		changed := false
		for i, playerCreds := range match.PlayerCreds {
			if playerCreds.PlayerUID == playerID {
				playerCreds.PollingDisabled = true
				match.PlayerCreds[i] = playerCreds
				changed = true
				break
			}
		}

		if !changed {
			log.Fatal("failed to find player", playerID)
		}

		err = db.SaveMatch(match)
		if err != nil {
			log.Fatal("failed to save match", err)
		}

		log.Println("disabled!")
	},
}
