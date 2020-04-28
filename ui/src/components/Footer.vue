<template>
    <v-container>
        <v-footer app>
            <v-row justify="start" no-gutters>
                <v-col cols="12" lg="6" md="12" sm="12">
                    <div :align="$vuetify.breakpoint.smAndDown ? 'center' : 'left'">
                        <small>Copyright &copy; {{ new Date().getFullYear() }}, Wg Gen Web. </small>
                        <small>This work is licensed under <a class="pr-1 pl-1" href="http://www.wtfpl.net/" target="_blank">WTFPL License.</a></small>
                    </div>
                </v-col>
            </v-row>
            <v-row justify="end" no-gutters>
                <v-col cols="12" lg="6" md="12" sm="12">
                    <div :align="$vuetify.breakpoint.smAndDown ? 'center' : 'right'">
                        <small>Created with</small>
                        <v-icon class="pr-1 pl-1">mdi-heart</v-icon><span>by</span><a class="pr-2 pl-1" href="mailto:vx3r@127-0-0-1.fr">vx3r</a>
                        <a :href="'https://github.com/vx3r/wg-gen-web/commit/' + version"><kbd>Version: {{ version.substring(0,7) }}</kbd></a>
                    </div>
                </v-col>
            </v-row>
        </v-footer>
    </v-container>
</template>

<script>
  import {mapActions, mapGetters} from "vuex";

  export default {
    name: 'Footer',

    data: () => ({

    }),

    computed:{
      ...mapGetters({
        version: 'server/version',
      }),
    },

    mounted() {
      this.versionServer()
    },

    methods: {
      ...mapActions('server', {
        versionServer: 'version',
      }),
    }
  }
</script>
