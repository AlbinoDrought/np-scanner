<template>
  <div class="game">
    <galaxy-map v-if="data" :data="data" />
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
import { APIResponse } from '@/types/api';

@Component({
  components: {
    GalaxyMap,
  },
})
export default class Game extends Vue {
  @Prop({ required: true }) public gameNumber!: string;

  private data: APIResponse|null = null;

  private accessCode = '';

  private temporaryAccessCode = '';

  private submitText = 'Submit';

  private requiresAuth = false;

  private error: Error|null = null;

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
    this.setAccessCodeFromRoute();
    await this.loadData();
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
    this.error = null;
    try {
      const resp = await fetch(`/api/matches/${this.gameNumber}/merged-snapshot?access_code=${this.accessCode}`);
      this.requiresAuth = resp.status === 401;
      if (!resp.ok) {
        throw new Error(await resp.text());
      }
      const json = await resp.json();
      this.data = json as APIResponse;
    } catch (ex) {
      console.error(ex);
      this.error = ex;
    }
  }
}
</script>

<style scoped lang="scss">
.game {
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
