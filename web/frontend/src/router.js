import Vue from 'vue';
import VueRouter from 'vue-router';

import Home from "./components/Home.vue";
import Login from "./components/Login.vue";
import Logout from "./components/Logout.vue";
import PrivacyPolicy from "./components/PrivacyPolicy.vue";
import ProPitch from "./components/ProPitch.vue";
import ProUpgrade from "./components/ProUpgrade.vue";
import Recent from "./components/Recent.vue";
import Submit from "./components/Submit.vue";
import ViewEntry from "./components/ViewEntry.vue";
import MissingPage from "./components/404.vue";

Vue.use(VueRouter)

const routes = [
  { path: '/recent', component: Recent },
  { path: '/pro', component: ProPitch },
  { path: '/upgrade', component: ProUpgrade },
  { path: '/submit', component: Submit },
  { path: '/login', component: Login },
  { path: '/logout', component: Logout },
  { path: '/privacy-policy', component: PrivacyPolicy },
  { path: '/:username', component: ViewEntry },
  { path: '/:username/:date', component: ViewEntry },
  { path: '/', component: Home },
  { path: '*', component: MissingPage },
]

export default new VueRouter({
  routes,
  mode: 'history'
})