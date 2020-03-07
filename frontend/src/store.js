import Vue from 'vue';
import Vuex from 'vuex';
import VuexPersistence from 'vuex-persist';

Vue.use(Vuex);

const vuexLocal = new VuexPersistence({
  storage: window.localStorage,
});

export const mutations = {
  setUsername(state, username) {
    state.username = username;
  },
  clearUsername(state) {
    state.username = null;
  },
  setRecent(state, entries) {
    state.recentEntries = entries;
  },
  setFollowing(state, following) {
    state.following = following;
  },
  clearFollowing(state) {
    state.following = [];
  },
};

export default new Vuex.Store({
  state: {
    username: null,
    recentEntries: null,
    following: [],
  },
  mutations,
  plugins: [vuexLocal.plugin],
});
