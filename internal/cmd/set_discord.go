package cmd

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	setDiscordCmdWipe bool
)

var setDiscordCmd = &cobra.Command{
	Use:   "set-discord [game number] [player uid] [discord user id]",
	Short: "Associate a player with their Discord user ID",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		db, err := openDB()
		if err != nil {
			log.Fatal("failed to open DB: ", err)
		}

		match, err := db.FindMatchOrFail(args[0])
		if err != nil {
			log.Fatal("failed finding match: ", err)
		}

		if setDiscordCmdWipe {
			match.DiscordUserIDs = nil
			log.Println("wiped old discord user IDs")
		}

		playerUID, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatal("failed parsing player uid: ", err)
		}

		if match.DiscordUserIDs == nil {
			match.DiscordUserIDs = map[int]string{}
		}

		match.DiscordUserIDs[playerUID] = args[2]

		err = db.SaveMatch(match)
		if err != nil {
			log.Fatal("failed saving match: ", err)
		}

		log.Println("set discord user ID mapping for match", match.GameNumber, "player", playerUID, "to", args[2])
	},
}

func init() {
	setDiscordCmd.Flags().BoolVar(&setDiscordCmdWipe, "wipe", false, "Remove all other Discord user IDs before setting this one")
}
