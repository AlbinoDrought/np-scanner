<template>
  <div>
    <div v-for="player in players" :key="player.player_uid" class="time">
      {{ player.player_alias }}

      <select v-model="selection[player.player_uid]">
        <option
          v-for="snapshot in niceSnapshots(player.player_uid)"
          :key="snapshot.value"
          :value="snapshot.value"
          v-text="snapshot.text"
        />
      </select>
    </div>

    <div class="buttons">
      <button type="button" @click.prevent="travel">Travel</button>
      <button type="button" @click.prevent="doReturn">Return</button>
    </div>
  </div>
</template>

<script lang="ts">
import {
  Component,
  Prop,
  Vue,
  Watch,
} from 'vue-property-decorator';
import { Match, PlayerCreds } from '@/types/api';

@Component({})
export default class TimeTravel extends Vue {
  @Prop() public accessCode!: string;

  @Prop() public match!: Match;

  private snapshots: { [key: string]: string[] } = {};

  private selection: { [key: string]: string } = {};

  private get players(): PlayerCreds[] {
    return Object.values(this.match.player_creds);
  }

  @Watch('match', { immediate: true })
  public async loadSnapshots() {
    if (!this.match) {
      return;
    }

    this.players.forEach((player) => {
      if (!this.selection[player.player_uid]) {
        this.$set(this.selection, player.player_uid, player.latest_snapshot);
      }
    });

    await this.players.map(
      (creds) => this.loadSnapshotsForPlayer(creds.player_uid),
    );
  }

  private async loadSnapshotsForPlayer(player: number) {
    const resp = await fetch(`/api/matches/${this.match.game_number}/player-snapshots/${player}?access_code=${this.accessCode}&limit=500`);
    if (!resp.ok) {
      throw new Error(await resp.text());
    }
    const json = await resp.json();
    this.$set(this.snapshots, player, json as string[]);
  }

  private niceSnapshots(player: number) {
    const snapshots = this.snapshots[player] || [];

    return snapshots.map((snapshot) => ({
      text: (new Date(snapshot)).toISOString(),
      value: snapshot,
    }));
  }

  private travel() {
    this.$emit('travel', this.selection);
  }

  private doReturn() {
    this.$emit('returnToPresent');
  }
}
</script>

<style scoped lang="scss">
.buttons {
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
  button {
    margin: 0.5em 0;
  }
  button+button {
    margin-left: 0.5em;
  }
}

.time {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;

  select {
    margin: 0.5em;
    margin-right: 0;
  }
}
</style>
