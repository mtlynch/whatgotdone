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
  clearUserState(state) {
    state.username = null;
    state.following = [];
    state.recentFollowingEntries = [];
  },
  setRichTextEditorChoice(state, useRichTextEditor) {
    state.useRichTextEditor = useRichTextEditor;
  },
  setRecent(state, entries) {
    state.recentEntries = entries;
  },
  setRecentFollowing(state, entries) {
    state.recentFollowingEntries = entries;
  },
  addFollowedUser(state, followedUser) {
    if (state.following.includes(followedUser)) {
      return;
    }
    state.following.push(followedUser);
  },
  removeFollowedUser(state, followedUser) {
    state.following = state.following.filter((item) => item !== followedUser);
  },
};

export default new Vuex.Store({
  state: {
    username: null,
    useRichTextEditor: true,
    recentEntries: [],
    recentFollowingEntries: [],
    following: [],
  },
  mutations,
  plugins: [vuexLocal.plugin],
});
