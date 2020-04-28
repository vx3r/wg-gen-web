<template>
    <v-container>
        <v-app-bar app>
            <img class="mr-3" :src="require('../assets/logo.png')" height="50" alt="Wg Gen Web"/>
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

            <v-menu
                    left
                    bottom
            >
                <template v-slot:activator="{ on }">
                    <v-btn icon v-on="on">
                        <v-icon>mdi-account-circle</v-icon>
                    </v-btn>
                </template>

                <v-card
                        class="mx-auto"
                        max-width="344"
                        outlined
                >
                    <v-list-item three-line>
                        <v-list-item-content>
                            <div class="overline mb-4">connected as</div>
                            <v-list-item-title class="headline mb-1">{{user.name}}</v-list-item-title>
                            <v-list-item-subtitle>Email: {{user.email}}</v-list-item-subtitle>
                            <v-list-item-subtitle>Issuer: {{user.issuer}}</v-list-item-subtitle>
                            <v-list-item-subtitle>Issued at: {{ user.issuedAt | formatDate }}</v-list-item-subtitle>
                        </v-list-item-content>
                    </v-list-item>
                    <v-card-actions>
                        <v-btn small
                                v-on:click="logout()"
                        >
                            logout
                            <v-icon small right dark>mdi-logout</v-icon>
                        </v-btn>
                    </v-card-actions>
                </v-card>
            </v-menu>

        </v-app-bar>
    </v-container>
</template>

<script>
  import {mapActions, mapGetters} from "vuex";

  export default {
    name: 'Header',

    computed:{
      ...mapGetters({
        user: 'auth/user',
      }),
    },

    methods: {
      ...mapActions('auth', {
        logout: 'logout',
      }),
    }
  }
</script>
