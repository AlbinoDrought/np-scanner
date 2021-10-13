package actions

import (
	"fmt"

	"go.albinodrought.com/neptunes-pride/internal/matches"
	"go.albinodrought.com/neptunes-pride/internal/notifications"
	"go.albinodrought.com/neptunes-pride/internal/opsec"
	"go.albinodrought.com/neptunes-pride/internal/types"
)

type notifiableThreat struct {
	baseID string
	threat *opsec.Threat
	match  *matches.Match
}

func (t *notifiableThreat) ID() string {
	return fmt.Sprintf("threat-%v-%v-%v-%v", t.baseID, t.threat.Fleet.UID, t.threat.Fleet.Strength, t.threat.TargetStarID)
}

func (t *notifiableThreat) createMessage(fleetOwner string, targetStarOwner string) string {
	return fmt.Sprintf(
		"%v's carrier %v is attacking %v's star %v with %v units",
		fleetOwner,
		t.threat.Fleet.Name,
		targetStarOwner,
		t.threat.TargetStar.Name,
		t.threat.Fleet.Strength,
	)
}

func (t *notifiableThreat) Message() string {
	return t.createMessage(t.threat.FleetOwner.Alias, t.threat.TargetStarOwner.Alias)
}

func (t *notifiableThreat) DiscordMessage() string {
	fleetOwner := t.threat.FleetOwner.Alias
	targetStarOwner := t.threat.TargetStarOwner.Alias

	if t.match.DiscordUserIDs != nil {
		/*
			// @mentions are useful for targeted notifications
			// notifying users of their own attacks against others is less useful
			fleetOwnerDiscordID, ok := t.match.DiscordUserIDs[t.threat.FleetOwner.UID]
			if ok {
				fleetOwner = fmt.Sprintf("<@%v>", fleetOwnerDiscordID)
			}
		*/

		targetStarOwnerDiscordID, ok := t.match.DiscordUserIDs[t.threat.TargetStarOwner.UID]
		if ok {
			targetStarOwner = fmt.Sprintf("<@%v>", targetStarOwnerDiscordID)
		}
	}

	return t.createMessage(fleetOwner, targetStarOwner)
}

func CheckNotifiables(match *matches.Match, resp *types.APIResponse) []notifications.Notifiable {
	notifiables := []notifications.Notifiable{}

	threats := opsec.FindThreats(resp)
	for i := range threats { // uses index to avoid loop variable overwriting
		notifiables = append(notifiables, &notifiableThreat{
			baseID: fmt.Sprintf("%v-%v", match.GameNumber, resp.ScanningData.Productions), // max 1 notification per production
			threat: &threats[i],
			match:  match,
		})
	}

	return notifiables
}
