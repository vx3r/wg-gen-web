import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter);

const routes = [
  {
    path: '/',
    name: 'index',
    component: function () {
      return import(/* webpackChunkName: "Index" */ '../views/Index.vue')
    },
  },
  {
    path: '/clients',
    name: 'clients',
    component: function () {
      return import(/* webpackChunkName: "Clients" */ '../views/Clients.vue')
    },
  },
  {
    path: '/server',
    name: 'server',
    component: function () {
      return import(/* webpackChunkName: "Server" */ '../views/Server.vue')
    },
  }
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
});

export default router
