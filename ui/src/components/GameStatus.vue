<template>
  <div class="game-status">
    <div class="panel">
      <div class="menu">
        <h1>{{ title }}</h1>
        <a href="#" @click.prevent="moreInfo = !moreInfo">
          Toggle
        </a>
      </div>

      <template v-if="moreInfo">
        <p>
          Last Updated: {{ lastUpdatedRelativeString }}
          <br>
          ({{ lastUpdatedString }})
        </p>

        <h2>Loaded Players</h2>
        <p v-for="player in privatePlayers" :key="player.uid" class="player">
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
          <span>{{ currentResearchText(player) }}</span>
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

  public niceTechName(tech: string) {
    return {
      scanning: 'Scanning',
      terraforming: 'Terraforming',
      research: 'Experimentation',
      weapons: 'Weapons',
      banking: 'Banking',
      manufacturing: 'Manufacturing',
    }[tech] || tech;
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
  margin: 1em;

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

      a {
        margin-left: 1em;
      }
    }

    h1, h2 {
      margin: 0;
    }

    .player {
      display: flex;
      flex-direction: column;
    }
  }
}
</style>
