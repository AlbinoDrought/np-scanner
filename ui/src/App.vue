<template>
  <div id="app">
    <router-view />
  </div>
</template>

<style lang="scss">
html, body {
  background-color: black;
  color: gold;
  margin: 0;
  min-height: 100vh;

  /* dragging a laggy vis.js network graph can trigger pull-refresh actions, squelch those */
  overscroll-behavior: contain;
}
#app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  height: 100%;
  width: 100%;
}
</style>
