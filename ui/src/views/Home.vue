<template>
  <div class="home">
    <galaxy-map :data="data" />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import GalaxyMap from '@/components/GalaxyMap.vue';
import { APIResponse } from '@/types/api';

import data from '../../data.jq.json';

@Component({
  components: {
    GalaxyMap,
  },
})
export default class Home extends Vue {
  private data: APIResponse = data;
}
</script>

<style scoped lang="scss">
.home {
  height: 100%;
}
</style>
