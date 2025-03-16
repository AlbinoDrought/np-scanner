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

      <a href="#" @click.prevent="timeTravel = !timeTravel">
        Ghetto Time Travel
      </a>

      <time-travel
        v-if="timeTravel"
        :accessCode="accessCode"
        :match="match"
        @travel="v => $emit('travel', v)"
        @returnToPresent="$emit('returnToPresent')"
      />

      <p
        class="threats"
        :class="{
          'threats--okay': majorThreatCount === 0,
          'threats--danger': majorThreatCount > 0,
        }">
        Threats: {{ privateFleetThreats.length }}
      </p>

      <template v-if="moreInfo">
        <p
          v-for="(threat, i) in privateFleetThreats"
          :key="`threat-${i}`"
          class="threat"
          :class="{
            'threat--minor': threat.battleResults.defenderWins,
            'threat--major': threat.battleResults.attackerWins,
          }"
        >
          <span>
            <strong>
              <a
                href="#"
                @click.prevent="$emit(
                  'selectStar',
                  data.scanning_data.stars[threat.fleetOwner.home || -1],
                )"
                v-text="threat.fleetOwner.alias"
              />
            </strong>
            vs
            <strong>
              <a
                href="#"
                @click.prevent="$emit(
                  'selectStar',
                  data.scanning_data.stars[threat.targetStarOwner.home || -1],
                )"
                v-text="threat.targetStarOwner.alias"
              />
            </strong>
          </span>
          <span>
            Fleet
            <strong>
              <a href="#" @click.prevent="$emit('selectFleet', threat.fleet)">
                {{ threat.fleet.uid }}
              </a>
            </strong>
            headed to
            <strong>
              <a href="#" @click.prevent="$emit('selectStar', threat.targetStar)">
                {{ threat.targetStar.n }}
              </a>
            </strong>
          </span>
          <span>
            Attacker
            <strong>{{ threat.fleet.st }}</strong>
            vs Defender
            <strong>{{ threat.targetStarTrueStrength }}</strong>
          </span>
          <span>({{ threat.travelTime }})</span>
          <span v-if="threat.battleResults.defenderWins" class="result">
            Defender wins with
            <strong>{{ threat.battleResults.defenderShipsRemaining }}</strong>
            ships left :)
          </span>
          <span v-else-if="threat.battleResults.attackerWins" class="result">
            Attacker wins with
            <strong>{{ threat.battleResults.attackerShipsRemaining }}</strong>
            ships left :(
            <br>
            Defender needs
            <strong>{{ threat.battleResults.defenderShipsNeeded }}</strong>
            more ships to win
          </span>
        </p>

        <h2>
          Loaded Players
          <a href="#" @click.prevent="addingPlayer = !addingPlayer">(add new)</a>
        </h2>
        <div v-if="addingPlayer" class="form-wrapper">
          <form class="form" @submit.prevent="tryApiKey">
            <label for="playerKey">Enter API Key:</label>
            <input name="playerKey" type="text" v-model="playerKey">
            <button type="submit">{{ addKeyText }}</button>
          </form>
        </div>
        <p
          v-for="player in privatePlayers"
          :key="`private-player-${player.uid}`"
          class="player player--private"
          :class="{ 'player--wiped-out': playerDead(player) }"
        >
          <span>
            <strong>
              <a
                href="#"
                @click.prevent="$emit('selectStar', data.scanning_data.stars[player.home])"
                v-text="player.alias"
              />
            </strong>
            <span>
              ({{ player.totalEconomy }}
              |
              {{ player.totalIndustry }}
              |
              {{ player.totalScience }})
            </span>
            <span>
              ${{ player.cash }}
            </span>
            <span v-if="player.ai">
              [AI]
            </span>
          </span>
          <span v-if="player.conceded === 3">(completely wiped out)</span>
          <span>Weapons Level {{ player.tech[TechKind.Weapons].level }}</span>
          <span>Total Ships: {{ player.totalStrength }}</span>
          <span>Total Carriers: {{ player.totalFleets }}</span>
          <span>{{ currentResearchText(player) }}</span>
          <span>{{ nextResearchText(player) }}</span>
        </p>

        <h2>Other Players</h2>
        <p
          v-for="player in publicPlayers"
          :key="`public-player-${player.uid}`"
          class="player player--public"
          :class="{ 'player--wiped-out': playerDead(player) }"
        >
          <span>
            <strong>
              <span v-text="player.alias" />
            </strong>
            <span>
              ({{ player.totalEconomy }}
              |
              {{ player.totalIndustry }}
              |
              {{ player.totalScience }})
            </span>
            <span v-if="player.ai">
              [AI]
            </span>
          </span>
          <span v-if="player.conceded === 3">(completely wiped out)</span>
          <span>Weapons Level {{ player.tech[TechKind.Weapons].level }}</span>
          <span>
            Total Ships: {{ player.totalStrength }},
            Visible: {{ visibleAndHiddenStrength[player.uid].visible }},
            Hidden: {{ visibleAndHiddenStrength[player.uid].hidden }}
          </span>
          <span>
            Total Carriers: {{ player.totalFleets }},
            Visible: {{ visibleAndHiddenFleets[player.uid].visible }},
            Hidden: {{ visibleAndHiddenFleets[player.uid].hidden }}
          </span>
        </p>

        <h2>Misc Actions</h2>
        <p>
          <a href="#" @click.prevent="$emit('createMassiveWarMap')">
            Download War Map
          </a>
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
  Fleet,
  isPrivatePlayer,
  Match,
  Player,
  PrivatePlayer,
  PrivateTechResearchStatus,
  PublicPlayer,
  PublicTechResearchStatus,
  Star,
  TechKind,
  TechKinds,
} from '@/types/api';
import {
  distanceBetween,
  guessBattle,
  pointsNeededForTechLevel,
  ticksNeededForResearch,
} from '@/types/algo';

