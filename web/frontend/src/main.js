import Vue from 'vue'
import axios from 'axios'
import VueAnalytics from 'vue-analytics'
import VueAxios from 'vue-axios'
import VueMoment from 'vue-moment'
import VueRouter from 'vue-router'
import BootstrapVue from 'bootstrap-vue'
import App from './App.vue'
import Home from "./components/Home.vue";
import Login from "./components/Login.vue";
import Logout from "./components/Logout.vue";
import ProUpgrade from "./components/ProUpgrade.vue";
import Submit from "./components/Submit.vue";
import ViewEntry from "./components/ViewEntry.vue";
import MissingPage from "./components/404.vue";

import 'bootswatch/dist/superhero/bootstrap.min.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

Vue.config.productionTip = false
Vue.use(VueAxios, axios)
Vue.use(BootstrapVue)
Vue.use(VueRouter)

Vue.use(VueMoment)
Vue.prototype.moment = VueMoment

const routes = [
  { path: '/pro', component: ProUpgrade },
  { path: '/submit', component: Submit },
  { path: '/login', component: Login },
  { path: '/logout', component: Logout },
  { path: '/:username', component: ViewEntry },
  { path: '/:username/:date', component: ViewEntry },
  { path: '/', component: Home },
  { path: '*', component: MissingPage },
]

const router = new VueRouter({
  routes,
  mode: 'history'
})

Vue.use(VueAnalytics, {
  id: process.env.VUE_APP_GOOGLE_ANALYTICS_ID,
  router
})

new Vue({
  router,
  render: h => h(App),
}).$mount('#app')