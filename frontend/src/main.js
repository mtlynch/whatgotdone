import Vue from 'vue';
import axios from 'axios';
import VueGtag from 'vue-gtag';
import VueAxios from 'vue-axios';
import VueMoment from 'vue-moment';
import BootstrapVue from 'bootstrap-vue';
import {config, library} from '@fortawesome/fontawesome-svg-core';
import {faTwitter} from '@fortawesome/free-brands-svg-icons';
import {FontAwesomeIcon} from '@fortawesome/vue-fontawesome';

import App from './App.vue';
import store from './store.js';
import router from './router.js';

import 'bootswatch/dist/superhero/bootstrap.min.css';
import 'bootstrap-vue/dist/bootstrap-vue.css';

Vue.config.productionTip = false;
Vue.use(VueAxios, axios);
Vue.use(BootstrapVue);

// Need to prevent FontAwesome from auto-adding CSS otherwise it violates
// Content Security Policy.
config.autoAddCss = false;

library.add(faTwitter);
Vue.component('font-awesome-icon', FontAwesomeIcon);

Vue.use(VueMoment);
Vue.prototype.moment = VueMoment;

if (process.env.VUE_APP_GOOGLE_ANALYTICS_ID.length > 1) {
  Vue.use(
    VueGtag,
    {
      config: {id: process.env.VUE_APP_GOOGLE_ANALYTICS_ID},
    },
    router
  );
}

// This callback runs before every route change, including on page load.
router.beforeEach((to, from, next) => {
  if (to.meta.title) {
    document.title = to.meta.title(to);
  } else {
    document.title = 'What Got Done';
  }
  next();
});

new Vue({
  store,
  router,
  render: h => h(App),
}).$mount('#app');
