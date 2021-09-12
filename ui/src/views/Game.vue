<template>
  <div class="game">
    <galaxy-map v-if="data" :data="data" />
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

  private error: Error|null = null;

  @Watch('gameNumber', { immediate: true })
  public async loadData() {
    this.data = null;
    this.error = null;
    try {
      const resp = await fetch(`/api/matches/${this.gameNumber}/merged-snapshot`);
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
</style>
