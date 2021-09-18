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

export const techs = [
  'scanning',
  'propulsion',
  'terraforming',
  'research',
  'weapons',
  'banking',
  'manufacturing',
];

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

export const isPrivatePlayer = (
  player: Player,
): player is PublicPlayer&PrivatePlayer => true
  && player.researching !== undefined
  && player.war !== undefined
  && player.stars_abandoned !== undefined
  && player.cash !== undefined
  && player.researching_next !== undefined
  && player.countdown_to_war !== undefined
  && player.researching !== '';

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

export interface PlayerCreds {
  player_uid: number;
  player_alias: string;
  last_poll: string;
  latest_snapshot: number;
}

export interface Match {
  game_number: string;
  name: string;
  last_poll: string;
  player_creds: { [key: string]: PlayerCreds };
}
