import Vue from 'vue';
import VueMeta from 'vue-meta';
import App from './App.vue';
import router from './router';

Vue.config.productionTip = false;

Vue.use(VueMeta);

console.log(`
# AlbinoDrought/np-scanner
Repo: https://github.com/AlbinoDrought/np-scanner
Source: ${window.location.origin}/source.tar.gz
`);

new Vue({
  router,
  render: (h) => h(App),
}).$mount('#app');
