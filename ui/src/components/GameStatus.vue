<template>
  <div class="game-status">
    <div class="panel">
      <div class="menu">
        <h1>
          <a :href="`https://np.ironhelmet.com/game/${gameNumber}`" target="_blank">
            {{ title }}
          </a>
        </h1>
        <a href="#" @click.prevent="moreInfo = !moreInfo">
          Toggle Info
        </a>
      </div>

      <p>
        Last Updated: {{ lastUpdatedRelativeString }}
        <br>
        ({{ lastUpdatedString }})
      </p>

      <template v-if="moreInfo">
        <h2>Loaded Players</h2>
        <p v-for="player in privatePlayers" :key="player.uid" class="player player--private">
          <span>
            <strong>{{ player.alias }}</strong>
            <span>
              ({{ player.total_economy }}
              |
              {{ player.total_industry }}
              |
              {{ player.total_science }})
            </span>
          </span>
          <span>Total Ships: {{ player.total_strength }}</span>
          <span>Total Carriers: {{ player.total_fleets }}</span>
          <span>{{ currentResearchText(player) }}</span>
        </p>

        <h2>Other Players</h2>
        <p v-for="player in publicPlayers" :key="player.uid" class="player player--public">
          <span>
            <strong>{{ player.alias }}</strong>
            <span>
              ({{ player.total_economy }}
              |
              {{ player.total_industry }}
              |
              {{ player.total_science }})
            </span>
          </span>
          <span>
            Total Ships: {{ player.total_strength }},
            Visible: {{ visibleAndHiddenStrength[player.uid].visible }},
            Hidden: {{ visibleAndHiddenStrength[player.uid].hidden }}
          </span>
          <span>
            Total Carriers: {{ player.total_fleets }},
            Visible: {{ visibleAndHiddenFleets[player.uid].visible }},
            Hidden: {{ visibleAndHiddenFleets[player.uid].hidden }}
          </span>
        </p>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
import {
  Component, Prop, Vue,
} from 'vue-property-decorator';

import {
  APIResponse,
  isPrivatePlayer,
  PrivatePlayer,
  PrivateTechResearchStatus,
  PublicPlayer,
  PublicTechResearchStatus,
} from '@/types/api';

const forceGrabTechState = (
  player: PublicPlayer&PrivatePlayer,
  tech: string,
): PublicTechResearchStatus&PrivateTechResearchStatus => {
  const techState = (player.tech as any)[tech];
  return techState as PublicTechResearchStatus&PrivateTechResearchStatus;
};

@Component({})
export default class GameStatus extends Vue {
  @Prop() private gameNumber!: number;

  @Prop() private data!: APIResponse;

  public moreInfo = false;

  public get title() {
    return this.data.scanning_data!.name;
  }

  private get lastUpdatedObject() {
    return new Date(this.data.scanning_data!.now);
  }

  public get lastUpdatedString() {
    return this.lastUpdatedObject.toISOString();
  }

  public get lastUpdatedRelativeString() {
    const msSinceUpdate = (new Date()).getTime() - this.lastUpdatedObject.getTime();
    const minutesSinceUpdate = Math.floor(msSinceUpdate / (60 * 1000));

    return minutesSinceUpdate === 1
      ? '1 minute ago'
      : `${minutesSinceUpdate} minutes ago`;
  }

  public get privatePlayers(): Array<PublicPlayer&PrivatePlayer> {
    const privatePlayers: Array<PublicPlayer&PrivatePlayer> = [];

    Object.values(this.data.scanning_data!.players).forEach((player) => {
      if (isPrivatePlayer(player)) {
        privatePlayers.push(player);
      }
    });

    return privatePlayers;
  }

  public get publicPlayers(): Array<PublicPlayer> {
    const publicPlayers: Array<PublicPlayer> = [];

    const privatePlayerIDs = new Set<number>();

    this.privatePlayers.forEach((player) => {
      privatePlayerIDs.add(player.uid);
    });

    return Object.values(this.data.scanning_data!.players)
      .filter((p) => !privatePlayerIDs.has(p.uid));
  }

  private get visibleStrength(): Map<number, number> {
    const visibleStrength = new Map<number, number>();

    Object.values(this.data.scanning_data!.players).forEach((player) => {
      visibleStrength.set(player.uid, 0);
    });

    Object.values(this.data.scanning_data!.fleets).forEach((fleet) => {
      let currentStrength = visibleStrength.get(fleet.puid) || 0;
      currentStrength += fleet.st;
      visibleStrength.set(fleet.puid, currentStrength);
    });

    Object.values(this.data.scanning_data!.stars).forEach((star) => {
      if (!star.st) {
        return;
      }

      let currentStrength = visibleStrength.get(star.puid) || 0;
      currentStrength += star.st;
      visibleStrength.set(star.puid, currentStrength);
    });

    return visibleStrength;
  }

