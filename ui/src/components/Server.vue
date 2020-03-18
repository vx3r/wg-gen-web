<template>
    <v-container v-if="server">
        <v-row>
            <v-col cols="12">
                <v-card dark>
                    <v-list-item>
                        <v-list-item-content>
                            <v-list-item-title class="headline">Server's interface configuration</v-list-item-title>
                        </v-list-item-content>
                    </v-list-item>
                    <div class="d-flex flex-no-wrap justify-space-between">
                        <v-col cols="12">
                            <v-text-field
                                    v-model="server.publicKey"
                                    label="Public key"
                                    disabled
                            />
                            <v-text-field
                                    v-model="server.listenPort"
                                    type="number"
                                    :rules="[
                          v => !!v || 'Listen port is required',
                        ]"
                                    label="Listen port"
                                    required
                            />
                            <v-combobox
                                    v-model="server.address"
                                    chips
                                    hint="Write IPv4 or IPv6 CIDR and hit enter"
                                    label="Server interface addresses"
                                    multiple
                                    dark
                            >
                                <template v-slot:selection="{ attrs, item, select, selected }">
                                    <v-chip
                                            v-bind="attrs"
                                            :input-value="selected"
                                            close
                                            @click="select"
                                            @click:close="server.address.splice(server.address.indexOf(item), 1)"
                                    >
                                        <strong>{{ item }}</strong>&nbsp;
                                    </v-chip>
                                </template>
                            </v-combobox>
                        </v-col>
                    </div>
                </v-card>
            </v-col>
            <v-col cols="12">
                <v-card dark>
                    <v-list-item>
                        <v-list-item-content>
                            <v-list-item-title class="headline">Client's global configuration</v-list-item-title>
                        </v-list-item-content>
                    </v-list-item>
                    <div class="d-flex flex-no-wrap justify-space-between">
                        <v-col cols="12">
                            <v-text-field
                                    v-model="server.endpoint"
                                    label="Public endpoint for clients to connect to"
                                    :rules="[
                          v => !!v || 'Public endpoint for clients to connect to is required',
                        ]"
                                    required
                            />
                            <v-combobox
                                    v-model="server.dns"
                                    chips
                                    hint="Write IPv4 or IPv6 address and hit enter"
                                    label="DNS servers for clients"
                                    multiple
                                    dark
                            >
                                <template v-slot:selection="{ attrs, item, select, selected }">
                                    <v-chip
                                            v-bind="attrs"
                                            :input-value="selected"
                                            close
                                            @click="select"
                                            @click:close="server.dns.splice(server.dns.indexOf(item), 1)"
                                    >
                                        <strong>{{ item }}</strong>&nbsp;
                                    </v-chip>
                                </template>
                            </v-combobox>
                            <v-text-field
                                    type="number"
                                    v-model="server.mtu"
                                    label="Define global MTU"
                                    hint="Leave at 0 and let wg-quick take care of MTU"
                            />
                            <v-text-field
                                    type="number"
                                    v-model="server.persistentKeepalive"
                                    label="Persistent keepalive"
                                    hint="Leave at 0 if you dont want to specify persistent keepalive"
                            />
                        </v-col>
                    </div>
                </v-card>
            </v-col>
        </v-row>
        <v-row>
            <v-col cols="12">
                <v-card dark>
                    <v-list-item>
                        <v-list-item-content>
                            <v-list-item-title class="headline">Interface configuration hooks</v-list-item-title>
                        </v-list-item-content>
                    </v-list-item>
                    <div class="d-flex flex-no-wrap justify-space-between">
                        <v-col cols="12">
                            <v-text-field
                                    v-model="server.preUp"
                                    label="PreUp: script snippets which will be executed by bash before setting up the interface"
                            />
                            <v-text-field
                                    v-model="server.postUp"
                                    label="PostUp: script snippets which will be executed by bash after setting up the interface"
                            />
                            <v-text-field
                                    v-model="server.preDown"
                                    label="PreDown: script snippets which will be executed by bash before setting down the interface"
                            />
                            <v-text-field
                                    v-model="server.postDown "
                                    label="PostDown : script snippets which will be executed by bash after setting down the interface"
                            />
                        </v-col>
                    </div>
                </v-card>
            </v-col>
        </v-row>
        <v-row>
            <v-divider dark/>
            <v-btn
                    class="ma-2"
                    color="success"
                    :href="`${apiBaseUrl}/server/config`"
            >
                Download server configuration
                <v-icon right dark>mdi-cloud-download-outline</v-icon>
            </v-btn>
            <v-spacer></v-spacer>
            <v-btn
                    class="ma-2"
                    color="warning"
                    @click="updateServer"
            >
                Update server configuration
                <v-icon right dark>mdi-update</v-icon>
            </v-btn>
            <v-divider dark/>
        </v-row>
        <Notification v-bind:notification="notification"/>
    </v-container>
</template>
<script>
  import {API_BASE_URL, ApiService} from "../services/ApiService";
  import Notification from '../components/Notification'

  export default {
    name: 'Server',

    components: {
      Notification
    },

    data: () => ({
      api: null,
      server: null,
      apiBaseUrl: API_BASE_URL,
      notification: {
        show: false,
        color: '',
        text: '',
      },
    }),

    mounted () {
      this.api = new ApiService();
      this.getServer()
    },

    methods: {
      getServer() {
        this.api.get('/server').then((res) => {
          this.server = res;
        }).catch((e) => {
          this.notify('error', e.response.status + ' ' + e.response.statusText);
        });
      },
      updateServer () {
        // convert int values
        this.server.listenPort = parseInt(this.server.listenPort, 10);
        this.server.persistentKeepalive = parseInt(this.server.persistentKeepalive, 10);
        this.server.mtu = parseInt(this.server.mtu, 10);

        // check server addresses
        if (this.server.address.length < 1) {
          this.notify('error', 'Please provide at least one valid CIDR address for server interface');
          return;
        }
        for (let i = 0; i < this.server.address.length; i++){
          if (this.$isCidr(this.server.address[i]) === 0) {
            this.notify('error', `Invalid CIDR detected, please correct ${this.server.address[i]} before submitting`);
            return
          }
        }

        // check DNS correct
        for (let i = 0; i < this.server.dns.length; i++){
          if (this.$isCidr(this.server.dns[i] + '/32') === 0) {
            this.notify('error', `Invalid IP detected, please correct ${this.server.dns[i]} before submitting`);
            return
          }
        }

        this.api.patch('/server', this.server).then((res) => {
          this.notify('success', "Server successfully updated");
          this.server = res;
        }).catch((e) => {
          this.notify('error', e.response.status + ' ' + e.response.statusText);
        });
      },
      notify(color, msg) {
        this.notification.show = true;
        this.notification.color = color;
        this.notification.text = msg;
      }
    }
  };
</script>
