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
                                @click="startCreate"
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
                                        <v-list-item-subtitle>Created: {{ client.created | formatDate }} by {{ client.createdBy }}</v-list-item-subtitle>
                                        <v-list-item-subtitle>Updated: {{ client.updated | formatDate }} by {{ client.updatedBy }}</v-list-item-subtitle>
                                    </v-list-item-content>

                                    <v-list-item-avatar
                                            tile
                                            size="150"
                                    >
                                        <v-img :src="'data:image/png;base64, ' + getClientQrcode(client.id)"/>
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
                                    <v-tooltip bottom>
                                        <template v-slot:activator="{ on }">
                                            <v-btn
                                                    text
                                                    v-on:click="forceFileDownload(client)"
                                                    v-on="on"
                                            >
                                                <span class="d-none d-lg-flex">Download</span>
                                                <v-icon right dark>mdi-cloud-download-outline</v-icon>
                                            </v-btn>
                                        </template>
                                        <span>Download</span>
                                    </v-tooltip>

                                    <v-tooltip bottom>
                                        <template v-slot:activator="{ on }">
                                            <v-btn
                                                    text
                                                    @click.stop="startUpdate(client)"
                                                    v-on="on"
                                            >
                                                <span class="d-none d-lg-flex">Edit</span>
                                                <v-icon right dark>mdi-square-edit-outline</v-icon>
                                            </v-btn>
                                        </template>
                                        <span>Edit</span>
                                    </v-tooltip>

                                    <v-tooltip bottom>
                                        <template v-slot:activator="{ on }">
                                            <v-btn
                                                    text
                                                    @click="remove(client)"
                                                    v-on="on"
                                            >
                                                <span class="d-none d-lg-flex">Delete</span>
                                                <v-icon right dark>mdi-trash-can-outline</v-icon>
                                            </v-btn>
                                        </template>
                                        <span>Delete</span>
                                    </v-tooltip>

                                    <v-tooltip bottom>
                                        <template v-slot:activator="{ on }">
                                            <v-btn
                                                    text
                                                    @click="email(client)"
                                                    v-on="on"
                                            >
                                                <span class="d-none d-lg-flex">Send Email</span>
                                                <v-icon right dark>mdi-email-send-outline</v-icon>
                                            </v-btn>
                                        </template>
                                        <span>Send Email</span>
                                    </v-tooltip>
                                    <v-spacer/>
                                    <v-tooltip right>
                                        <template v-slot:activator="{ on }">
                                            <v-switch
                                                    dark
                                                    v-on="on"
                                                    color="success"
                                                    v-model="client.enable"
                                                    v-on:change="update(client)"
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
                v-model="dialogCreate"
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
                            @click="create(client)"
                    >
                        Submit
                        <v-icon right dark>mdi-check-outline</v-icon>
                    </v-btn>
                    <v-btn
                            color="primary"
                            @click="dialogCreate = false"
                    >
                        Cancel
                        <v-icon right dark>mdi-close-circle-outline</v-icon>
                    </v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>
        <v-dialog
                v-if="client"
                v-model="dialogUpdate"
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
                            @click="update(client)"
                    >
                        Submit
                        <v-icon right dark>mdi-check-outline</v-icon>
                    </v-btn>
                    <v-btn
                            color="primary"
                            @click="dialogUpdate = false"
                    >
                        Cancel
                        <v-icon right dark>mdi-close-circle-outline</v-icon>
                    </v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>
    </v-container>
</template>
<script>
  import { mapActions, mapGetters } from 'vuex'

  export default {
    name: 'Clients',

    data: () => ({
      dialogCreate: false,
      dialogUpdate: false,
      client: null,
      valid: false,
    }),

    computed:{
      ...mapGetters({
        getClientQrcode: 'client/getClientQrcode',
        getClientConfig: 'client/getClientConfig',
        server: 'server/server',
        clients: 'client/clients',
        clientQrcodes: 'client/clientQrcodes',
      }),
    },

    mounted () {
      this.readAllClients()
      this.readServer()
    },

    methods: {
      ...mapActions('client', {
        errorClient: 'error',
        readAllClients: 'readAll',
        creatClient: 'create',
        updateClient: 'update',
        deleteClient: 'delete',
        emailClient: 'email',
      }),
      ...mapActions('server', {
        readServer: 'read',
      }),

      startCreate() {
        this.client = {
          name: "",
          email: "",
          enable: true,
          allowedIPs: this.server.allowedips,
          address: this.server.address,
        }
        this.dialogCreate = true;
      },

      create(client) {
        if (client.allowedIPs.length < 1) {
          this.errorClient('Please provide at least one valid CIDR address for client allowed IPs')
          return;
        }
        for (let i = 0; i < client.allowedIPs.length; i++){
          if (this.$isCidr(client.allowedIPs[i]) === 0) {
            this.errorClient('Invalid CIDR detected, please correct before submitting')
            return
          }
        }
        this.dialogCreate = false;
        this.creatClient(client)
      },

      remove(client) {
        if(confirm(`Do you really want to delete ${client.name} ?`)){
          this.deleteClient(client)
        }
      },

      email(client) {
        if (!client.email){
          this.errorClient('Client email is not defined')
          return
        }

        if(confirm(`Do you really want to send email to ${client.email} with all configurations ?`)){
          this.emailClient(client)
        }
      },

      startUpdate(client) {
        this.client = client;
        this.dialogUpdate = true;
      },

      update(client) {
        // check allowed IPs
        if (client.allowedIPs.length < 1) {
          this.errorClient('Please provide at least one valid CIDR address for client allowed IPs');
          return;
        }
        for (let i = 0; i < client.allowedIPs.length; i++){
          if (this.$isCidr(client.allowedIPs[i]) === 0) {
            this.errorClient('Invalid CIDR detected, please correct before submitting');
            return
          }
        }
        // check address
        if (client.address.length < 1) {
          this.errorClient('Please provide at least one valid CIDR address for client');
          return;
        }
        for (let i = 0; i < client.address.length; i++){
          if (this.$isCidr(client.address[i]) === 0) {
            this.errorClient('Invalid CIDR detected, please correct before submitting');
            return
          }
        }
        // all good, submit
        this.dialogUpdate = false;
        this.updateClient(client)
      },

      forceFileDownload(client){
        let config = this.getClientConfig(client.id)
        if (!config) {
          this.errorClient('Failed to download client config');
          return
        }
        const url = window.URL.createObjectURL(new Blob([config]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', 'wg0.conf') //or any other extension
        document.body.appendChild(link)
        link.click()
      },
    }
  };
</script>
