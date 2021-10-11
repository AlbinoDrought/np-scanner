package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"
	"go.albinodrought.com/neptunes-pride/internal/web"
)

var (
	serveCmdAddress    string
	serveCmdPollPeriod time.Duration
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the API and web UI",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := openDB()
		if err != nil {
			log.Fatal("failed to open DB", err)
		}

		guard, err := openNotificationGuard()
		if err != nil {
			log.Fatal("failed to open notification guard", err)
		}

		sinks := buildSinks()

		client := openClient()

		err = web.Run(context.Background(), db, client, guard, sinks, &web.WebOptions{
			Address:    serveCmdAddress,
			PollPeriod: serveCmdPollPeriod,
		})
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	serveCmd.Flags().StringVar(&serveCmdAddress, "address", web.DefaultWebOptions.Address, "Address to listen on")
	serveCmd.Flags().DurationVar(&serveCmdPollPeriod, "poll-period", web.DefaultWebOptions.PollPeriod, "Check for match updates this often")
}
