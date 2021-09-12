package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var protectCmd = &cobra.Command{
	Use:   "protect [game number] [code]",
	Short: "Protect a game with an access code",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		db, err := openDB()
		if err != nil {
			log.Fatal("failed to open DB: ", err)
		}

		match, err := db.FindMatchOrFail(args[0])
		if err != nil {
			log.Fatal("failed finding match: ", err)
		}

		err = match.SetAccessCode([]byte(args[1]))
		if err != nil {
			log.Fatal("failed setting access code: ", err)
		}

		err = db.SaveMatch(match)
		if err != nil {
			log.Fatal("failed saving match: ", err)
		}

		log.Println("set access code for match", match.GameNumber, "to", args[1])
	},
}
