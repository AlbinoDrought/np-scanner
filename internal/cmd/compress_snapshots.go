package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var compressSnapshotsCmd = &cobra.Command{
	Use:   "compress-snapshots",
	Short: "Poll games",
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		db, err := openDB()
		if err != nil {
			log.Fatal("failed to open DB", err)
		}
		if err := db.CompressSnapshots(log.Println); err != nil {
			log.Fatal("failed to compress snapshots", err)
		}
		log.Println("ok")
	},
}
