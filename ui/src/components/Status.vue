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
            <template v-slot:item.connected="{ item }">
              <v-icon left v-if="item.connected" color="success">mdi-lan-connect</v-icon>
              <v-icon left v-else>mdi-lan-disconnect</v-icon>
            </template>
            <template v-slot:item.receivedBytes="{ item }">
              {{ humanFileSize(item.receivedBytes) }}
            </template>
            <template v-slot:item.transmittedBytes="{ item }">
              {{ humanFileSize(item.transmittedBytes) }}
            </template>
            <template v-slot:item.allowedIPs="{ item }">
              <v-chip
                  v-for="(ip, i) in item.allowedIPs"
                  :key="i"
                  color="indigo"
                  text-color="white"
              >
                <v-icon left>mdi-ip-network</v-icon>
                {{ ip }}
              </v-chip>
            </template>
            <template v-slot:item.lastHandshake="{ item }">
              <v-row>
                <p>{{ item.lastHandshake | formatDate }} ({{ item.lastHandshakeRelative }})</p>
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

      humanFileSize(bytes, si=false, dp=1) {
        const thresh = si ? 1000 : 1024;

        if (Math.abs(bytes) < thresh) {
          return bytes + ' B';
        }

        const units = si
            ? ['kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
            : ['KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
        let u = -1;
        const r = 10**dp;

        do {
          bytes /= thresh;
          ++u;
        } while (Math.round(Math.abs(bytes) * r) / r >= thresh && u < units.length - 1);


        return bytes.toFixed(dp) + ' ' + units[u];
      }
    }
  };
</script>
