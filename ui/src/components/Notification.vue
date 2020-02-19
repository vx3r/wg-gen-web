<template>
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
</template>
<script>
  export default {
    name: 'Notification',
    props: ['notification'],
    data: () => ({
    }),
  };
</script>
