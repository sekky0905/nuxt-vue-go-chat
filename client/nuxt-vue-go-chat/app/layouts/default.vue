<template>
  <div id="app">
    <v-app id="inspire" dark>
      <Navigation />
      <Toolbar />
      <Container />
      <Footer />
    </v-app>
  </div>
</template>
<script>
import Navigation from '../components/Navigation'
import Toolbar from '../components/Toolbar'
import Container from '../components/Container'
import Footer from '../components/Footer'
export default {
  components: {
    Navigation,
    Toolbar,
    Container,
    Footer
  }
}
</script>
