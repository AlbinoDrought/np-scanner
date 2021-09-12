package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "np-scanner",
	Short: "NP Scanner does stuff with th Neptune's Pride API",
	Long: `NP Scanner does stuff with th Neptune's Pride API.
                See np.ironhelmet.com, github.com/AlbinoDrought`,
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	addGlobalConfigFlags(rootCmd)
	rootCmd.AddCommand(pollCmd)
	rootCmd.AddCommand(protectCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(serveCmd)
}
