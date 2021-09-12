package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"go.albinodrought.com/neptunes-pride/internal/actions"
)

var setCmd = &cobra.Command{
	Use:   "set [game number] [key]",
	Short: "Set player credentials for a game",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		db, err := openDB()
		if err != nil {
			log.Fatal("failed to open DB", err)
		}

		client := openClient()

		err = actions.SetCredentials(context.Background(), db, client, args[0], args[1])
		if err != nil {
			log.Fatal(err)
		}
	},
}
