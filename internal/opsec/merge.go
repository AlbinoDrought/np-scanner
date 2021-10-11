package opsec

import (
	"sort"

	"go.albinodrought.com/neptunes-pride/internal/types"
)

func Merge(responses ...*types.APIResponse) *types.APIResponse {
	if len(responses) == 0 {
		return nil
	}

	sort.Slice(responses, func(i, j int) bool {
		return responses[i].ScanningData.Now < responses[j].ScanningData.Now
	})

	base := responses[0]

	for i, other := range responses {
		if i == 0 {
			continue
		}

		// add missing fleets
		for fleetIndex, fleet := range other.ScanningData.Fleets {
			if _, ok := base.ScanningData.Fleets[fleetIndex]; !ok {
				base.ScanningData.Fleets[fleetIndex] = fleet
			}
		}

		// more info on stars
		for starIndex, star := range other.ScanningData.Stars {
			if _, ok := base.ScanningData.Stars[starIndex]; !ok {
				base.ScanningData.Stars[starIndex] = star
			}

			if star.PrivateStar.Useful() {
				baseStar := base.ScanningData.Stars[starIndex]
				baseStar.PrivateStar = star.PrivateStar
				baseStar.Visible = star.Visible
				base.ScanningData.Stars[starIndex] = baseStar
			}
		}

		// more info on players
		for playerIndex, player := range other.ScanningData.Players {
			if _, ok := base.ScanningData.Players[playerIndex]; !ok {
				base.ScanningData.Players[playerIndex] = player
			}

			if player.PrivatePlayer.Useful() {
				basePlayer := base.ScanningData.Players[playerIndex]
				basePlayer.PrivatePlayer = player.PrivatePlayer

				// assume this also means private tech data needs to be imported
				basePlayer.Tech.Scanning.PrivateTechResearchStatus = player.Tech.Scanning.PrivateTechResearchStatus
				basePlayer.Tech.Propulsion.PrivateTechResearchStatus = player.Tech.Propulsion.PrivateTechResearchStatus
				basePlayer.Tech.Terraforming.PrivateTechResearchStatus = player.Tech.Terraforming.PrivateTechResearchStatus
				basePlayer.Tech.Research.PrivateTechResearchStatus = player.Tech.Research.PrivateTechResearchStatus
				basePlayer.Tech.Weapons.PrivateTechResearchStatus = player.Tech.Weapons.PrivateTechResearchStatus
				basePlayer.Tech.Banking.PrivateTechResearchStatus = player.Tech.Banking.PrivateTechResearchStatus
				basePlayer.Tech.Manufacturing.PrivateTechResearchStatus = player.Tech.Manufacturing.PrivateTechResearchStatus

				base.ScanningData.Players[playerIndex] = basePlayer
			}
		}
	}

	return base
}