import TimeTravel from './TimeTravel.vue';

const forceGrabTechState = (
  player: PublicPlayer,
  tech: TechKind,
): PublicTechResearchStatus&PrivateTechResearchStatus => ({
  kind: tech,
  level: 0,
  research: 0,
  cost: 0,
  ...player.tech[tech],
});

@Component({
  components: {
    TimeTravel,
  },
})
export default class GameStatus extends Vue {
  @Prop() public accessCode!: string;

  @Prop() private gameNumber!: number;

  @Prop() private data!: APIResponse;

  @Prop() private match!: Match;

  private addingPlayer = false;

  private playerKey = '';

  private addKeyText = 'Add Key';

  private timeTravel = false;

  public moreInfo = false;

  TechKind = TechKind;

  public async tryApiKey() {
    try {
      const resp = await fetch(`/api/matches/${this.gameNumber}/api-key?access_code=${this.accessCode}&api-key=${this.playerKey}`, {
        method: 'POST',
      });
      if (!resp.ok) {
        throw new Error(await resp.text());
      }
      this.addingPlayer = false;
      this.addKeyText = 'Add Key';
      this.$emit('forceRefresh');
    } catch (ex) {
      console.error(ex);
      this.addKeyText = `${this.playerKey} didn't work`;
    }
  }

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

