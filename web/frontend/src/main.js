import Vue from 'vue'
import axios from 'axios'
import VueAxios from 'vue-axios'
import VueRouter from 'vue-router'
import BootstrapVue from 'bootstrap-vue'
import App from './App.vue'
import Login from "./components/Login.vue";
import Submit from "./components/Submit.vue";
import ViewEntry from "./components/ViewEntry.vue";
import MissingPage from "./components/404.vue";

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

Vue.config.productionTip = false
Vue.use(VueAxios, axios)
Vue.use(BootstrapVue)
Vue.use(VueRouter)

const routes = [
  { path: '/:username/:date', component: ViewEntry },
  { path: '/submit', component: Submit },
  { path: '/login', component: Login },
  { path: '*', component: MissingPage },
]

const router = new VueRouter({
  routes,
  mode: 'history'
})

new Vue({
  router,
  render: h => h(App),
}).$mount('#app')