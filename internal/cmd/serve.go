package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"go.albinodrought.com/neptunes-pride/internal/web"
)

var (
	serveCmdAddress string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the API and web UI",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := openDB()
		if err != nil {
			log.Fatal("failed to open DB", err)
		}

		client := openClient()

		err = web.Run(context.Background(), db, client, &web.WebOptions{
			Address: serveCmdAddress,
		})
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	serveCmd.Flags().StringVar(&serveCmdAddress, "address", ":38080", "Address to listen on")
}
