<template>
  <div class="galaxy-map" ref="parent">
    <div class="container" ref="container" />
  </div>
</template>

<script lang="ts">
import {
  Component, Prop, Vue, Watch,
} from 'vue-property-decorator';
import { Network, Node, Edge } from 'vis-network';

import { APIResponse } from '@/types/api';

const neptuneColor = (slot: number) => {
  if (slot === -1) {
    return 'grey';
  }

  return [
    'rgba(0, 0, 255, 1)',
    'rgba(0, 159, 223, 1)',
    'rgba(64, 192, 0, 1)',
    'rgba(255, 192, 0, 1)',
    'rgba(223, 95, 0, 1)',
    'rgba(192, 0, 0, 1)',
    'rgba(192, 0, 192, 1)',
    'rgba(96, 0, 192, 1)',
  ][slot % 8];
};

const groups: { [key: string]: Partial<Node> } = {};

for (let i = -1; i < 8; i += 1) {
  const color = neptuneColor(i);
  groups[`${i}`] = {
    color: {
      background: color,
      border: color,
    },
  };
}

@Component
export default class GalaxyMap extends Vue {
  @Prop() private data!: APIResponse;

  private lastChart: Network|null = null;

  private scale = 1000;

  @Watch('data', { immediate: true })
  private renderMap(canRetry = true) {
    if (!this.$refs.container) {
      if (!canRetry) {
        throw new Error('container still missing after retry');
      }

      console.warn('tried to render map before container existed');
      this.$nextTick(() => this.renderMap(false));
      return;
    }

    if (this.lastChart) {
      this.lastChart.destroy();
      this.lastChart = null;
    }

    const data = Object.values(this.data.scanning_data!.stars).map((star): Node => ({
      id: `star-${star.uid}`,
      label: `${star.n}\n${star.st === undefined ? '?' : star.st}`,
      group: `${star.puid}`,
      fixed: { x: true, y: true },
      x: parseFloat(star.x) * this.scale,
      y: parseFloat(star.y) * this.scale,

      shape: 'dot',
      font: {
        color: 'white',
        strokeColor: 'black',
        strokeWidth: 2,
        size: 32,
      },
      size: 30,
    }));

    const edges: Edge[] = [];

    Object.values(this.data.scanning_data!.fleets).forEach((fleet) => {
      data.push({
        id: `fleet-${fleet.uid}`,
        label: `${fleet.n}\n${fleet.st}`,
        group: `${fleet.puid}`,
        fixed: { x: true, y: true },
        x: parseFloat(fleet.x) * this.scale,
        y: parseFloat(fleet.y) * this.scale,

        shape: 'diamond',
        font: {
          color: 'white',
          strokeColor: 'black',
          strokeWidth: 2,
          size: 16,
        },
        size: 15,
      });

      let start = `fleet-${fleet.uid}`;
      fleet.o.forEach((order, i) => {
        const end = `star-${order[1]}`;

        edges.push({
          from: start,
          to: end,
          color: {
            color: neptuneColor(fleet.puid),
          },
          width: 5 - ((5 / fleet.o.length) * i),
        });

        start = end;
      });
    });

    this.lastChart = new Network(
      this.$refs.container as HTMLElement,
      {
        nodes: data,
        edges,
      },
      {
        physics: {
          enabled: false,
        },
        groups,
      },
    );
    this.hackAroundNetworkChartSizeIssues();
  }

  private hackAroundNetworkChartSizeIssues() {
    this.$nextTick(() => {
      const container = this.$refs.container as HTMLElement;
      const parent = this.$refs.parent as HTMLElement;

      container.style.height = `${parent.clientHeight}px`;
    });
  }

  public mounted() {
    window.addEventListener('resize', () => this.hackAroundNetworkChartSizeIssues);
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.galaxy-map {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  width: 100%;
  height: 100%;
}
</style>
