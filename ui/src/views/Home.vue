<template>
  <div class="home">
    <pre class="motd">
 _   _ ____    ____
| \ | |  _ \  / ___|  ___ __ _ _ __  _ __   ___ _ __
|  \| | |_) | \___ \ / __/ _` | '_ \| '_ \ / _ \ '__|
| |\  |  __/   ___) | (_| (_| | | | | | | |  __/ |
|_| \_|_|     |____/ \___\__,_|_| |_|_| |_|\___|_|
    </pre>
    <p v-if="loading">
      Loading (probably)
    </p>
    <p v-else-if="error">
      Error: {{ error }}
    </p>
    <ul v-else class="matches">
      <li
        v-for="match in matches"
        :key="match.game_number"
        class="match"
        :class="{ 'match--finished': match.finished }"
      >
        <router-link :to="{ name: 'Game', params: { gameNumber: match.game_number } }">
          {{ match.name }}
          <span v-if="match.finished">
            (finished)
          </span>
        </router-link>
      </li>
    </ul>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import GalaxyMap from '@/components/GalaxyMap.vue';
import { Match } from '@/types/api';

@Component({
  components: {
    GalaxyMap,
  },
  metaInfo: {
    title: 'Home',
  },
})
export default class Home extends Vue {
  private error: Error|null = null;

  private matches: Match[] = [];

  private loading = false;

  public async mounted() {
    this.loading = true;
    this.matches = [];
    this.error = null;
    try {
      const resp = await fetch('/api/matches');
      if (!resp.ok) {
        throw new Error(await resp.text());
      }
      const json = await resp.json();
      this.matches = json as Match[];
    } catch (ex) {
      console.error(ex);
      this.error = ex;
    }
    this.loading = false;
  }
}
</script>

<style scoped lang="scss">
.home {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  align-items: center;
  width: 100%;
  height: 100%;
}

.motd {
  font-weight: bold;
  font-size: 150%;
}

.form {
  display: flex;
  flex-direction: column;
  width: 100%;
  max-width: 500px;

  input { margin: 1em 0; }
}

.matches {
  list-style: none;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.match {
  margin: 0.5em;
  font-size: larger;
}

.match--finished {
  opacity: 0.5;
}
</style>
