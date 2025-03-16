package types

import (
	"encoding/json"
	"testing"
)

func TestApiDeserialize(t *testing.T) {
	var apiResp APIResponse

	err := json.Unmarshal([]byte(`
{
  "scanning_data": {
    "playerUid": 3,
    "config": {
      "version": "np4",
      "name": "foo",
      "description": "",
      "password": "snipped",
      "players": 6,
      "starsForVictory": 50,
      "sfvDecay": 0,
      "buildGates": 1,
      "randomGates": 0,
      "randomWorms": 0,
      "playerType": 0,
      "alliances": 1,
      "anonymity": 1,
      "darkGalaxy": 3,
      "starsPerPlayer": 24,
      "homeStarDistance": 3,
      "naturalResources": 2,
      "starfield": "hexgrid",
      "starScatter": "random",
      "customStarfield": "",
      "prodTicks": 20,
      "mirror": 0,
      "startStars": 6,
      "startCash": 500,
      "startShips": 10,
      "startInfEco": 10,
      "startInfInd": 5,
      "startInfSci": 2,
      "devCostEco": 2,
      "devCostInd": 2,
      "devCostSci": 2,
      "tradeCost": 25,
      "tradeScanned": 0,
      "fleetCost": 25,
      "newBnk": 1,
      "newRng": 1,
      "expBonus": 72,
      "noExp": 0,
      "noScn": 1,
      "noTer": 1,
      "startTechTer": 1,
      "startTechWep": 1,
      "startTechScn": 1,
      "startTechRng": 1,
      "startTechMan": 1,
      "startTechExp": 1,
      "startTechBnk": 1,
      "resCostTer": 2,
      "resCostWep": 2,
      "resCostScn": 2,
      "resCostRng": 2,
      "resCostMan": 2,
      "resCostExp": 2,
      "resCostBnk": 2,
      "turnBased": 0,
      "tickRate": 60,
      "turnJumpTicks": 5,
      "turnTime": 24,
      "turnTimeType": 0,
      "autoStart": 0,
      "nonDefaultSettings": [
        "players",
        "anonymity",
        "autoStart"
      ],
      "adminUserId": "abc-dddd-ddd-dd-dddd",
      "fealty": 0,
      "buildWorms": 0,
      "devCostGate": 2,
      "fleetInc": 0,
      "chatId": "dddd-ddd-dd-ddd-ddd"
    },
    "now": 1742092025549,
    "tickFragment": 0,
    "paused": true,
    "started": false,
    "gameOver": false,
    "victoryPoints": {},
    "startTime": 1542091243276,
    "productions": 0,
    "productionRate": 20,
    "productionCounter": 0,
    "tick": 0,
    "admin": 3,
    "name": "2025-03",
    "starsForVictory": 72,
    "starsForFealty": 48,
    "totalStars": 144,
    "turnBased": 0,
    "turnDeadline": 0,
    "tradeCost": 25,
    "tradeScanned": 0,
    "tickRate": 60,
    "fleetSpeed": 0.041666666666666664,
    "players": {
      "1": {
        "uid": 1,
        "alias": "",
        "avatar": 0,
        "race": [
          0,
          0
        ],
        "color": 1,
        "shape": 0,
        "totalStars": 6,
        "totalFleets": 0,
        "totalStrength": 60,
        "totalEconomy": 10,
        "totalIndustry": 5,
        "totalScience": 2,
        "acceptedVassal": 0,
        "offersOfFealty": [],
        "vassals": {},
        "karmaToGive": 0,
        "ready": 0,
        "missedTurns": 0,
        "conceded": 0,
        "ai": 0,
        "regard": 0,
        "tech": {
          "0": {
            "kind": 0,
            "level": 1
          },
          "1": {
            "kind": 1,
            "level": 1
          },
          "2": {
            "kind": 2,
            "level": 1
          },
          "3": {
            "kind": 3,
            "level": 1
          },
          "5": {
            "kind": 5,
            "level": 1
          }
        }
      },
      "2": {
        "uid": 2,
        "alias": "",
        "avatar": 0,
        "race": [
          0,
          0
        ],
        "color": 6,
        "shape": 1,
        "totalStars": 6,
        "totalFleets": 0,
        "totalStrength": 60,
        "totalEconomy": 10,
        "totalIndustry": 5,
        "totalScience": 2,
        "acceptedVassal": 0,
        "offersOfFealty": [],
        "vassals": {},
        "karmaToGive": 0,
        "ready": 0,
        "missedTurns": 0,
        "conceded": 0,
        "ai": 0,
        "regard": 0,
        "tech": {
          "0": {
            "kind": 0,
            "level": 1
          },
          "1": {
            "kind": 1,
            "level": 1
          },
          "2": {
            "kind": 2,
            "level": 1
          },
          "3": {
            "kind": 3,
            "level": 1
          },
          "5": {
            "kind": 5,
            "level": 1
          }
        }
      },
      "3": {
        "cash": 500,
        "war": {
          "1": 3,
          "2": 3,
          "3": 3,
          "4": 3,
          "5": 3,
          "6": 3
        },
        "countdown_to_war": {
          "1": 0,
          "2": 0,
          "3": 0,
          "4": 0,
          "5": 0,
          "6": 0
        },
        "starsAbandoned": 0,
        "ledger": {
          "1": 0,
          "2": 0,
          "3": 0,
          "4": 0,
          "5": 0,
          "6": 0
        },
        "home": 1,
        "researching": 3,
        "researchingNext": 3,
        "uid": 3,
        "alias": "Alrakis Enactment",
        "avatar": 31,
        "race": [
          2,
          13
        ],
        "color": 0,
        "shape": 2,
        "totalStars": 6,
        "totalFleets": 1,
        "totalStrength": 60,
        "totalEconomy": 10,
        "totalIndustry": 5,
        "totalScience": 2,
        "acceptedVassal": 0,
        "offersOfFealty": [],
        "vassals": {},
        "karmaToGive": 0,
        "ready": 0,
        "missedTurns": 0,
        "conceded": 0,
        "ai": 0,
        "regard": 0,
        "tech": {
          "0": {
            "kind": 0,
            "level": 1,
            "research": 0,
            "cost": 144
          },
          "1": {
            "kind": 1,
            "level": 1,
            "research": 0,
            "cost": 160
          },
          "2": {
            "kind": 2,
            "level": 1,
            "research": 0,
            "cost": 144
          },
          "3": {
            "kind": 3,
            "level": 1,
            "research": 0,
            "cost": 128
          },
          "5": {
            "kind": 5,
            "level": 1,
            "research": 0,
            "cost": 144
          }
        }
      },
      "4": {
        "uid": 4,
        "alias": "",
        "avatar": 0,
        "race": [
          0,
          0
        ],
        "color": 3,
        "shape": 3,
        "totalStars": 6,
        "totalFleets": 0,
        "totalStrength": 60,
        "totalEconomy": 10,
        "totalIndustry": 5,
        "totalScience": 2,
        "acceptedVassal": 0,
        "offersOfFealty": [],
        "vassals": {},
        "karmaToGive": 0,
        "ready": 0,
        "missedTurns": 0,
        "conceded": 0,
        "ai": 0,
        "regard": 0,
        "tech": {
          "0": {
            "kind": 0,
            "level": 1
          },
          "1": {
            "kind": 1,
            "level": 1
          },
          "2": {
            "kind": 2,
            "level": 1
          },
          "3": {
            "kind": 3,
            "level": 1
          },
          "5": {
            "kind": 5,
            "level": 1
          }
        }
      },
      "5": {
        "uid": 5,
        "alias": "",
        "avatar": 0,
        "race": [
          0,
          0
        ],
        "color": 4,
        "shape": 4,
        "totalStars": 6,
        "totalFleets": 0,
        "totalStrength": 60,
        "totalEconomy": 10,
        "totalIndustry": 5,
        "totalScience": 2,
        "acceptedVassal": 0,
        "offersOfFealty": [],
        "vassals": {},
        "karmaToGive": 0,
        "ready": 0,
        "missedTurns": 0,
        "conceded": 0,
        "ai": 0,
        "regard": 0,
        "tech": {
          "0": {
            "kind": 0,
            "level": 1
          },
          "1": {
            "kind": 1,
            "level": 1
          },
          "2": {
            "kind": 2,
            "level": 1
          },
          "3": {
            "kind": 3,
            "level": 1
          },
          "5": {
            "kind": 5,
            "level": 1
          }
        }
      },
      "6": {
        "uid": 6,
        "alias": "Alderamin Identity",
        "avatar": 33,
        "race": [
          4,
          12
        ],
        "color": 2,
        "shape": 5,
        "totalStars": 6,
        "totalFleets": 1,
        "totalStrength": 60,
        "totalEconomy": 10,
        "totalIndustry": 5,
        "totalScience": 2,
        "acceptedVassal": 0,
        "offersOfFealty": [],
        "vassals": {},
        "karmaToGive": 0,
        "ready": 0,
        "missedTurns": 0,
        "conceded": 0,
        "ai": 0,
        "regard": 0,
        "tech": {
          "0": {
            "kind": 0,
            "level": 1
          },
          "1": {
            "kind": 1,
            "level": 1
          },
          "2": {
            "kind": 2,
            "level": 1
          },
          "3": {
            "kind": 3,
            "level": 1
          },
          "5": {
            "kind": 5,
            "level": 1
          }
        }
      }
    },
    "stars": {
      "1": {
        "uid": 1,
        "x": 0,
        "y": 0,
        "n": "Alrakis",
        "exp": 0,
        "puid": 3,
        "v": 1,
        "r": 50,
        "nr": 50,
        "yard": 0,
        "e": 10,
        "i": 5,
        "s": 2,
        "ga": 0,
        "st": 9
      },
      "7": {
        "uid": 7,
        "x": 0.2995013139127818,
        "y": -0.05235242888701741,
        "n": "Maasym",
        "exp": 0,
        "puid": 3,
        "v": 1,
        "r": 35,
        "nr": 35,
        "yard": 0,
        "e": 0,
        "i": 0,
        "s": 0,
        "ga": 0,
        "st": 10
      },
      "8": {
        "uid": 8,
        "x": 0.160359848353965,
        "y": -0.2029954259992639,
        "n": "Izar",
        "exp": 0,
        "puid": 3,
        "v": 1,
        "r": 25,
        "nr": 25,
        "yard": 0,
        "e": 0,
        "i": 0,
        "s": 0,
        "ga": 0,
        "st": 10
      },
      "9": {
        "uid": 9,
        "x": -0.0506243006776772,
        "y": -0.5653102516168558,
        "n": "Quo",
        "exp": 0,
        "puid": -1,
        "v": 1,
        "r": 22,
        "nr": 22,
        "yard": 0,
        "e": 0,
        "i": 0,
        "s": 0,
        "ga": 0,
        "st": 0
      },
      "10": {
        "uid": 10,
        "x": -0.5187499737040024,
        "y": -0.24671245071276982,
        "n": "Adhafera",
        "exp": 0,
        "puid": -1,
        "v": 1,
        "r": 10,
        "nr": 10,
        "yard": 0,
        "e": 0,
        "i": 0,
        "s": 0,
        "ga": 0,
        "st": 0
      },
      "11": {
        "uid": 11,
        "x": 0.5092904833475888,
        "y": -0.20061493876106606,
        "n": "Atlas",
        "exp": 0,
        "puid": 3,
        "v": 1,
        "r": 10,
        "nr": 10,
        "yard": 0,
        "e": 0,
        "i": 0,
        "s": 0,
        "ga": 0,
        "st": 10
      },
      "12": {
        "uid": 12,
        "x": -0.05026820476143323,
        "y": -0.7725168307750729,
        "n": "Markab",
        "exp": 0,
        "puid": -1,
        "v": 1,
        "r": 3,
        "nr": 3,
        "yard": 0,
        "e": 0,
        "i": 0,
        "s": 0,
        "ga": 0,
        "st": 0
      },
      "16": {
        "uid": 16,
        "x": 0.9671558540415333,
        "y": -0.5651099869797169,
        "n": "Beid",
        "exp": 0,
        "puid": -1,
        "v": 1,
        "r": 23,
        "nr": 23,
        "yard": 0,
        "e": 0,
        "i": 0,
        "s": 0,
        "ga": 0,
        "st": 0
      },
      "34": {
        "uid": 34,
        "x": 0.7482483797099546,
        "y": 0.3248379551695191,
        "n": "Zosma",
        "exp": 0,
        "puid": -1,
        "v": 1,
        "r": 7,
        "nr": 7,
        "yard": 0,
        "e": 0,
        "i": 0,
        "s": 0,
        "ga": 0,
        "st": 0
      },
      "35": {
        "uid": 35,
        "x": 0.6529353465227437,
        "y": -0.23614081118275188,
        "n": "Kraz",
        "exp": 0,
        "puid": -1,
        "v": 1,
        "r": 3,
        "nr": 3,
        "yard": 0,
        "e": 0,
        "i": 0,
        "s": 0,
        "ga": 0,
        "st": 0
      },
      "38": {
        "uid": 38,
        "x": 0.3801252239047054,
        "y": 0.0412007321853034,
        "n": "Pee",
        "exp": 0,
        "puid": 3,
        "v": 1,
        "r": 30,
        "nr": 30,
        "yard": 0,
        "e": 0,
        "i": 0,
        "s": 0,
        "ga": 0,
        "st": 10
      },
      "49": {
        "uid": 49,
        "x": -0.6141393045856991,
        "y": 0.1859766761779742,
        "n": "Mog",
        "exp": 0,
        "puid": -1,
        "v": 1,
        "r": 44,
        "nr": 44,
        "yard": 0,
        "e": 0,
        "i": 0,
        "s": 0,
        "ga": 0,
        "st": 0
      },
      "107": {
        "uid": 107,
        "x": 1.1307440560407676,
        "y": -0.1778913557776436,
        "n": "Zu",
        "exp": 0,
        "puid": 0,
        "v": 1,
        "r": 30,
        "nr": 30,
        "yard": 0,
        "e": 0,
        "i": 0,
        "s": 0,
        "ga": 0,
        "st": 0
      },
      "127": {
        "uid": 127,
        "x": 0.27836817844076833,
        "y": -0.5362451428477966,
        "n": "Kajam",
        "exp": 0,
        "puid": -1,
        "v": 1,
        "r": 43,
        "nr": 43,
        "yard": 0,
        "e": 0,
        "i": 0,
        "s": 0,
        "ga": 0,
        "st": 0
      },
      "132": {
        "uid": 132,
        "x": -0.2899562202173563,
        "y": -0.0019745532983419523,
        "n": "Algedi",
        "exp": 0,
        "puid": 3,
        "v": 1,
        "r": 15,
        "nr": 15,
        "yard": 0,
        "e": 0,
        "i": 0,
        "s": 0,
        "ga": 0,
        "st": 10
      }
    },
    "fleets": {
      "1": {
        "uid": 1,
        "puid": 3,
        "x": 0,
        "y": 0,
        "lx": 0,
        "ly": 0,
        "exp": 0,
        "speed": 0,
        "st": 1,
        "lsuid": 1,
        "ouid": 1,
        "o": [],
        "l": 0
      }
    }
  }
}
	`), &apiResp)
	if err != nil {
		t.Error(err)
	}
}
