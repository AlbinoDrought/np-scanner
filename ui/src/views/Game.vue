<template>
  <div class="game">
    <div class="details" v-if="data">
      <galaxy-map
        :accessCode="accessCode"
        :data="data"
        :match="match"
        :selectedStar="selectedStar"
        @selectStar="v => { selectedStar = v; selectedFleet = null }"
        :selectedFleet="selectedFleet"
        @selectFleet="v => { selectedFleet = v; selectedStar = null; }"
        @travel="v => timeSelection = v"
        @returnToPresent="timeSelection = null"
      />
      <game-status
        :accessCode="accessCode"
        :gameNumber="gameNumber"
        :match="match"
        :data="data"
        :selectedStar="selectedStar"
        @selectStar="v => { selectedStar = v; selectedFleet = null }"
        :selectedFleet="selectedFleet"
        @selectFleet="v => { selectedFleet = v; selectedStar = null; }"
        @travel="v => timeSelection = v"
        @returnToPresent="timeSelection = null"
      />
    </div>
    <div class="form-wrapper" v-else-if="requiresAuth">
      <form class="form" @submit.prevent="tryCode">
        <label for="accessCode">Enter Access Code:</label>
        <input v-model="temporaryAccessCode" name="accessCode" placeholder="hunter2">
        <button type="submit">{{ submitText }}</button>
      </form>
    </div>
    <p v-else-if="error">Something broke: {{ error }}</p>
    <p v-else>Loading (probably)</p>
  </div>
</template>

<script lang="ts">
import {
  Component, Vue, Prop, Watch,
} from 'vue-property-decorator';
import GalaxyMap from '@/components/GalaxyMap.vue';
import GameStatus from '@/components/GameStatus.vue';
import {
  APIResponse,
  Fleet,
  Match,
  Star,
} from '@/types/api';

@Component({
  components: {
    GalaxyMap,
    GameStatus,
  },
  metaInfo() {
    const me = this as any;

    const title = me.data && me.data.scanning_data
      ? me.data.scanning_data.name
      : 'Loading';

    return {
      title: `${title} (${me.gameNumber})`,
    };
  },
})
export default class Game extends Vue {
  @Prop({ required: true }) public gameNumber!: string;

  private match: Match|null = null;

  private data: APIResponse|null = null;

  private timeSelection: { [key: string]: string|number }|null = null;

  private accessCode = '';

  private temporaryAccessCode = '';

  private submitText = 'Submit';

  private requiresAuth = false;

  private error: Error|null = null;

  private interval: number|null = null;

  private selectedStar: Star|null = null;

  private selectedFleet: Fleet|null = null;

  public async tryCode() {
    this.submitText = 'Trying code...';
    this.accessCode = this.temporaryAccessCode;

    await this.loadData();

    if (this.requiresAuth) {
      this.submitText = `Nope, not ${this.accessCode}`;
      return;
    }

    this.submitText = 'Submit';
    this.$router.replace({
      name: 'Game',
      params: {
        gameNumber: this.gameNumber,
      },
      query: {
        code: this.temporaryAccessCode,
      },
    });
  }

  public async mounted() {
    this.interval = setInterval(() => this.loadDataNoWipe(), 5 * 60 * 1000);
    this.setAccessCodeFromRoute();
    await this.loadData();
  }

  public beforeDestroy() {
    clearInterval(this.interval!);
  }

  @Watch('$route')
  public setAccessCodeFromRoute() {
    const rawQueryCode = this.$route.query.code;
    if (Array.isArray(rawQueryCode) && rawQueryCode[0]) {
      [this.accessCode] = rawQueryCode;
    } else if (typeof rawQueryCode === 'string') {
      this.accessCode = rawQueryCode;
    } else {
      this.accessCode = '';
    }
  }

  @Watch('gameNumber')
  public async loadData() {
    this.data = null;
    await this.loadDataNoWipe();
  }

  private async fetchMatch() {
    const resp = await fetch(`/api/matches/${this.gameNumber}?access_code=${this.accessCode}`);
    this.requiresAuth = resp.status === 401;
    if (!resp.ok) {
      throw new Error(await resp.text());
    }
    const json = await resp.json();
    this.match = json as Match;
  }

  private async fetchSnapshot() {
    let timeKeyString = '';
    const timeSelectionParams = Object.entries(this.timeSelection || {})
      .map(([player, snapshot]) => `${player}=${snapshot}`);

    if (timeSelectionParams.length > 0) {
      timeKeyString = `&${timeSelectionParams.join('&')}`;
    }

    const resp = await fetch(`/api/matches/${this.gameNumber}/merged-snapshot?access_code=${this.accessCode}${timeKeyString}`);
    this.requiresAuth = resp.status === 401;
    if (!resp.ok) {
      throw new Error(await resp.text());
    }
    const json = await resp.json();
    this.data = json as APIResponse;
  }

  @Watch('timeSelection')
  private async loadDataNoWipe() {
    this.error = null;
    try {
      await Promise.all([
        this.fetchMatch(),
        this.fetchSnapshot(),
      ]);
    } catch (ex) {
      console.error(ex);
      this.error = ex;
    }
  }
}
</script>

<style scoped lang="scss">
.game, .details {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  width: 100%;
  height: 100%;
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
