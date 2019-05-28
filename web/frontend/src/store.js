import Vue from 'vue';
import Vuex from 'vuex';
import VuexPersistence from 'vuex-persist'

Vue.use(Vuex);

const vuexLocal = new VuexPersistence({
  storage: window.localStorage
})

export default new Vuex.Store({
  state: {
    username: null
  },
  mutations: {
    setUsername(state, username) {
      state.username = username
    },
    clearUsername(state) {
      state.username = null;
    }
  },
  plugins: [vuexLocal.plugin]
});