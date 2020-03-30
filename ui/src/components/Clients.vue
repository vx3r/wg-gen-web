<template>
    <v-container>
        <v-row>
            <v-col cols="12">
                <v-card dark>
                    <v-list-item>
                        <v-list-item-content>
                            <v-list-item-title class="headline">Clients</v-list-item-title>
                        </v-list-item-content>
                        <v-btn
                                color="success"
                                @click="startAddClient"
                        >
                            Add new client
                            <v-icon right dark>mdi-account-multiple-plus-outline</v-icon>
                        </v-btn>
                    </v-list-item>
                    <v-row>
                        <v-col
                                v-for="(client, i) in clients"
                                :key="i"
                                sm12 lg6
                        >
                            <v-card
                                    :color="client.enable ? '#1F7087' : 'warning'"
                                    class="mx-auto"
                                    raised
                                    shaped
                            >
                                <v-list-item>
                                    <v-list-item-content>
                                        <v-list-item-title class="headline">{{ client.name }}</v-list-item-title>
                                        <v-list-item-subtitle>{{ client.email }}</v-list-item-subtitle>
                                        <v-list-item-subtitle>Created: {{ client.created | formatDate }}</v-list-item-subtitle>
                                        <v-list-item-subtitle>Updated: {{ client.updated | formatDate }}</v-list-item-subtitle>
                                    </v-list-item-content>

                                    <v-list-item-avatar
                                            tile
                                            size="150"
                                    >
                                        <v-img :src="`${apiBaseUrl}/client/${client.id}/config?qrcode=true`"/>
                                    </v-list-item-avatar>
                                </v-list-item>

                                <v-card-text class="text--primary">
                                    <v-chip
                                            v-for="(ip, i) in client.address"
                                            :key="i"
                                            color="indigo"
                                            text-color="white"
                                    >
                                        <v-icon left>mdi-ip-network</v-icon>
                                        {{ ip }}
                                    </v-chip>
                                </v-card-text>
                                <v-card-actions>
                                    <v-btn
                                            text
                                            :href="`${apiBaseUrl}/client/${client.id}/config?qrcode=false`"
                                    >
                                        <span class="d-none d-lg-flex">Download</span>
                                        <v-icon right dark>mdi-cloud-download-outline</v-icon>
                                    </v-btn>
                                    <v-btn
                                            text
                                            @click.stop="startUpdateClient(client)"
                                    >
                                        <span class="d-none d-lg-flex">Edit</span>
                                        <v-icon right dark>mdi-square-edit-outline</v-icon>
                                    </v-btn>
                                      <v-btn
                                              text
                                              @click="deleteClient(client)"
                                      >
                                          <span class="d-none d-lg-flex">Delete</span>
                                          <v-icon right dark>mdi-trash-can-outline</v-icon>
                                      </v-btn>
                                    <v-btn
                                            text
                                            @click="sendEmailClient(client.id)"
                                    >
                                        <span class="d-none d-lg-flex">Send Email</span>
                                        <v-icon right dark>mdi-email-send-outline</v-icon>
                                    </v-btn>
                                    <v-spacer/>
                                    <v-tooltip right>
                                        <template v-slot:activator="{ on }">
                                            <v-switch
                                                    dark
                                                    v-on="on"
                                                    color="success"
                                                    v-model="client.enable"
                                                    v-on:change="updateClient(client)"
                                            />
                                        </template>
                                        <span> {{client.enable ? 'Disable' : 'Enable'}} this client</span>
                                    </v-tooltip>

                                </v-card-actions>
                            </v-card>
                        </v-col>
                    </v-row>
                </v-card>
            </v-col>
        </v-row>
        <v-dialog
                v-if="client"
                v-model="dialogAddClient"
                max-width="550"
        >
            <v-card>
                <v-card-title class="headline">Add new client</v-card-title>
                <v-card-text>
                    <v-row>
                        <v-col
                                cols="12"
                        >
                            <v-form
                                    ref="form"
                                    v-model="valid"
                            >
                                <v-text-field
                                        v-model="client.name"
                                        label="Client friendly name"
                                        :rules="[ v => !!v || 'Client name is required', ]"
                                        required
                                />
                                <v-text-field
                                        v-model="client.email"
                                        label="Client email"
                                        :rules="[ v => (/.+@.+\..+/.test(v) || v === '') || 'E-mail must be valid',]"
                                />
                                <v-select
                                        v-model="client.address"
                                        :items="server.address"
                                        label="Client IP will be chosen from these networks"
                                        :rules="[ v => !!v || 'Network is required', ]"
                                        multiple
                                        chips
                                        persistent-hint
                                        required
                                />
                                <v-combobox
                                        v-model="client.allowedIPs"
                                        chips
                                        hint="Write IPv4 or IPv6 CIDR and hit enter"
                                        label="Allowed IPs"
                                        multiple
                                        dark
                                >
                                    <template v-slot:selection="{ attrs, item, select, selected }">
                                        <v-chip
                                                v-bind="attrs"
                                                :input-value="selected"
                                                close
                                                @click="select"
                                                @click:close="client.allowedIPs.splice(client.allowedIPs.indexOf(item), 1)"
                                        >
                                            <strong>{{ item }}</strong>&nbsp;
                                        </v-chip>
                                    </template>
                                </v-combobox>

                                <v-switch
                                        v-model="client.enable"
                                        color="red"
                                        inset
                                        :label="client.enable ? 'Enable client after creation': 'Disable client after creation'"
                                />
                                <v-switch
                                        v-model="client.ignorePersistentKeepalive"
                                        color="red"
                                        inset
                                        :label="'Ignore global persistent keepalive: ' + (client.ignorePersistentKeepalive ? 'Yes': 'NO')"
                                />
                            </v-form>
                        </v-col>
                    </v-row>
                </v-card-text>
                <v-card-actions>
                    <v-spacer/>
                    <v-btn
                            :disabled="!valid"
                            color="success"
                            @click="addClient(client)"
                    >
                        Submit
                        <v-icon right dark>mdi-check-outline</v-icon>
                    </v-btn>
                    <v-btn
                            color="primary"
                            @click="dialogAddClient = false"
                    >
                        Cancel
                        <v-icon right dark>mdi-close-circle-outline</v-icon>
                    </v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>
        <v-dialog
                v-if="client"
                v-model="dialogEditClient"
                max-width="550"
        >
            <v-card>
                <v-card-title class="headline">Edit client</v-card-title>
                <v-card-text>
                    <v-row>
                        <v-col
                                cols="12"
                        >
                            <v-form
                                    ref="form"
                                    v-model="valid"
                            >
                                <v-text-field
                                        v-model="client.name"
                                        label="Friendly name"
                                        :rules="[ v => !!v || 'Client name is required',]"
                                        required
                                />
                                <v-text-field
                                        v-model="client.email"
                                        label="Email"
                                        :rules="[ v => (/.+@.+\..+/.test(v) || v === '') || 'E-mail must be valid',]"
                                        required
                                />
                                <v-combobox
                                        v-model="client.address"
                                        chips
                                        hint="Write IPv4 or IPv6 CIDR and hit enter"
                                        label="Addresses"
                                        multiple
                                        dark
                                >
                                    <template v-slot:selection="{ attrs, item, select, selected }">
                                        <v-chip
                                                v-bind="attrs"
                                                :input-value="selected"
                                                close
                                                @click="select"
                                                @click:close="client.address.splice(client.address.indexOf(item), 1)"
                                        >
                                            <strong>{{ item }}</strong>&nbsp;
                                        </v-chip>
                                    </template>
                                </v-combobox>
                                <v-combobox
                                        v-model="client.allowedIPs"
                                        chips
                                        hint="Write IPv4 or IPv6 CIDR and hit enter"
                                        label="Allowed IPs"
                                        multiple
                                        dark
                                >
                                    <template v-slot:selection="{ attrs, item, select, selected }">
                                        <v-chip
                                                v-bind="attrs"
                                                :input-value="selected"
                                                close
                                                @click="select"
                                                @click:close="client.allowedIPs.splice(client.allowedIPs.indexOf(item), 1)"
                                        >
                                            <strong>{{ item }}</strong>&nbsp;
                                        </v-chip>
                                    </template>
                                </v-combobox>
                                <v-switch
                                        v-model="client.ignorePersistentKeepalive"
                                        color="red"
                                        inset
                                        :label="'Ignore global persistent keepalive: ' + (client.ignorePersistentKeepalive ? 'Yes': 'NO')"
                                />
                            </v-form>
                        </v-col>
                    </v-row>
                </v-card-text>
                <v-card-actions>
                    <v-spacer/>
                    <v-btn
                            :disabled="!valid"
                            color="success"
                            @click="updateClient(client)"
                    >
                        Submit
                        <v-icon right dark>mdi-check-outline</v-icon>
                    </v-btn>
                    <v-btn
                            color="primary"
                            @click="dialogEditClient = false"
                    >
                        Cancel
                        <v-icon right dark>mdi-close-circle-outline</v-icon>
                    </v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>
        <Notification v-bind:notification="notification"/>
    </v-container>
