<template>
  <div class="not-found">
    <p>Greetings Brother, I'm not sure what you were looking for but it isn't here.</p>
    <router-link to="/">Go Home</router-link>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';

@Component({
  metaInfo: {
    title: 'Not Found',
  },
})
export default class NotFound extends Vue {
}
</script>

<style scoped lang="scss">
.not-found {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  align-items: center;
  width: 100%;
  height: 100%;
}
</style>
