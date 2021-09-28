<template>
  <div class="war-map">
    <h1>Generating War Map (probably). Refresh if this takes too long</h1>
    <div class="wrapper">
      <galaxy-map
        chartSize="8000px"
        :data="data"
        :selectedStar="selectedStar"
        :selectedFleet="selectedFleet"
        @rendered="saveMap"
      />
    </div>
  </div>
</template>

<script lang="ts">
import {
  Component, Vue, Prop,
} from 'vue-property-decorator';
import GalaxyMap from '@/components/GalaxyMap.vue';
import { APIResponse, Fleet, Star } from '@/types/api';

@Component({
  components: {
    GalaxyMap,
  },
})
export default class MassiveWarMap extends Vue {
  @Prop() private data!: APIResponse;

  @Prop()
  public selectedStar!: Star|null;

  @Prop()
  public selectedFleet!: Fleet|null;

  public saveMap(canvasContainer: HTMLDivElement) {
    // I always thought $nextTick did this same thing, but it doesn't seem to.
    // requestAnimationFrame seems to be the most sensible thing
    // setTimeout(() => {}, 0) also works
    // Accepting PRs for a better solution :)
    requestAnimationFrame(() => {
      try {
        const canvas = canvasContainer.querySelector('canvas');
        if (!canvas) {
          throw new Error('Canvas missing after render?');
        }

        const link = document.createElement('a');
        link.download = `WarMap_${this.data.scanning_data!.name.replace(/[^A-Za-z0-9]/g, '')}_${this.data.scanning_data!.tick}.png`;
        link.href = canvas.toDataURL('image/png');
        link.click();
      } catch (ex) {
        console.error(ex);
      } finally {
        this.$emit('created');
      }
    });
  }
}
</script>

<style scoped lang="scss">
.war-map {
  position: absolute;
  background-color: black;
  top: 0;
  left: 0;
  z-index: 1337;
}
</style>
