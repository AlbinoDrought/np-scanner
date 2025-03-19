export interface Fleet {
  uid: number;
  puid: number;

  x: number;
  y: number;

  lx: number;
  ly: number;

  exp: number;
  speed: number;
  st: number;

  lsuid: number;
  ouid: number;
  o: number[][];

  l: unknown;
}

export interface PublicStar {
  uid: number;
  n: string;
  puid: number;
  v: number;
  x: number;
  y: number;

  exp: unknown;
}

export interface PrivateStar {
  r: number;
  nr: number;
  yard: number;
  e: number;
  i: number;
  s: number;
  ga: number;
  st: number;
}

export type Star = PublicStar & Partial<PrivateStar>;

export enum TechKind {
  Banking = 0,
  Experimentation = 1,
  Manufacturing = 2,
  Range = 3,
  Scan = 4,
  Weapons = 5,
  Terraforming = 6,
}

export const TechKinds = [
  TechKind.Banking,
  TechKind.Experimentation,
  TechKind.Manufacturing,
  TechKind.Range,
  TechKind.Scan,
  TechKind.Weapons,
  TechKind.Terraforming,
];

export interface PublicTechResearchStatus {
  kind: TechKind;
  level: number;
}

export interface PrivateTechResearchStatus {
  research: number;
  cost: number;
}

export type TechResearchStatus = PublicTechResearchStatus & Partial<PrivateTechResearchStatus>;

export interface Tech {
  [TechKind.Banking]: TechResearchStatus;
  [TechKind.Experimentation]: TechResearchStatus;
  [TechKind.Manufacturing]: TechResearchStatus;
  [TechKind.Range]: TechResearchStatus;
  [TechKind.Scan]?: TechResearchStatus;
  [TechKind.Weapons]: TechResearchStatus;
  [TechKind.Terraforming]?: TechResearchStatus;
}

export interface PublicPlayer {
  uid: number;
  alias: string;
  avatar: number;
  race: unknown;
  color: number;
  shape: number;
  totalStars: number;
  totalFleets: number;
  totalStrength: number;
  totalEconomy: number;
  totalIndustry: number;
  totalScience: number;
  acceptedVassal: unknown;
  offersOfFealty: unknown;
  vassals: unknown;
  karmaToGive: number;
  ready: number;
  missedTurns: number;
  conceded: number;
  ai: number;
  regard: number;
  tech: Tech;
}

export interface PrivatePlayer {
  cash: number;
  researching: TechKind;
  researchingNext: TechKind;
  war: { [key: string]: number };
  // NP4 *sic*, it's still countdown_to_war and not countdownToWar
  countdown_to_war: { [key: string]: number };
  starsAbandoned: number;
  ledger: unknown;
  home: number;
}

export type Player = PublicPlayer & Partial<PrivatePlayer>;

export const isPrivatePlayer = (
  player: Player,
): player is PublicPlayer&PrivatePlayer => true
  && player.researching !== undefined
  && !!player.war
  && player.starsAbandoned !== undefined
  && player.cash !== undefined
  && player.researchingNext !== undefined
  && !!player.countdown_to_war
  && player.researching !== undefined
  && !!player.home;

export interface ScanningData {
  playerUid: number;
  now: number;
  tickFragment: number;
  paused: boolean;
  started: boolean;
  gameOver: boolean;
  startTime: number;
  productions: number;
  productionRate: number;
  productionCounter: number;
  tick: number;
  admin: number;
  name: string;
  starsForVictory: number;
  totalStars: number;
  tickRate: number;
  fleetSpeed: number;
  players: { [key: string]: Player };
  stars: { [key: string]: Star };
  fleets: { [key: string]: Fleet };

  config: unknown;
  victoryPoints: unknown;
  starsForFealty: unknown;
  turnBased: unknown;
  turnDeadline: unknown;
  tradeCost: unknown;
  tradeScanned: unknown;
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
  finished: boolean;
  last_poll: string;
  player_creds: { [key: string]: PlayerCreds };
}
