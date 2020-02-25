<template>
  <v-app id="inspire">

    <v-app-bar app>
      <img class="mr-3" :src="require('./assets/logo.png')" height="50" alt="Wg Gen Web"/>
      <v-toolbar-title to="/">Wg Gen Web</v-toolbar-title>

      <v-spacer />

      <v-toolbar-items>
        <v-btn to="/clients">
          Clients
          <v-icon right dark>mdi-account-network-outline</v-icon>
        </v-btn>
        <v-btn to="/server">
          Server
          <v-icon right dark>mdi-vpn</v-icon>
        </v-btn>
      </v-toolbar-items>

    </v-app-bar>

    <v-content>
      <v-container>
        <router-view />
      </v-container>
    </v-content>

    <v-footer app>
      <span>License <a class="pr-1 pl-1" href="http://www.wtfpl.net/" target="_blank">WTFPL</a> &copy; {{ new Date().getFullYear() }} Created with</span><v-icon class="pr-1 pl-1">mdi-heart</v-icon><span>by</span><a class="pr-1 pl-1" href="mailto:vx3r@127-0-0-1.fr">vx3r</a><v-spacer></v-spacer><span>Version: {{ VersionGitCommit }}</span>
    </v-footer>

  </v-app>
</template>

<script>
export default {
  name: 'App',

  data: () => ({
    VersionGitCommit: process.env.VUE_APP_GIT_HASH
  }),

  mounted() {
    console.log("Starting Wg Gen Web version: " + process.env.VUE_APP_GIT_HASH)
  },

  created () {
    this.$vuetify.theme.dark = true
  },
};
</script>
