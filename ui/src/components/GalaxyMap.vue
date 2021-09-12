<template>
  <div class="galaxy-map">
    <div class="container" ref="container" />
  </div>
</template>

<script lang="ts">
import {
  Component, Prop, Vue, Watch,
} from 'vue-property-decorator';
import { Network, Node } from 'vis-network';

import { APIResponse } from '@/types/api';

const groups = {
  '-1': { color: { background: 'grey', border: 'black' } },
  0: { color: { background: 'rgba(0, 0, 255, 1)', border: 'black' } },
  1: { color: { background: 'rgba(0, 159, 223, 1)', border: 'black' } },
  2: { color: { background: 'rgba(64, 192, 0, 1)', border: 'black' } },
  3: { color: { background: 'rgba(255, 192, 0, 1)', border: 'black' } },
  4: { color: { background: 'rgba(223, 95, 0, 1)', border: 'black' } },
  5: { color: { background: 'rgba(192, 0, 0, 1)', border: 'black' } },
  6: { color: { background: 'rgba(192, 0, 192, 1)', border: 'black' } },
  7: { color: { background: 'rgba(96, 0, 192, 1)', border: 'black' } },
};

@Component
export default class GalaxyMap extends Vue {
  @Prop() private data!: APIResponse;

  private lastChart: Network|null = null;

  private scale = 500;

  @Watch('data', { immediate: true })
  private renderMap(canRetry = true) {
    console.log(this.data);

    if (!this.$refs.container) {
      if (!canRetry) {
        throw new Error('container still missing after retry');
      }

      console.warn('tried to render map before container existed');
      this.$nextTick(() => this.renderMap(false));
      return;
    }

    const data = Object.values(this.data.scanning_data!.stars).map((star): Node => ({
      id: star.uid,
      label: star.n,
      group: `${star.puid}`,
      fixed: { x: true, y: true },
      x: parseFloat(star.x) * this.scale,
      y: parseFloat(star.y) * this.scale,

      shape: 'dot',
      font: {
        color: 'white',
        strokeColor: 'black',
        strokeWidth: 2,
      },
    }));

    this.lastChart = new Network(
      this.$refs.container as HTMLElement,
      {
        nodes: data,
        edges: [],
      },
      {
        physics: {
          enabled: false,
        },
        height: '100%',
        groups,
      },
    );
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.galaxy-map {
  width: 100%;
  height: 100%;
  .container {
    width: 1000px;
    height: 1000px
  }
}
</style>
