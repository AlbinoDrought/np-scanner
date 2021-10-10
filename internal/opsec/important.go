package opsec

import (
	"strconv"

	"go.albinodrought.com/neptunes-pride/internal/types"
)

type Threat struct {
	Fleet        *types.Fleet
	Order        []int
	FleetOwnerID string
	FleetOwner   *types.Player
	TargetStarID string
	TargetStar   *types.Star
	// todo: TargetStartTrueStrength
	TargetStarOwnerID string
	TargetStarOwner   *types.Player
}

func FindThreats(resp *types.APIResponse) []Threat {
	threats := []Threat{}

	for _, fleet := range resp.ScanningData.Fleets {
		fleetOwnerID := strconv.Itoa(fleet.PlayerID)
		fleetOwner, ok := resp.ScanningData.Players[fleetOwnerID]
		if !ok {
			// can't find fleet owner, ignore
			// (should probably never happen)
			continue
		}

		for _, order := range fleet.Orders {
			if len(order) < 2 {
				// malformed order, ignore
				continue
			}

			targetStarID := strconv.Itoa(order[1])
			targetStar, ok := resp.ScanningData.Stars[targetStarID]
			if !ok {
				// can't find star, ignore
				continue
			}

			if targetStar.PlayerID == -1 || targetStar.PlayerID == fleet.PlayerID {
				// star is unowned or owned by same player, ignore
				continue
			}

			targetStarOwnerID := strconv.Itoa(targetStar.PlayerID)
			targetStarOwner, ok := resp.ScanningData.Players[targetStarOwnerID]
			if !ok {
				// can't find star owner, ignore
				// (should probably never happen)
				continue
			}

			threats = append(threats, Threat{
				Fleet:             &fleet,
				Order:             order,
				FleetOwnerID:      fleetOwnerID,
				FleetOwner:        &fleetOwner,
				TargetStarID:      targetStarID,
				TargetStar:        &targetStar,
				TargetStarOwnerID: targetStarOwnerID,
				TargetStarOwner:   &targetStarOwner,
			})
		}
	}

	return threats
}
