<template>
  <v-app>
     <v-content>
      <FragClient/>
    </v-content>
  </v-app>
</template>

<script>
import FragClient from './components/FragClient'

export default {
  name: 'App',
  components: {
    FragClient
  },
  data () {
    return {
      //
    }
  }
}
</script>
