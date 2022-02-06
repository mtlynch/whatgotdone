import Vue from 'vue';
import VueRouter from 'vue-router';

import Home from '@/views/Home.vue';
import LoginPage from '@/views/LoginPage.vue';
import LogoutPage from '@/views/LogoutPage.vue';
import PersonalizedFeedPage from '@/views/PersonalizedFeedPage.vue';
import PrivacyPolicyPage from '@/views/PrivacyPolicyPage.vue';
import RecentEntriesPage from '@/views/RecentEntriesPage.vue';
import EditEntry from '@/views/EditEntry.vue';
import ViewEntry from '@/views/ViewEntry.vue';
import ViewProject from '@/views/ViewProject.vue';
import EditUserProfile from '@/views/EditUserProfile.vue';
import UserPreferences from '@/views/UserPreferences.vue';
import UserProfile from '@/views/UserProfile.vue';
import MissingPage from '@/views/404.vue';

Vue.use(VueRouter);

const routes = [
  {path: '/about', component: Home},
  {path: '/recent', component: RecentEntriesPage},
  {path: '/feed', component: PersonalizedFeedPage},
  {path: '/entry/edit/:date', component: EditEntry, name: 'EditEntry'},
  {path: '/login', component: LoginPage},
  {path: '/logout', component: LogoutPage},
  {path: '/preferences', component: UserPreferences, name: 'Preferences'},
  {path: '/privacy-policy', component: PrivacyPolicyPage},
  {
    path: '/:username',
    component: UserProfile,
    meta: {
      title: (route) => {
        return `${route.params.username} - What Got Done`;
      },
    },
  },
  {path: '/profile/edit', component: EditUserProfile, name: 'EditProfile'},
  {
    path: '/:username/:date',
    component: ViewEntry,
    meta: {
      title: (route) => {
        return `${route.params.username}'s What Got Done for the week of ${route.params.date}`;
      },
    },
  },
  {
    path: '/:username/project/:project',
    component: ViewProject,
    meta: {
      title: (route) => {
        return `${route.params.username}'s What Got Done | ${route.params.project}`;
      },
    },
  },
  {path: '/', component: Home},
  {path: '*', component: MissingPage},
];

export default new VueRouter({
  routes,
  mode: 'history',
});
