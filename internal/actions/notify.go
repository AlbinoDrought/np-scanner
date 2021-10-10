package actions

import (
	"fmt"

	"go.albinodrought.com/neptunes-pride/internal/notifications"
	"go.albinodrought.com/neptunes-pride/internal/opsec"
	"go.albinodrought.com/neptunes-pride/internal/types"
)

type notifiableThreat struct {
	baseID string
	threat *opsec.Threat
}

func (t *notifiableThreat) ID() string {
	return fmt.Sprintf("threat-%v-%v-%v-%v", t.baseID, t.threat.Fleet.UID, t.threat.Fleet.Strength, t.threat.TargetStarID)
}

func (t *notifiableThreat) Message() string {
	return fmt.Sprintf(
		"%v's carrier %v is attacking %v's star %v with %v units",
		t.threat.FleetOwner.Alias,
		t.threat.Fleet.Name,
		t.threat.TargetStarOwner.Alias,
		t.threat.TargetStar.Name,
		t.threat.Fleet.Strength,
	)
}

func CheckNotifiables(gameNumber string, resp *types.APIResponse) []notifications.Notifiable {
	notifiables := []notifications.Notifiable{}

	threats := opsec.FindThreats(resp)
	for i := range threats { // uses index to avoid loop variable overwriting
		notifiables = append(notifiables, &notifiableThreat{
			baseID: fmt.Sprintf("%v-%v", gameNumber, resp.ScanningData.Productions), // max 1 notification per production
			threat: &threats[i],
		})
	}

	return notifiables
}
