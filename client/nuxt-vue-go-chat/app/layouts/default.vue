<template>
  <div id="app">
    <v-app id="inspire" dark>
      <Navigation />
      <Toolbar />
      <v-content>
        <v-container fluid fill-height>
          <nuxt />
        </v-container>
      </v-content>
      <Footer />
    </v-app>
  </div>
</template>
<script>
import Navigation from '../components/Navigation'
import Toolbar from '../components/Toolbar'
import Footer from '../components/Footer'
export default {
  components: {
    Navigation,
    Toolbar,
    Footer
  }
}
</script>
