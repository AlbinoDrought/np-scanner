package cmd

import (
	"net/http"

	"github.com/spf13/cobra"

	"go.albinodrought.com/neptunes-pride/internal/matchstore"
	"go.albinodrought.com/neptunes-pride/internal/npapi"
)

var (
	dbPath string
)

func addGlobalConfigFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&dbPath, "db-path", "np.db", "Database Path")
}

func openDB() (matchstore.MatchStore, error) {
	return matchstore.Open(dbPath)
}

func openClient() npapi.NeptunesPrideClient {
	return npapi.NewClient(http.DefaultClient)
}
