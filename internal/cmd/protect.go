package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"go.albinodrought.com/neptunes-pride/internal/matches"
)

var (
	protectCmdWipe        bool
	protectCmdAllowedUIDs []int
)

var protectCmd = &cobra.Command{
	Use:   "protect [game number] [code]",
	Short: "Protect a game with an access profile",
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

		if protectCmdWipe {
			match.WipeAccessCodes()
			log.Println("wiped old access code and profiles")
		}

		plaintext := []byte(args[1])

		newProfile, err := matches.NewAccessProfile(plaintext)
		if err != nil {
			log.Fatal("failed creating new access profile: ", err)
		}

		if len(protectCmdAllowedUIDs) > 0 {
			log.Println("new profile can only view these players: ", protectCmdAllowedUIDs)
			newProfile.AllowedPlayerIDs = protectCmdAllowedUIDs
		} else {
			log.Println("new profile can view every player")
			newProfile.CanViewEveryPlayer = true
		}

		err = match.AddAccessProfile(newProfile, plaintext)
		if err != nil {
			log.Fatal("failed adding access profile: ", err)
		}

		err = db.SaveMatch(match)
		if err != nil {
			log.Fatal("failed saving match: ", err)
		}

		log.Println("set access code for match", match.GameNumber, "to", args[1])
	},
}

func init() {
	protectCmd.Flags().BoolVar(&protectCmdWipe, "wipe", false, "Remove all other access profiles before adding this one")
	protectCmd.Flags().IntSliceVar(&protectCmdAllowedUIDs, "allowed-uid", []int{}, "If set, access profile can only see these whitelisted users")
}
