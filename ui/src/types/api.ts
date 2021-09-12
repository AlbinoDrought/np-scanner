export interface Fleet {
  uid: Int16Array;
  o: number[][];
  n: string;
  puid: number;
  ouid: number;
  w: number;
  x: string;
  y: string;
  st: number;
  lx: string;
  ly: string;
}

export interface PublicStar {
  uid: number;
  n: string;
  puid: number;
  v: string;
  x: string;
  y: string;
}

export interface PrivateStar {
  c: number;
  e: number;
  i: number;
  s: number;
  r: number;
  ga: number;
  nr: number;
  st: number;
}

export type Star = PublicStar & Partial<PrivateStar>;

export interface PublicTechResearchStatus {
  level: number;
  value: number;
}

export interface PrivateTechResearchStatus {
  sv: number;
  research: number;
  bv: number;
  brr: number;
}

export type TechResearchStatus = PublicTechResearchStatus & Partial<PrivateTechResearchStatus>;

export interface Tech {
  scanning: TechResearchStatus;
  propulsion: TechResearchStatus;
  terraforming: TechResearchStatus;
  research: TechResearchStatus;
  weapons: TechResearchStatus;
  banking: TechResearchStatus;
  manufacturing: TechResearchStatus;
}

export interface PublicPlayer {
  total_industry: number;
  regard: number;
  total_science: number;
  uid: number;
  ai: number;
  huid: number;
  total_stars: number;
  total_fleets: number;
  total_strength: number;
  alias: string;
  tech: Tech;
  avatar: number;
  conceded: number;
  ready: number;
  total_economy: number;
  missed_turns: number;
  karma_to_give: number;
}

export interface PrivatePlayer {
  researching: string;
  war: { [key: string]: number };
  stars_abandoned: number;
  cash: number;
  researching_next: string;
  countdown_to_war: { [key: string]: number };
}

export type Player = PublicPlayer & Partial<PrivatePlayer>;

export interface ScanningData {
  fleets: { [key: string]: Fleet };
  fleet_speed: number;
  paused: boolean;
  productions: number;
  tick_fragment: number;
  now: number;
  tick_rate: number;
  production_rate: number;
  stars: { [key: string]: Star };
  stars_for_victory: number;
  game_over: number;
  started: boolean;
  start_time: number;
  total_stars: number;
  production_counter: number;
  trade_scanned: number;
  tick: number;
  trade_cost: number;
  name: string;
  player_uid: number;
  admin: number;
  turn_based: number;
  war: number;
  players: { [key: string]: Player };
  turn_based_time_out: number;
}

export interface APIResponse {
  error?: string;
  scanning_data?: ScanningData;
}
