package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"
	"go.albinodrought.com/neptunes-pride/internal/actions"
)

var (
	pollCmdForce         bool
	pollCmdMinTimePassed time.Duration
)

var pollCmd = &cobra.Command{
	Use:   "poll [...game numbers, or \"all\"]",
	Short: "Poll games",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		db, err := openDB()
		if err != nil {
			log.Fatal("failed to open DB", err)
		}

		client := openClient()

		pollOptions := &actions.PollOptions{
			Force:         pollCmdForce,
			MinTimePassed: pollCmdMinTimePassed,
		}

		if len(args) == 1 && args[0] == "all" {
			err = actions.PollAllMatches(context.Background(), db, client, pollOptions)
		} else {
			err = actions.PollMatches(context.Background(), db, client, args, pollOptions)
		}

		if err != nil {
			log.Fatal("polling failed", err)
		}

		log.Println("polled")
	},
}

func init() {
	pollCmd.Flags().BoolVar(&pollCmdForce, "force", false, "Poll matches without regard for last-polled time")
	pollCmd.Flags().DurationVar(&pollCmdMinTimePassed, "min-time-passed", 15*time.Minute, "This much time must pass before we poll again")
}
