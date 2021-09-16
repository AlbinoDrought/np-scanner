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

const colorBase = [
  'rgba(0, 0, 255, 1)',
  'rgba(0, 159, 223, 1)',
  'rgba(64, 192, 0, 1)',
  'rgba(255, 192, 0, 1)',
  'rgba(223, 95, 0, 1)',
  'rgba(192, 0, 0, 1)',
  'rgba(192, 0, 192, 1)',
  'rgba(96, 0, 192, 1)',
];

const colorBaseExtended = [
  'rgba(179, 191, 255, 1)',
  'rgba(179, 255, 255, 1)',
  'rgba(179, 255, 179, 1)',
  'rgba(255, 255, 179, 1)',
  'rgba(255, 204, 153, 1)',
  'rgba(255, 153, 153, 1)',
  'rgba(255, 153, 204, 1)',
  'rgba(204, 153, 255, 1)',
];

const neptuneColor = (slot: number) => {
  if (slot === -1) {
    return 'grey';
  }

  if (slot < 56) {
    return colorBase[slot % colorBase.length];
  }

  return colorBaseExtended[slot % colorBaseExtended.length];
};

const shapes = [
  'dot',
  'square',
  'hexagon',
  'triangle',
  'triangleDown',
  'diamond',
  'star',
  'dot', // should be pill
];

const neptuneShape = (slot: number) => {
  if (slot === -1) {
    return 'dot';
  }

  return shapes[Math.floor(slot / shapes.length) % shapes.length];
};

const groups: { [key: string]: Partial<Node> } = {};

for (let i = -1; i < 64; i += 1) {
  const color = neptuneColor(i);
  const shape = neptuneShape(i);
  groups[`${i}-stars`] = {
    color: {
      background: color,
      border: 'black',
    },
    shape,
    fixed: { x: true, y: true },
    font: {
      color: 'white',
      strokeColor: 'black',
      strokeWidth: 2,
      size: 32,
    },
    size: 30,
  };
  groups[`${i}-fleets`] = {
    color: {
      background: color,
      border: 'black',
    },
    shape,
    fixed: { x: true, y: true },
    font: {
      color: 'white',
      strokeColor: 'black',
      strokeWidth: 2,
      size: 16,
    },
    borderWidth: 1,
    size: 15,
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

    const data: Node[] = [];
    const edges: Edge[] = [];

    const fleetPowerDockedAtStars = new Map<number, number>();

    Object.values(this.data.scanning_data!.fleets).forEach((fleet) => {
      data.push({
        id: `fleet-${fleet.uid}`,
        label: `${fleet.n}\n${fleet.ouid ? '' : fleet.st}`,
        group: `${fleet.puid}-fleets`,
        fixed: { x: true, y: true },
        x: parseFloat(fleet.x) * this.scale,
        y: parseFloat(fleet.y) * this.scale,
      });

      if (fleet.ouid) {
        let currentDockedPower = fleetPowerDockedAtStars.get(fleet.ouid) || 0;
        currentDockedPower += fleet.st;
        fleetPowerDockedAtStars.set(fleet.ouid, currentDockedPower);
      }

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

    Object.values(this.data.scanning_data!.stars).forEach((star) => {
      const fleetPowerDockedAtStar = fleetPowerDockedAtStars.get(star.uid) || 0;
      let altText = '?';

      if (star.v === '1') {
        const powerAtStar = (star.st || 0) + fleetPowerDockedAtStar;
        let powerLine = `${powerAtStar}`;
        if ((star.c || 0) > 0) {
          const newPowerPerTick = (star.c || 0.0).toFixed(2);
          powerLine += ` (+${newPowerPerTick}/t)`;
        }
        const econ = star.e || 0;
        const industry = star.i || 0;
        const science = star.s || 0;
        altText = [
          powerLine,
          `(${econ} | ${industry} | ${science})`,
        ].join('\n');
      }
      data.push({
        id: `star-${star.uid}`,
        label: `${star.n}\n${altText}`,
        group: `${star.puid}-stars`,
        x: parseFloat(star.x) * this.scale,
        y: parseFloat(star.y) * this.scale,
      });
    });

    this.lastChart = new Network(
      this.$refs.container as HTMLElement,
      {
        nodes: data.reverse(),
        edges,
      },
      {
        physics: {
          enabled: false,
        },
        interaction: {
          dragNodes: false,
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
