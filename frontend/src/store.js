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
    // Create a new object instead of simply appending so that Vue's reactivity
    // detection notices the change.
    state.following = this.state.following.concat([followedUser]);
  },
  removeFollowedUser(state, followedUser) {
    this.state.following = this.state.following.filter(
      (item) => item !== followedUser
    );
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
