<template>
  <v-content>
    <v-row v-if="server">
      <v-col cols="12">
        <v-card dark>
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title class="headline">Server configurations</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <div class="d-flex flex-no-wrap justify-space-between">
            <v-col cols="6">
              <v-text-field
                      v-model="server.name"
                      :rules="[
                          v => !!v || 'Friendly name is required',
                        ]"
                      label="Friendly name"
                      required
              />
              <v-text-field
                      type="number"
                      v-model="server.persistentKeepalive"
                      label="Persistent keepalive"
                      :rules="[
                          v => !!v || 'Persistent keepalive is required',
                        ]"
                      required
              />
              <v-text-field
                      v-model="server.endpoint"
                      label="Public endpoint for clients to connect to"
                      :rules="[
                          v => !!v || 'Public endpoint for clients to connect to is required',
                        ]"
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
              <v-text-field
                      v-model="server.preUp"
                      label="PreUp: script snippets which will be executed by bash before setting up the interface"
              />
              <v-text-field
                      v-model="server.postUp"
                      label="PostUp: script snippets which will be executed by bash after setting up the interface"
              />
            </v-col>
            <v-col cols="6">
              <v-text-field
                      v-model="server.publicKey"
                      label="Public key"
                      disabled
              />
              <v-text-field
                      v-model="server.presharedKey"
                      label="Preshared key"
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
                      v-model="server.preDown"
                      label="PreDown: script snippets which will be executed by bash before setting down the interface"
              />
              <v-text-field
                      v-model="server.postDown "
                      label="PostDown : script snippets which will be executed by bash after setting down the interface"
              />
            </v-col>
          </div>

          <v-card-actions>
            <v-spacer/>
            <v-btn
                    class="ma-2"
                    color="warning"
                    @click="updateServer"
            >
              Update server configuration
              <v-icon right dark>mdi-update</v-icon>
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
    <v-divider dark/>
    <v-row>
      <v-col cols="12">
        <v-card dark>
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title class="headline">Clients</v-list-item-title>
            </v-list-item-content>
            <v-btn
                    color="success"
                    @click.stop="startAddClient"
            >
              Add new client
              <v-icon right dark>mdi-account-multiple-plus-outline</v-icon>
            </v-btn>
          </v-list-item>
          <v-row>
            <v-col
                    v-for="(client, i) in clients"
                    :key="i"
                    cols="6"
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
                    <v-img :src="getUrlToConfig(client.id, true)"/>
                  </v-list-item-avatar>
                </v-list-item>

                <v-card-text class="text--primary">
                  <v-chip
                          v-for="(ip, i) in client.address.split(',')"
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
                          :href="getUrlToConfig(client.id, false)"
                  >
                    Download
                    <v-icon right dark>mdi-cloud-download-outline</v-icon>
                  </v-btn>
                  <v-btn
                          text
                          @click.stop="editClient(client.id)"
                  >
                    Edit
                    <v-icon right dark>mdi-square-edit-outline</v-icon>
                  </v-btn>
                  <v-btn
                          text
                          @click="deleteClient(client.id)"
                  >
                    Delete
                    <v-icon right dark>mdi-trash-can-outline</v-icon>
                  </v-btn>
                  <v-btn
                          text
                          @click="sendEmailClient(client.id)"
                  >
                    Send email
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
                              v-on:change="disableClient(client)"
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
                        :rules="[
                          v => !!v || 'Client name is required',
                        ]"
                        required
                />
                <v-text-field
                        v-model="client.email"
                        label="Client email"
                        :rules="[
                        v => !!v || 'E-mail is required',
                        v => /.+@.+\..+/.test(v) || 'E-mail must be valid',
                      ]"
                        required
                />
                <v-select
                        v-model="clientAddress"
                        :items="serverAddress"
                        label="Client IP will be chosen from these networks"
                        :rules="[
                                v => !!v || 'Network is required',
                        ]"
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
                        :rules="[
                          v => !!v || 'Client name is required',
                        ]"
                        required
                />
                <v-text-field
                        v-model="client.email"
                        label="Email"
                        :rules="[
                        v => !!v || 'Email is required',
                        v => /.+@.+\..+/.test(v) || 'Email must be valid',
                        ]"
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
    <v-snackbar
            v-model="notification.show"
            :right="true"
            :top="true"
            :color="notification.color"
    >
      {{ notification.text }}
      <v-btn
              dark
              text
              @click="notification.show = false"
      >
        Close
      </v-btn>
    </v-snackbar>
  </v-content>
