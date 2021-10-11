package cmd

import (
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"go.albinodrought.com/neptunes-pride/internal/matchstore"
	"go.albinodrought.com/neptunes-pride/internal/notifications"
	"go.albinodrought.com/neptunes-pride/internal/npapi"
)

var (
	matchStoreDbPath    string
	notificationsDbPath string
	discordWebhookURL   string
)

func addGlobalConfigFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&matchStoreDbPath, "db-path", "np.db", "Database Path (Match Store)")
	cmd.PersistentFlags().StringVar(&notificationsDbPath, "notifications-db-path", "np-notifications.db", "Database Path (Notifications)")
	cmd.PersistentFlags().StringVar(&discordWebhookURL, "discord-webhook-url", os.Getenv("NP_SCANNER_DISCORD_WEBHOOK_URL"), "Discord Webhook URL")
}

func openDB() (matchstore.MatchStore, error) {
	return matchstore.Open(matchStoreDbPath)
}

func openNotificationGuard() (notifications.Guard, error) {
	return notifications.OpenBoltGuard(notificationsDbPath)
}

func buildSinks() []notifications.Sink {
	sinks := []notifications.Sink{}

	if discordWebhookURL != "" {
		sinks = append(sinks, notifications.NewDiscordWebhookSink(discordWebhookURL, http.DefaultClient))
	}

	return sinks
}

func openClient() npapi.NeptunesPrideClient {
	return npapi.NewClient(http.DefaultClient)
}
