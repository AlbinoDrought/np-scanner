package cmd

import (
	"context"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"go.albinodrought.com/neptunes-pride/internal/npapi"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch [game number] [key]",
	Short: "Fetch game details immediately",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := npapi.NewClient(http.DefaultClient)
		resp, err := client.State(context.Background(), &npapi.Request{
			GameNumber: args[0],
			APIKey:     args[0],
		})

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("%+v", resp)
	},
}
