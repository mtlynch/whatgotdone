import Vue from 'vue';
import VueRouter from 'vue-router';

import LandingPage from '@/views/LandingPage.vue';
import LoginPage from '@/views/LoginPage.vue';
import LogoutPage from '@/views/LogoutPage.vue';
import PersonalizedFeedPage from '@/views/PersonalizedFeedPage.vue';
import PrivacyPolicyPage from '@/views/PrivacyPolicyPage.vue';
import RecentEntriesPage from '@/views/RecentEntriesPage.vue';
import ShutdownNoticePage from '@/views/ShutdownNoticePage.vue';
import EditEntryPage from '@/views/EditEntryPage.vue';
import ViewEntryPage from '@/views/ViewEntryPage.vue';
import ViewProjectPage from '@/views/ViewProjectPage.vue';
import EditUserProfilePage from '@/views/EditUserProfilePage.vue';
import UserPreferencesPage from '@/views/UserPreferencesPage.vue';
import UserProfilePage from '@/views/UserProfilePage.vue';
import ExportPage from '@/views/ExportPage.vue';
import MissingPage from '@/views/404.vue';

Vue.use(VueRouter);

const routes = [
  {path: '/about', component: LandingPage},
  {path: '/recent', component: RecentEntriesPage},
  {path: '/feed', component: PersonalizedFeedPage},
  {path: '/entry/edit/:date', component: EditEntryPage, name: 'EditEntry'},
  {path: '/login', component: LoginPage},
  {path: '/logout', component: LogoutPage},
  {path: '/preferences', component: UserPreferencesPage, name: 'Preferences'},
  {path: '/export', component: ExportPage, name: 'Export'},
  {path: '/privacy-policy', component: PrivacyPolicyPage},
  {path: '/shutdown-notice', component: ShutdownNoticePage},
  {
    path: '/:username',
    component: UserProfilePage,
    meta: {
      title: (route) => {
        return `${route.params.username} - What Got Done`;
      },
    },
  },
  {path: '/profile/edit', component: EditUserProfilePage, name: 'EditProfile'},
  {
    path: '/:username/:date',
    component: ViewEntryPage,
    meta: {
      title: (route) => {
        return `${route.params.username}'s What Got Done for the week of ${route.params.date}`;
      },
    },
  },
  {
    path: '/:username/project/:project',
    component: ViewProjectPage,
    meta: {
      title: (route) => {
        return `${route.params.username}'s What Got Done | ${route.params.project}`;
      },
    },
  },
  {path: '/', component: LandingPage},
  {path: '*', component: MissingPage},
];

export default new VueRouter({
  routes,
  mode: 'history',
});