</template>
<script>
  import {ApiService, API_BASE_URL} from '../services/ApiService'
  import Notification from '../components/Notification'

  export default {
    name: 'Clients',

    components: {
      Notification
    },

    data: () => ({
      api: null,
      apiBaseUrl: API_BASE_URL,
      clients: [],
      notification: {
        show: false,
        color: '',
        text: '',
      },
      dialogAddClient: false,
      dialogEditClient: false,
      client: null,
      server: null,
      valid: false,
    }),

    mounted () {
      this.api = new ApiService();
      this.getClients();
      this.getServer()
    },

    methods: {
      getClients() {
        this.api.get('/client').then((res) => {
          this.clients = res
        }).catch((e) => {
          this.notify('error', e.response.status + ' ' + e.response.statusText);
        });
      },

      getServer() {
        this.api.get('/server').then((res) => {
          this.server = res;
        }).catch((e) => {
          this.notify('error', e.response.status + ' ' + e.response.statusText);
        });
      },

      startAddClient() {
        this.dialogAddClient = true;
        this.client = {
          name: "",
          email: "",
          enable: true,
          allowedIPs: ["0.0.0.0/0", "::/0"],
          address: this.server.address,
        }
      },
      addClient(client) {
        if (client.allowedIPs.length < 1) {
          this.notify('error', 'Please provide at least one valid CIDR address for client allowed IPs');
          return;
        }
        for (let i = 0; i < client.allowedIPs.length; i++){
          if (this.$isCidr(client.allowedIPs[i]) === 0) {
            this.notify('error', 'Invalid CIDR detected, please correct before submitting');
            return
          }
        }
        this.dialogAddClient = false;

        this.api.post('/client', client).then((res) => {
          this.notify('success', `Client ${res.name} successfully added`);
          this.getClients()
        }).catch((e) => {
          this.notify('error', e.response.status + ' ' + e.response.statusText);
        });
      },

      deleteClient(client) {
        if(confirm(`Do you really want to delete ${client.name} ?`)){
          this.api.delete(`/client/${client.id}`).then((res) => {
            this.notify('success', "Client successfully deleted");
            this.getClients()
          }).catch((e) => {
            this.notify('error', e.response.status + ' ' + e.response.statusText);
          });
        }
      },

      sendEmailClient(id) {
        this.api.get(`/client/${id}/email`).then((res) => {
          this.notify('success', "Email successfully sent");
          this.getClients()
        }).catch((e) => {
          this.notify('error', e.response.status + ' ' + e.response.statusText);
        });
      },

      startUpdateClient(client) {
        this.client = client;
        this.dialogEditClient = true;
      },
      updateClient(client) {
        // check allowed IPs
        if (client.allowedIPs.length < 1) {
          this.notify('error', 'Please provide at least one valid CIDR address for client allowed IPs');
          return;
        }
        for (let i = 0; i < client.allowedIPs.length; i++){
          if (this.$isCidr(client.allowedIPs[i]) === 0) {
            this.notify('error', 'Invalid CIDR detected, please correct before submitting');
            return
          }
        }
        // check address
        if (client.address.length < 1) {
          this.notify('error', 'Please provide at least one valid CIDR address for client');
          return;
        }
        for (let i = 0; i < client.address.length; i++){
          if (this.$isCidr(client.address[i]) === 0) {
            this.notify('error', 'Invalid CIDR detected, please correct before submitting');
            return
          }
        }
        // all good, submit
        this.dialogEditClient = false;

        this.api.patch(`/client/${client.id}`, client).then((res) => {
          this.notify('success', `Client ${res.name} successfully updated`);
          this.getClients()
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
