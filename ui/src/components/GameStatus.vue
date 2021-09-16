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
                  data.scanning_data.stars[threat.fleetOwner.huid],
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
                  data.scanning_data.stars[threat.targetStarOwner.huid],
                )"
                v-text="threat.targetStarOwner.alias"
              />
            </strong>
          </span>
          <span>
            Fleet
            <strong>
              <a href="#" @click.prevent="$emit('selectFleet', threat.fleet)">
                {{ threat.fleet.n }}
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
          <!--
          <span>({{ threat.travelTime }})</span>
          -->
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

        <h2>Loaded Players</h2>
        <p
          v-for="player in privatePlayers"
          :key="`private-player-${player.uid}`"
          class="player player--private"
        >
          <span>
            <strong>
              <a
                href="#"
                @click.prevent="$emit('selectStar', data.scanning_data.stars[player.huid])"
                v-text="player.alias"
              />
            </strong>
            <span>
              ({{ player.total_economy }}
              |
              {{ player.total_industry }}
              |
              {{ player.total_science }})
            </span>
          </span>
          <span>Weapons Level {{ player.tech.weapons.level }}</span>
          <span>Total Ships: {{ player.total_strength }}</span>
          <span>Total Carriers: {{ player.total_fleets }}</span>
          <span>{{ currentResearchText(player) }}</span>
        </p>

        <h2>Other Players</h2>
        <p
          v-for="player in publicPlayers"
          :key="`public-player-${player.uid}`"
          class="player player--public"
        >
          <span>
            <strong>
              <a
                href="#"
                @click.prevent="$emit('selectStar', data.scanning_data.stars[player.huid])"
                v-text="player.alias"
              />
            </strong>
            <span>
              ({{ player.total_economy }}
              |
              {{ player.total_industry }}
              |
              {{ player.total_science }})
            </span>
          </span>
          <span>Weapons Level {{ player.tech.weapons.level }}</span>
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
  Fleet,
  isPrivatePlayer,
  Player,
  PrivatePlayer,
  PrivateTechResearchStatus,
  PublicPlayer,
  PublicTechResearchStatus,
  Star,
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

  public get visibleAndHiddenFleets() {
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
      /*
      distance: number,
      distanceLY: number,
      fleetSpeed: number,
      travelTicks: number,
      travelTime: string,
      */
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

          /*
          const distanceX = Math.abs(parseFloat(targetStar.x) - parseFloat(fleet.x));
          const distanceY = Math.abs(parseFloat(targetStar.y) - parseFloat(fleet.y));
          const distance = Math.sqrt(
            (distanceX * distanceX) * (distanceY * distanceY),
          );
          const distanceLY = distance * 8;

          let fleetSpeed = this.data.scanning_data!.fleet_speed;
          if (fleet.w) {
            fleetSpeed *= 3;
          }

          const travelTicks = distanceLY / fleetSpeed; // todo: this is wrong
          console.log(
            targetStar.x,
            fleet.x,
            Math.abs(parseFloat(targetStar.x) - parseFloat(fleet.x)),
          );
          console.log(
            targetStar.y,
            fleet.y,
            Math.abs(parseFloat(targetStar.y) - parseFloat(fleet.y)),
          );
          console.log(distance, distanceLY, fleetSpeed, travelTicks);
          */

          fleetThreats.push({
            fleet,
            order,
            fleetOwnerID,
            fleetOwner,
            targetStarID,
            targetStar,
            targetStarTrueStrength: this.trueStarStrengths.get(targetStarID) || targetStar.st || 0,
            targetStarOwner,
            /*
            distance,
            distanceLY,
            fleetSpeed,
            travelTicks,
            travelTime: this.adjustedTicksToTime(travelTicks),
            */
            battleResults: this.guessBattle(
              fleet.st,
              fleetOwner.tech.weapons.level,
              this.trueStarStrengths.get(targetStarID) || 0,
              targetStarOwner.tech.weapons.level,
            ),
          });
        }
      });
    });

    return fleetThreats;
  }

  public get privateFleetThreats() {
    const privatePlayerIDs = new Set<number>();

    this.privatePlayers.forEach((player) => {
      privatePlayerIDs.add(player.uid);
    });

    return this.fleetThreats.filter((threat) => privatePlayerIDs.has(threat.targetStar.puid));
  }

  public get majorThreatCount() {
    return this.privateFleetThreats.filter((t) => t.battleResults.attackerWins).length;
  }

  public get minorThreatCount() {
    return this.privateFleetThreats.filter((t) => t.battleResults.defenderWins).length;
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
    if (player.total_science <= 0) {
      return 'stalled (user has no science)';
    }

    const amountNeeded = 144 * (targetLevel - 1);
    const ticksNeeded = Math.ceil((amountNeeded - status.research) / player.total_science);
    return `${ticksNeeded} ticks (${this.adjustedTicksToTime(ticksNeeded)})`;
  }

  private adjustedTicksToTime(unadjustedTicks: number): string {
    const adjustedTicks = Math.max(0, unadjustedTicks - this.data.scanning_data!.tick_fragment);
    return this.ticksToTime(adjustedTicks);
  }

  private ticksToTime(ticks: number): string {
    const ticksAsMinutes = Math.ceil(ticks * this.data.scanning_data!.tick_rate);

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

  private guessBattle(
    attackerShips: number,
    attackerWeaponsLevel: number,
    defenderShips: number,
    defenderWeaponsLevel: number,
  ) {
    const defenderWeaponsLevelWithBonus = defenderWeaponsLevel + 1; // defenders bonus

    const ticksToKillAttacker = attackerShips / defenderWeaponsLevelWithBonus; // 14.75
    const ticksToKillDefender = defenderShips / attackerWeaponsLevel; // 0

    const lowestTicks = Math.ceil(Math.min(ticksToKillAttacker, ticksToKillDefender)); // 0

    const attackerShipsRemaining = Math.max(
      0,
      attackerShips - (lowestTicks * defenderWeaponsLevelWithBonus), // 59
    );
    let defenderShipsRemaining = Math.max(
      0,
      defenderShips - (lowestTicks * attackerWeaponsLevel), // 0
    );

    if (attackerShipsRemaining === 0 && defenderShipsRemaining === 0) {
      // defender wins
      defenderShipsRemaining += attackerWeaponsLevel;
    }

    const attackerWins = attackerShipsRemaining > 0;
    const defenderWins = defenderShipsRemaining > 0;

    if (attackerWins === defenderWins) {
      throw new Error('encountered draw but this should be impossible');
    }

    // the rounding could be incorrect, not 100% sure
    const defenderShipsNeeded = Math.floor(Math.max(
      0,
      (ticksToKillAttacker * attackerWeaponsLevel) - defenderShips,
    ));

    return {
      attackerWins,
      defenderWins,

      attackerShipsRemaining,
      defenderShipsRemaining,

      lowestTicks,
      defenderShipsNeeded,
    };
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

    .player, .threat {
      display: flex;
      flex-direction: column;
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
</style>
