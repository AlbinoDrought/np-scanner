package types

type Fleet struct {
	UID         int     `json:"uid"`
	Unknown     int     `json:"l"`
	Orders      [][]int `json:"o"`
	Name        string  `json:"n"`
	PlayerID    int     `json:"puid"`
	CurrentStar int     `json:"ouid"`
	WarpSpeed   int     `json:"w"`
	CurrentX    string  `json:"x"`
	CurrentY    string  `json:"y"`
	Strength    int     `json:"st"`
	LastX       string  `json:"lx"`
	LastY       string  `json:"ly"`
}

type PublicStar struct {
	UID      int    `json:"uid"`
	Name     string `json:"n"`
	PlayerID int    `json:"puid"`
	Visible  string `json:"v"`
	X        string `json:"x"`
	Y        string `json:"y"`
}

type PrivateStar struct {
	ShipsPerTick     float64 `json:"c"`
	Economy          int     `json:"e"`
	Industry         int     `json:"i"`
	Science          int     `json:"s"`
	Resources        int     `json:"r"`
	WarpGate         int     `json:"ga"`
	NaturalResources int     `json:"nr"`
	Strength         int     `json:"st"`
}

type Star struct {
	PublicStar
	PrivateStar
}

type PublicTechResearchStatus struct {
	Level int     `json:"level"`
	Value float64 `json:"value"`
}

type PrivateTechResearchStatus struct {
	Sv               float64 `json:"sv"`
	Research         int     `json:"research"`
	Bv               float64 `json:"bv"`
	CostPerTechLevel int     `json:"brr"`
}

type TechResearchStatus struct {
	PublicTechResearchStatus
	PrivateTechResearchStatus
}

type Tech struct {
	Scanning      TechResearchStatus `json:"scanning"`
	Propulsion    TechResearchStatus `json:"propulsion"`
	Terraforming  TechResearchStatus `json:"terraforming"`
	Research      TechResearchStatus `json:"research"`
	Weapons       TechResearchStatus `json:"weapons"`
	Banking       TechResearchStatus `json:"banking"`
	Manufacturing TechResearchStatus `json:"manufacturing"`
}

type PublicPlayer struct {
	TotalIndustry int    `json:"total_industry"`
	Regard        int    `json:"regard"`
	TotalScience  int    `json:"total_science"`
	UID           int    `json:"uid"`
	Ai            int    `json:"ai"`
	Huid          int    `json:"huid"`
	TotalStars    int    `json:"total_stars"`
	TotalFleets   int    `json:"total_fleets"`
	TotalStrength int    `json:"total_strength"`
	Alias         string `json:"alias"`
	Tech          Tech   `json:"tech"`
	Avatar        int    `json:"avatar"`
	Conceded      int    `json:"conceded"`
	Ready         int    `json:"ready"`
	TotalEconomy  int    `json:"total_economy"`
	MissedTurns   int    `json:"missed_turns"`
	KarmaToGive   int    `json:"karma_to_give"`
}

type PrivatePlayer struct {
	Researching     string         `json:"researching"`
	War             map[string]int `json:"war"`
	StarsAbandoned  int            `json:"stars_abandoned"`
	Cash            int            `json:"cash"`
	ResearchingNext string         `json:"researching_next"`
	CountdownToWar  map[string]int `json:"countdown_to_war"`
}

type Player struct {
	PublicPlayer
	PrivatePlayer
}

type ScanningData struct {
	Fleets            map[string]Fleet  `json:"fleets"`
	FleetSpeed        float64           `json:"fleet_speed"`
	Paused            bool              `json:"paused"`
	Productions       int               `json:"productions"`
	TickFragment      float64           `json:"tick_fragment"`
	Now               int64             `json:"now"`
	TickRate          int               `json:"tick_rate"`
	ProductionRate    int               `json:"production_rate"`
	Stars             map[string]Star   `json:"stars"`
	StarsForVictory   int               `json:"stars_for_victory"`
	GameOver          int               `json:"game_over"`
	Started           bool              `json:"started"`
	StartTime         int64             `json:"start_time"`
	TotalStars        int               `json:"total_stars"`
	ProductionCounter int               `json:"production_counter"`
	TradeScanned      int               `json:"trade_scanned"`
	Tick              int               `json:"tick"`
	TradeCost         int               `json:"trade_cost"`
	Name              string            `json:"name"`
	PlayerUID         int               `json:"player_uid"`
	Admin             int               `json:"admin"`
	TurnBased         int               `json:"turn_based"`
	War               int               `json:"war"`
	Players           map[string]Player `json:"players"`
	TurnBasedTimeOut  int               `json:"turn_based_time_out"`
}

type APIResponse struct {
	Error        string       `json:"error"`
	ScanningData ScanningData `json:"scanning_data"`
}
