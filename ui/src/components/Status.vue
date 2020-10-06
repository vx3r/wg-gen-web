<template>
    <v-container>
      <v-row v-if="dataLoaded">
        <v-col cols="12">
          <v-card>
            <v-card-title>
              WireGuard Interface Status: {{ interface.name }}
            </v-card-title>
            <v-list-item>
              <v-list-item-content>
                <v-list-item-subtitle>Public Key: {{ interface.publicKey }}</v-list-item-subtitle>
                <v-list-item-subtitle>Listening Port: {{ interface.listenPort }}</v-list-item-subtitle>
                <v-list-item-subtitle>Device Type: {{ interface.type }}</v-list-item-subtitle>
                <v-list-item-subtitle>Number of Peers: {{ interface.numPeers }}</v-list-item-subtitle>
              </v-list-item-content>
            </v-list-item>
          </v-card>
        </v-col>
      </v-row>
      <v-row v-if="dataLoaded">
        <v-col cols="12">
          <v-card>
            <v-card-title>
              WireGuard Client Status
              <v-spacer></v-spacer>
              <v-text-field
                  v-model="search"
                  append-icon="mdi-magnify"
                  label="Search"
                  single-line
                  hide-details
              ></v-text-field>
              <v-spacer></v-spacer>
              <v-btn
                  color="success"
                  @click="reload"
              >
                Reload
                <v-icon right dark>mdi-reload</v-icon>
              </v-btn>
            </v-card-title>
            <v-data-table
              :headers="headers"
              :items="clients"
              :search="search"
          >
            <template v-slot:item.address="{ item }">
              <v-chip
                  v-for="(ip, i) in item.address"
                  :key="i"
                  color="indigo"
                  text-color="white"
              >
                <v-icon left>mdi-ip-network</v-icon>
                {{ ip }}
              </v-chip>
            </template>
            <template v-slot:item.tags="{ item }">
              <v-chip
                  v-for="(tag, i) in item.tags"
                  :key="i"
                  color="blue-grey"
                  text-color="white"
              >
                <v-icon left>mdi-tag</v-icon>
                {{ tag }}
              </v-chip>
            </template>
            <template v-slot:item.created="{ item }">
              <v-row>
                <p>At {{ item.created | formatDate }} by {{ item.createdBy }}</p>
              </v-row>
            </template>
            <template v-slot:item.updated="{ item }">
              <v-row>
                <p>At {{ item.updated | formatDate }} by {{ item.updatedBy }}</p>
              </v-row>
            </template>
            <template v-slot:item.action="{ item }">
              <v-row>
                <v-icon
                    class="pr-1 pl-1"
                    @click.stop="startUpdate(item)"
                >
                  mdi-square-edit-outline
                </v-icon>
                <v-icon
                    class="pr-1 pl-1"
                    @click.stop="forceFileDownload(item)"
                >
                  mdi-cloud-download-outline
                </v-icon>
                <v-icon
                    class="pr-1 pl-1"
                    @click.stop="email(item)"
                >
                  mdi-email-send-outline
                </v-icon>
                <v-icon
                    class="pr-1 pl-1"
                    @click="remove(item)"
                >
                  mdi-trash-can-outline
                </v-icon>
                <v-switch
                    dark
                    class="pr-1 pl-1"
                    color="success"
                    v-model="item.enable"
                    v-on:change="update(item)"
                />
              </v-row>
            </template>

          </v-data-table>
          </v-card>
        </v-col>
      </v-row>
      <v-row v-else>
        <v-col cols="12">
          <v-card>
            <v-card-title>
              No stats available...
            </v-card-title>
            <v-card-text>{{ error }}</v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
</template>
<script>
  import { mapActions, mapGetters } from 'vuex'

  export default {
    name: 'Status',

    data: () => ({
      search: '',
      headers: [
        { text: 'Connected', value: 'connected', },
        { text: 'Name', value: 'name', },
        { text: 'Endpoint', value: 'endpoint', },
        { text: 'IP addresses', value: 'allowedIPs', sortable: false, },
        { text: 'Received Bytes', value: 'receivedBytes', },
        { text: 'Transmitted Bytes', value: 'transmittedBytes', },
        { text: 'Last Handshake', value: 'lastHandshake',} ,
      ],
    }),

    computed:{
      ...mapGetters({
        interface: 'status/interfaceStatus',
        clients: 'status/clientStatus',
        error: 'status/error',
      }),
      dataLoaded: function () {
        return this.interface != null && this.interface.name !== "";
      }
    },

    mounted () {
      this.readStatus()
    },

    methods: {
      ...mapActions('status', {
        readStatus: 'read',
      }),

      reload() {
        this.readStatus()
      },
    }
  };
</script>
