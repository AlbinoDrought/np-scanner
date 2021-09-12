<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script lang="ts">
import {
  Component, Vue,
} from 'vue-property-decorator';

@Component({
  metaInfo: {
    titleTemplate: '%s | NP Scanner',
  },
})
export default class App extends Vue {
}
</script>

<style lang="scss">
html, body {
  background-color: black;
  color: gold;
  margin: 0;
  min-height: 100vh;

  /* dragging a laggy vis.js network graph can trigger pull-refresh actions, squelch those */
  overscroll-behavior: contain;
}

a {
  color: orangered;
}

#app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  height: 100%;
  width: 100%;
}
</style>
