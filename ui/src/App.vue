<template>
  <v-app id="inspire">

    <v-app-bar app>
      <img class="mr-3" :src="require('./assets/logo.png')" height="50"/>
      <v-toolbar-title>Wg Gen Web</v-toolbar-title>
    </v-app-bar>

    <v-content>
      <v-container>
        <router-view />
      </v-container>
    </v-content>

    <v-footer app>
      <span>License <a class="pr-1 pl-1" href="http://www.wtfpl.net/" target="_blank">WTFPL</a> &copy; {{ new Date().getFullYear() }} Created with</span><v-icon class="pr-1 pl-1">mdi-heart</v-icon><span>by</span><a class="pr-1 pl-1" href="mailto:vx3r@127-0-0-1.fr">vx3r</a>
    </v-footer>

  </v-app>
</template>

<script>
export default {
  name: 'App',

  data: () => ({
    //
  }),

  created () {
    this.$vuetify.theme.dark = true
  },
};
</script>
