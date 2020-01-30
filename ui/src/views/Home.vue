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
                          v => !!v || 'Name is required',
                        ]"
                      label="Friendly server name"
                      required
              />
              <v-text-field
                      type="number"
                      v-model="server.persistentKeepalive"
                      label="Persistent keepalive for clients"
                      :rules="[
                          v => !!v || 'Persistent keepalive is required',
                        ]"
                      required
              />
              <v-text-field
                      v-model="server.endpoint"
                      label="Endpoint for clients to connect to"
                      :rules="[
                          v => !!v || 'Endpoint is required',
                        ]"
                      required
              />
              <v-text-field
                      v-model="server.address"
                      label="Server interface addresses"
                      :rules="[
                          v => !!v || 'Server interface address is required',
                        ]"
                      required
              />
            </v-col>
            <v-col cols="6">
              <v-text-field
                      v-model="server.publicKey"
                      label="Server public key"
                      disabled
              />
              <v-text-field
                      v-model="server.presharedKey"
                      label="Preshared Key key"
                      disabled
              />
              <v-text-field
                      v-model="server.dns"
                      label="DNS servers for clients"
                      :rules="[
                          v => !!v || 'DNS server is required',
                        ]"
                      required
              />
              <v-text-field
                      v-model="server.listenPort"
                      type="number"
                      :rules="[
                          v => !!v || 'Listen port is required',
                        ]"
                      label="Server listen port"
                      required
              />
            </v-col>
          </div>

          <v-card-actions>
            <v-spacer/>
            <v-btn
                    color="warning"
                    @click="updateServer"
            >
              Update server configuration
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
                    @click.stop="dialogAddClient = true"
            >
              Add new client
            </v-btn>
          </v-list-item>
          <v-row>
            <v-col
                    v-for="(client, i) in clients"
                    :key="i"
                    cols="6"
            >
              <v-card
                      color="#1F7087"
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
                    Download configuration
                    <v-icon right dark>mdi-cloud-download</v-icon>
                  </v-btn>
                  <v-btn
                          text
                          @click.stop="dialogEditClient = true; clientToEdit = client"
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
                  <v-spacer/>
                  <v-tooltip right>
                    <template v-slot:activator="{ on }">
                      <v-switch
                              dark
                              v-on="on"
                              color="success"
                              v-model="client.enable"
                      />
                    </template>
                    <span>Enable or disable this client</span>
                  </v-tooltip>

                </v-card-actions>
              </v-card>
            </v-col>
          </v-row>
        </v-card>
      </v-col>
    </v-row>
    <v-dialog
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
                  @click="addClient()"
          >
            Submit
          </v-btn>
          <v-btn
                  color="primary"
                  @click="dialogAddClient = false"
          >
            Cancel
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog
            v-if="clientToEdit"
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
                        v-model="clientToEdit.name"
                        label="Client friendly name"
                        :rules="[
                          v => !!v || 'Client name is required',
                        ]"
                        required
                />
                <v-text-field
                        v-model="clientToEdit.email"
                        label="Client email"
                        :rules="[
                        v => !!v || 'E-mail is required',
                        v => /.+@.+\..+/.test(v) || 'E-mail must be valid',
                      ]"
                        required
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
                  @click="updateClient()"
          >
            Submit
          </v-btn>
          <v-btn
                  color="primary"
                  @click="dialogEditClient = false"
          >
            Cancel
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
      clientToEdit: null,
      client: {
        name: "",
        email: "",
        enable: true,
        allowedIPs: "0.0.0.0/0,::/0",
        address: "",
      }
    }),

    methods: {
      getData() {
        this.$get('/server').then((res) => {
          this.server = res;
          this.clientAddress = this.serverAddress = this.server.address.split(',')
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
        this.$patch('/server', this.server).then((res) => {
          this.notify('success', "Server successfully updated");
          this.getData()
        }).catch((e) => {
          this.notify('error', e.response.status + ' ' + e.response.statusText);
        });
      },
      addClient () {
        this.dialogAddClient = false;
        this.client.address = this.clientAddress.join(',');
        this.$post('/client', this.client).then((res) => {
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
      updateClient() {
        this.dialogEditClient = false;
        this.$patch(`/client/${this.clientToEdit.id}`, this.clientToEdit).then((res) => {
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
