<template>
  <v-content>
    <Server/>
    <Clients/>
  </v-content>
</template>

<script>
  import Server from '../components/Server'
  import Clients from '../components/Clients'

  export default {
    name: 'home',
    components: {
      Server,
      Clients
    }
  }
</script>
