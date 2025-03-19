package types

import "encoding/json"

type Fleet struct {
	UID      int `json:"uid"`
	PlayerID int `json:"puid"`

	CurrentX float64 `json:"x"`
	CurrentY float64 `json:"y"`

	LastX float64 `json:"lx"`
	LastY float64 `json:"ly"`

	XP       float64 `json:"exp"`
	Speed    float64 `json:"speed"`
	Strength int     `json:"st"`

	LastStar    int     `json:"lsuid"`
	CurrentStar int     `json:"ouid"`
	Orders      [][]int `json:"o"`

	// Unknown bool `json:"l"`
}

type PublicStar struct {
	UID      int         `json:"uid"`
	Name     string      `json:"n"`
	PlayerID int         `json:"puid"`
	Visible  json.Number `json:"v"` // `1` for visible, `"0"` for not visible?
	X        float64     `json:"x"`
	Y        float64     `json:"y"`
}

type PrivateStar struct {
	Resources        int     `json:"r"`
	NaturalResources int     `json:"nr"`
	Yard             float64 `json:"yard"`
	Economy          int     `json:"e"`
	Industry         int     `json:"i"`
	Science          int     `json:"s"`
	WarpGate         int     `json:"ga"`
	Strength         int     `json:"st"`
}

func (ps PrivateStar) Useful() bool {
	return ps.Economy > 0 || ps.Industry > 0 || ps.Science > 0 || ps.Resources > 0 || ps.WarpGate > 0 || ps.NaturalResources > 0 || ps.Strength > 0
}

type Star struct {
	PublicStar
	PrivateStar
}

type TechKind int

const TechBanking = TechKind(0)
const TechExperimentation = TechKind(1)
const TechManufacturing = TechKind(2)
const TechRange = TechKind(3)
const TechScan = TechKind(4)
const TechWeapons = TechKind(5)
const TechTerraforming = TechKind(6)

type PublicTechResearchStatus struct {
	Kind  TechKind `json:"kind"`
	Level int      `json:"level"`
}

type PrivateTechResearchStatus struct {
	Research int `json:"research"`
	Cost     int `json:"cost"`
}

func (ptrs PrivateTechResearchStatus) Useful() bool {
	return ptrs.Research > 0 || ptrs.Cost > 0
}

type TechResearchStatus struct {
	PublicTechResearchStatus
	PrivateTechResearchStatus
}

type Tech struct {
	Banking         TechResearchStatus `json:"0"`
	Experimentation TechResearchStatus `json:"1"`
	Manufacturing   TechResearchStatus `json:"2"`
	Range           TechResearchStatus `json:"3"`
	Scan            TechResearchStatus `json:"4"`
	Weapons         TechResearchStatus `json:"5"`
	Terraforming    TechResearchStatus `json:"6"`
}

const (
	ConcededNo       = 0
	ConcededYes      = 1
	ConcededInactive = 2
	ConcededWipedOut = 3
)

type PublicPlayer struct {
	UID           int    `json:"uid"`
	Alias         string `json:"alias"`
	Avatar        int    `json:"avatar"`
	Color         int    `json:"color"`
	Shape         int    `json:"shape"`
	TotalStars    int    `json:"totalStars"`
	TotalFleets   int    `json:"totalFleets"`
	TotalStrength int    `json:"totalStrength"`
	TotalEconomy  int    `json:"totalEconomy"`
	TotalIndustry int    `json:"totalIndustry"`
	TotalScience  int    `json:"totalScience"`
	KarmaToGive   int    `json:"karmaToGive"`
	Ready         int    `json:"ready"`
	MissedTurns   int    `json:"missedTurns"`
	Conceded      int    `json:"conceded"`
	Ai            int    `json:"ai"`
	Regard        int    `json:"regard"`
	Tech          Tech   `json:"tech"`
}

type PrivatePlayer struct {
	Cash            int            `json:"cash"`
	Researching     TechKind       `json:"researching"`
	ResearchingNext TechKind       `json:"researchingNext"`
	War             map[string]int `json:"war"`
	CountdownToWar  map[string]int `json:"countdown_to_war"` // NP4 *sic*, it's still countdown_to_war and not countdownToWar
	StarsAbandoned  int            `json:"starsAbandoned"`
	Home            int            `json:"home"`
}

func (pp PrivatePlayer) Useful() bool {
	return pp.Cash > 0 || pp.Researching != 0 || pp.ResearchingNext != 0 || len(pp.War) > 0 || len(pp.CountdownToWar) > 0 || pp.StarsAbandoned > 0
}

type Player struct {
	PublicPlayer
	PrivatePlayer
}

type ScanningData struct {
	PlayerUID         int               `json:"playerUid"`
	Now               int64             `json:"now"`
	TickFragment      float64           `json:"tickFragment"`
	Paused            bool              `json:"paused"`
	Started           bool              `json:"started"`
	GameOver          bool              `json:"gameOver"`
	StartTime         int64             `json:"startTime"`
	Productions       int               `json:"productions"`
	ProductionRate    int               `json:"productionRate"`
	ProductionCounter int               `json:"productionCounter"`
	Tick              int               `json:"tick"`
	Admin             int               `json:"admin"`
	Name              string            `json:"name"`
	StarsForVictory   int               `json:"starsForVictory"`
	TotalStars        int               `json:"totalStars"`
	TickRate          int               `json:"tickRate"`
	FleetSpeed        float64           `json:"fleetSpeed"`
	Players           map[string]Player `json:"players"`
	Stars             map[string]Star   `json:"stars"`
	Fleets            map[string]Fleet  `json:"fleets"`
}

type APIResponse struct {
	Error        string       `json:"error"`
	ScanningData ScanningData `json:"scanning_data"`
}
