import Vue from 'vue'
import VueRouter from 'vue-router'
import store from "../store";

Vue.use(VueRouter);

const routes = [
  {
    path: '/clients',
    name: 'clients',
    component: function () {
      return import(/* webpackChunkName: "Clients" */ '../views/Clients.vue')
    },
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/server',
    name: 'server',
    component: function () {
      return import(/* webpackChunkName: "Server" */ '../views/Server.vue')
    },
    meta: {
      requiresAuth: true
    }
  }
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
});

router.beforeEach((to, from, next) => {
  if(to.matched.some(record => record.meta.requiresAuth)) {
    if (store.getters["auth/isAuthenticated"]) {
      next()
      return
    }
    next('/')
  } else {
    next()
  }
})

export default router