  public playerDead(player: PublicPlayer): boolean {
    // conceded===3 is supposed to work, but sometimes it stays at 1/2
    return player.conceded === 3 || (player.totalStrength === 0 && player.totalStars === 0);
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
      const hidden = player.totalStrength - visible;
      strengths[player.uid] = { visible, hidden };
    });

    return strengths;
  }

  public get visibleAndHiddenFleets() {
    const visibleFleets = new Map<number, number>();

    Object.values(this.data.scanning_data!.fleets).forEach((fleet) => {
      const current = visibleFleets.get(fleet.puid) || 0;
      visibleFleets.set(fleet.puid, current + 1);
    });

    const fleets: { [key: number]: { visible: number, hidden: number } } = {};

    Object.values(this.data.scanning_data!.players).forEach((player) => {
      const visible = visibleFleets.get(player.uid) || 0;
      const hidden = player.totalFleets - visible;
      fleets[player.uid] = { visible, hidden };
    });

    return fleets;
  }

  private get trueStarStrengths() {
    // map of docked star strength + hovering fleet strength
    const trueStarStrengths = new Map<number, number>();

    Object.values(this.data.scanning_data!.stars).forEach((star) => {
      trueStarStrengths.set(star.uid, star.st || 0);
    });

    Object.values(this.data.scanning_data!.fleets).forEach((fleet) => {
      if (!fleet.ouid) {
        return;
      }

      trueStarStrengths.set(fleet.ouid, (trueStarStrengths.get(fleet.ouid) || 0) + fleet.st);
    });

    return trueStarStrengths;
  }

  public get fleetThreats() {
    const fleetThreats: {
      fleet: Fleet,
      order: number[],
      fleetOwnerID: number,
      fleetOwner: Player,
      targetStarID: number,
      targetStar: Star,
      targetStarTrueStrength: number,
      targetStarOwner: Player,
      distance: number,
      fleetSpeed: number,
      travelTicks: number,
      travelTime: string,
      battleResults: {
          attackerWins: boolean;
          defenderWins: boolean;
          attackerShipsRemaining: number;
          defenderShipsRemaining: number;
          lowestTicks: number;
          defenderShipsNeeded: number;
      },
    }[] = [];

    Object.values(this.data.scanning_data!.fleets).forEach((fleet) => {
      const fleetOwnerID = fleet.puid;
      const fleetOwner = this.data.scanning_data!.players[fleetOwnerID]!;

      fleet.o.forEach((order) => {
        const targetStarID = order[1];
        const targetStar = this.data.scanning_data!.stars[`${targetStarID}`];
        if (!targetStar) {
          return;
        }

        if (targetStar.puid !== -1 && targetStar.puid !== fleetOwnerID) {
          const targetStarOwner = this.data.scanning_data!.players[targetStar.puid]!;

          const distance = distanceBetween(targetStar, fleet);

          const { fleetSpeed } = this.data.scanning_data!;
          // todo: reimpl new warpgate logic
          // if (fleet.w) {
          //   fleetSpeed *= 3;
          // }

          const travelTicks = Math.ceil(distance / fleetSpeed);

          const targetStarTrueStrength = this.trueStarStrengths.get(targetStarID)
            || targetStar.st
            || 0;

          fleetThreats.push({
            fleet,
            order,
            fleetOwnerID,
            fleetOwner,
            targetStarID,
            targetStar,
            targetStarTrueStrength,
            targetStarOwner,
            distance,
            fleetSpeed,
            travelTicks,
            travelTime: this.adjustedTicksToTime(travelTicks),
            battleResults: guessBattle(
              fleet.st,
              fleetOwner.tech[TechKind.Weapons].level,
              targetStarTrueStrength,
              targetStarOwner.tech[TechKind.Weapons].level,
            ),
          });
        }
      });
    });

    return fleetThreats;
  }

  public get highestTechOwners() {
    const highestTechOwnerMap = new Map<TechKind, number>();
    const privateTechOwnerMap = new Map<TechKind, number>();

    this.publicPlayers.forEach((player) => {
      for (const tech of TechKinds) {
        const currentHigh = highestTechOwnerMap.get(tech);
        const value = forceGrabTechState(player, tech).level;

        if (value > (currentHigh || 0)) {
          highestTechOwnerMap.set(tech, player.uid);
        }
      }
    });

    this.privatePlayers.forEach((player) => {
      for (const tech of TechKinds) {
        const currentHigh = highestTechOwnerMap.get(tech);
        const currentPrivateHigh = privateTechOwnerMap.get(tech);
        const value = forceGrabTechState(player, tech).level;

        if (value > (currentHigh || 0)) {
          highestTechOwnerMap.set(tech, player.uid);
        }

        if (value > (currentPrivateHigh || 0)) {
          privateTechOwnerMap.set(tech, player.uid);
        }
      }
    });

    return {
      all: highestTechOwnerMap,
      private: privateTechOwnerMap,
    };
  }

  private get privateFleetThreatsAndAttacks() {
    const privatePlayerIDs = new Set<number>();

    this.privatePlayers.forEach((player) => {
      privatePlayerIDs.add(player.uid);
    });

    return {
      threats: this.fleetThreats.filter((threat) => privatePlayerIDs.has(threat.targetStar.puid)),
      attacks: this.fleetThreats.filter((threat) => privatePlayerIDs.has(threat.fleet.puid)),
    };
  }

  public get privateFleetThreats() {
    return this.privateFleetThreatsAndAttacks.threats;
  }

  public get privateFleetAttacks() {
    return this.privateFleetThreatsAndAttacks.attacks;
  }

  public get majorThreatCount() {
    return this.privateFleetThreats.filter((t) => t.battleResults.attackerWins).length;
  }

  public niceTechName(tech: TechKind) {
    if (tech === TechKind.Banking) {
      return 'Banking';
    }
    if (tech === TechKind.Experimentation) {
      return 'Experimentation';
    }
    if (tech === TechKind.Manufacturing) {
      return 'Manufacturing';
    }
    if (tech === TechKind.Range) {
      return 'Hyperspace Range';
    }
    if (tech === TechKind.Scan) {
      return 'Scanning';
    }
    if (tech === TechKind.Terraforming) {
      return 'Terraforming';
    }
    if (tech === TechKind.Weapons) {
      return 'Weapons';
    }
    return `${tech}`;
  }

  private techProgress(
    status: PublicTechResearchStatus&PrivateTechResearchStatus,
    targetLevel: number,
  ) {
    const amountNeeded = pointsNeededForTechLevel(targetLevel);
    return `${status.research}/${amountNeeded}`;
  }

  private techETA(
    player: PublicPlayer,
    status: PublicTechResearchStatus&PrivateTechResearchStatus,
    targetLevel: number,
  ) {
    if (player.totalScience <= 0) {
      return 'stalled (user has no science)';
    }

    const ticksNeeded = ticksNeededForResearch(targetLevel, status.research, player.totalScience);
    return `${ticksNeeded} ticks (${this.adjustedTicksToTime(ticksNeeded)})`;
  }

  private adjustedTicksToTime(unadjustedTicks: number): string {
    const adjustedTicks = Math.max(
      0,
      unadjustedTicks - this.data.scanning_data!.tickFragment,
    );
    return this.ticksToTime(adjustedTicks);
  }

  private ticksToTime(ticks: number): string {
    const ticksAsMinutes = ticks * this.data.scanning_data!.tickRate;

    let minutesRemaining = ticksAsMinutes;

    const days = Math.floor(minutesRemaining / (60 * 24));
    minutesRemaining -= days * (60 * 24);

    const hours = Math.floor(minutesRemaining / 60);
    minutesRemaining -= hours * 60;

    const minutes = Math.ceil(minutesRemaining);
    minutesRemaining = 0;

    return `~${days}d${hours}h${minutes}m`;
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

  private nextResearchText(player: PublicPlayer&PrivatePlayer) {
    let tech = forceGrabTechState(player, player.researchingNext);
    let targetLevel = tech.level + 1;

    if (player.researching === player.researchingNext) {
      // player has the same research queued twice
      // the regular research is for `tech.level + 1`, this is for `tech.level + 2`
      targetLevel += 1;
      // once the player finishes their regular research, progress will reset to 0.
      // mimic that:
      tech = {
        ...tech,
        research: 0,
      };
    }

    return [
      'Next ',
      this.niceTechName(player.researchingNext),
      ' ',
      targetLevel,
      ', currently ',
      this.techProgress(tech, targetLevel),
    ].join('');
  }
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

    .player, .threat {
      display: flex;
      flex-direction: column;
    }

    .player--wiped-out {
      // opacity: 0.5;
      display: none;
    }

    .threats {
      &.threats--danger {
        color: red;
        font-size: larger;
        font-weight: bold;
        animation: blink 1s ease-in-out infinite;
      }
    }

    .threat {
      &.threat--minor .result {
        color: lime;
      }
      &.threat--major .result {
        color: red;
      }
    }
  }
}

.form-wrapper {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  align-items: center;
}

.form {
  display: flex;
  flex-direction: column;
  width: 100%;
  max-width: 500px;

  input { margin: 1em 0; }
}
</style>