  public get visibleAndHiddenStrength() {
    const { visibleStrength } = this;
    const strengths: { [key: number]: { visible: number, hidden: number } } = {};

    Object.values(this.data.scanning_data!.players).forEach((player) => {
      const visible = visibleStrength.get(player.uid) || 0;
      const hidden = player.total_strength - visible;
      strengths[player.uid] = { visible, hidden };
    });

    return strengths;
  }

  private get visibleAndHiddenFleets() {
    const visibleFleets = new Map<number, number>();

    Object.values(this.data.scanning_data!.fleets).forEach((fleet) => {
      const current = visibleFleets.get(fleet.puid) || 0;
      visibleFleets.set(fleet.puid, current + 1);
    });

    const fleets: { [key: number]: { visible: number, hidden: number } } = {};

    Object.values(this.data.scanning_data!.players).forEach((player) => {
      const visible = visibleFleets.get(player.uid) || 0;
      const hidden = player.total_fleets - visible;
      fleets[player.uid] = { visible, hidden };
    });

    return fleets;
  }

  public niceTechName(tech: string) {
    return ({
      scanning: 'Scanning',
      terraforming: 'Terraforming',
      propulsion: 'Hyperspace Range',
      research: 'Experimentation',
      weapons: 'Weapons',
      banking: 'Banking',
      manufacturing: 'Manufacturing',
    } as any)[tech] || tech;
  }

  private techProgress(
    status: PublicTechResearchStatus&PrivateTechResearchStatus,
    targetLevel: number,
  ) {
    const amountNeeded = 144 * (targetLevel - 1);

    return `${status.research}/${amountNeeded}`;
  }

  private techETA(
    player: PublicPlayer,
    status: PublicTechResearchStatus&PrivateTechResearchStatus,
    targetLevel: number,
  ) {
    const amountNeeded = 144 * (targetLevel - 1);

    const ticksNeeded = Math.ceil((amountNeeded - status.research) / player.total_science);
    const adjustedTicksNeeded = ticksNeeded - this.data.scanning_data!.tick_fragment;
    const ticksAsMinutes = Math.ceil(adjustedTicksNeeded * this.data.scanning_data!.tick_rate);

    let minutesRemaining = ticksAsMinutes;
    let days = 0;
    let hours = 0;
    let minutes = 0;

    while (minutesRemaining >= (60 * 24)) {
      minutesRemaining -= 60 * 24;
      days += 1;
    }

    while (minutesRemaining >= 60) {
      minutesRemaining -= 60;
      hours += 1;
    }

    minutes = minutesRemaining;
    minutesRemaining = 0;

    return `${ticksNeeded} ticks (~${days}d${hours}h${minutes}m)`;
  }

  private currentResearchText(player: PublicPlayer&PrivatePlayer) {
    const tech = forceGrabTechState(player, player.researching);
    const targetLevel = tech.level + 1;

    return [
      'Researching ',
      this.niceTechName(player.researching),
      ' ',
      targetLevel,
      ', ',
      this.techProgress(tech, targetLevel),
      ', ',
      this.techETA(player, tech, targetLevel),
    ].join('');
  }

  /*
  private nextResearchText(player: PublicPlayer&PrivatePlayer) {
    const tech = forceGrabTechState(player, player.researching_next);
    const targetLevel = tech.level + (player.researching === player.researching_next ? 2 : 1);

    return [
      'Next ',
      this.niceTechName(player.researching_next),
      ' ',
      targetLevel,
      ', ',
      this.techProgress(tech, targetLevel),
      ', ',
      this.techETA(player, tech, targetLevel),
    ].join('');
  }
  */
}
</script>

<style scoped lang="scss">
.game-status {
  position: fixed;
  top: 0;
  left: 0;
  margin: 1vh;
  max-height: 98vh;
  overflow-y: auto;

  .panel {
    padding: 1em;
    background-color: black;
    border: 1px solid orange;
    display: flex;
    flex-direction: column;

    .menu {
      display: flex;
      flex-direction: row;
      align-items: center;
      justify-content: space-between;

      &>a {
        margin-left: 1em;
      }
    }

    h1, h2 {
      margin: 0;
    }
    h2 {
      margin-top: 1em;
    }

    .player {
      display: flex;
      flex-direction: column;
    }
  }
}
</style>
