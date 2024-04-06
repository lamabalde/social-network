import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import ProfileView from '../views/ProfileView.vue';
import PageNotFound from '../views/PageNotFound.vue';
import LoginView from '../views/LoginView.vue';
import RegisterView from '../views/RegisterView.vue';
import checkAuth from '../middleware/auth.js';
import GroupView from '../views/GroupView.vue'
import EventView from "../views/EventView.vue";
import SearchView from "../views/SearchView.vue";
import Logout from "../js_modules/Logout.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/:catchAll(.*)*",
      name: "PageNotFound",
      component: PageNotFound,
    },
    {
      path: "/",
      name: "home",
      component: HomeView,
      /*meta: {
        middleware: auth,
      }*/
    },
    {
      path: "/profile/:id",
      name: "profile",
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: ProfileView,
    },
    {
      path: "/login",
      name: "login",
      component: LoginView,
    },
    {
      path: "/register",
      name: "register",
      component: RegisterView,
    },
    {
      path: "/logout",
      name: "logout",
      component: Logout,
    },
    {
      path: "/groups/:id",
      name: "groups",
      component: GroupView,
    },
    {
      path: "/event/:id",
      name: "event",
      component: EventView,
    },
    {
      path: "/search/:id?",
      name: "search",
      component: SearchView,
      
    },
  ],
});

router.beforeEach(async function (to, from) {
  const authenticated = await checkAuth()
  if (to.path!="/login" && to.path!="/register") {
    if (!authenticated) return '/login'
  } else {
    if (authenticated) return '/'
  }
})



export default router