</template>

<script>
  export default {
    name: 'home',
    mounted () {
      this.getData()
    },
    data: () => ({
      notification: {
        show: false,
        color: '',
        text: '',
      },
      valid: true,
      checkbox: false,
      server: null,
      clients: [],
      ipDns: "",
      ipAddress: "",
      clientAddress: [],
      serverAddress: [],
      dialogAddClient: false,
      dialogEditClient: false,
      client: null,
    }),
    methods: {
      startAddClient() {
        this.dialogAddClient = true;
        this.client = {
          name: "",
          email: "",
          enable: true,
          allowedIPs: ["0.0.0.0/0", "::/0"],
          address: "",
        }
      },
      editClient(id) {
        this.$get(`/client/${id}`).then((res) => {
          this.dialogEditClient = true;
          res.allowedIPs = res.allowedIPs.split(',');
          res.address = res.address.split(',');
          this.client = res
        }).catch((e) => {
          this.notify('error', e.response.status + ' ' + e.response.statusText);
        });
      },
      disableClient(client) {
        if(!Array.isArray(client.allowedIPs)){
          client.allowedIPs = client.allowedIPs.split(',');
        }
        if(!Array.isArray(client.address)){
          client.address = client.address.split(',');
        }
        this.updateClient(client)
      },
      getData() {
        this.$get('/server').then((res) => {
          res.address = res.address.split(',');
          res.dns = res.dns.split(',');
          this.server = res;
          this.clientAddress = this.serverAddress = this.server.address
        }).catch((e) => {
          this.notify('error', e.response.status + ' ' + e.response.statusText);
        });

        this.$get('/client').then((res) => {
          this.clients = res
        }).catch((e) => {
          this.notify('error', e.response.status + ' ' + e.response.statusText);
        });
      },
      updateServer () {
        // convert int values
        this.server.listenPort = parseInt(this.server.listenPort, 10);
        this.server.persistentKeepalive = parseInt(this.server.persistentKeepalive, 10);
        // check server addresses
        if (this.server.address.length < 1) {
          this.notify('error', 'Please provide at least one valid CIDR address for server interface');
          return;
        }
        for (let i = 0; i < this.server.address.length; i++){
          if (this.$isCidr(this.server.address[i]) === 0) {
            this.notify('error', 'Invalid CIDR detected, please correct before submitting');
            return
          }
        }
        this.server.address = this.server.address.join(',');
        this.server.dns = this.server.dns.join(',');

        this.$patch('/server', this.server).then((res) => {
          this.notify('success', "Server successfully updated");
          this.getData()
        }).catch((e) => {
          this.notify('error', e.response.status + ' ' + e.response.statusText);
        });
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
        client.address = this.clientAddress.join(',');
        client.allowedIPs = this.client.allowedIPs.join(',');

        this.$post('/client', client).then((res) => {
          this.notify('success', "Client successfully added");
          this.getData()
        }).catch((e) => {
          this.notify('error', e.response.status + ' ' + e.response.statusText);
        });
      },
      deleteClient(id) {
        if(confirm("Do you really want to delete?")){
          this.$delete(`/client/${id}`).then((res) => {
            this.notify('success', "Client successfully deleted");
            this.getData()
          }).catch((e) => {
            this.notify('error', e.response.status + ' ' + e.response.statusText);
          });
        }
      },
      sendEmailClient(id) {
        this.$get(`/client/${id}/email`).then((res) => {
          this.notify('success', "Email successfully sent");
          this.getData()
        }).catch((e) => {
          this.notify('error', e.response.status + ' ' + e.response.statusText);
        });
      },
      getUrlToConfig(id, qrcode){
        let base = "/api/v1.0";
        if (process.env.NODE_ENV === "development"){
          base = process.env.VUE_APP_API_BASE_URL
        }
        if (qrcode){
          return `${base}/client/${id}/config?qrcode=true`
        } else {
          return `${base}/client/${id}/config`
        }
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
        client.allowedIPs = client.allowedIPs.join(',');
        client.address = client.address.join(',');

        this.$patch(`/client/${client.id}`, client).then((res) => {
          this.notify('success', "Client successfully updated");
          this.getData()
        }).catch((e) => {
          this.notify('error', e.response.status + ' ' + e.response.statusText);
        });
      },
      notify(color, msg) {
        this.notification.show = true;
        this.notification.color = color;
        this.notification.text = msg;
      },

    },
  }
</script>
