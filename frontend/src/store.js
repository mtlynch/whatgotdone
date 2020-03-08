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
  clearLoginState(state) {
    state.username = null;
    state.following = [];
    state.recentFollowingEntries = [];
  },
  setRecent(state, entries) {
    state.recentEntries = entries;
  },
  setRecentFollowing(state, entries) {
    state.recentFollowingEntries = entries;
  },
  setFollowing(state, following) {
    state.following = following;
  },
};

export default new Vuex.Store({
  state: {
    username: null,
    recentEntries: [],
    recentFollowingEntries: [],
    following: [],
  },
  mutations,
  plugins: [vuexLocal.plugin],
});
