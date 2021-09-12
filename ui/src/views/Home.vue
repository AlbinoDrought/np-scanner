<template>
  <div class="home">
    <pre class="motd">
 _   _ ____    ____
| \ | |  _ \  / ___|  ___ __ _ _ __  _ __   ___ _ __
|  \| | |_) | \___ \ / __/ _` | '_ \| '_ \ / _ \ '__|
| |\  |  __/   ___) | (_| (_| | | | | | | |  __/ |
|_| \_|_|     |____/ \___\__,_|_| |_|_| |_|\___|_|
    </pre>
    <form class="form" @submit.prevent="submit">
      <label for="gameNumber">Enter Game Number or URL:</label>
      <input v-model="gameNumber" name="gameNumber" placeholder="https://np.ironhelmet.com/game/643054096969696">
      <button type="submit">{{ submitText }}</button>
    </form>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import GalaxyMap from '@/components/GalaxyMap.vue';
import { APIResponse } from '@/types/api';

@Component({
  components: {
    GalaxyMap,
  },
})
export default class Home extends Vue {
  private gameNumber = '';

  private submitText = 'Submit';

  public submit() {
    this.submitText = 'Submit';

    const pieces = this.gameNumber.split('/');
    const gameNumber = pieces[pieces.length - 1];
    const parsedGameNumber = parseInt(gameNumber, 10);

    if (Number.isNaN(parsedGameNumber)) {
      this.submitText = 'That number is broken yo, try again';
      return;
    }

    this.$router.push({
      name: 'Game',
      params: {
        gameNumber,
      },
    });
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
</style>
