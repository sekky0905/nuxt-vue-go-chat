<template>
  <v-navigation-drawer v-model="drawer" clipped fixed>
    <v-list dense>
      <v-list-tile>
        <v-list-tile-action>
          <v-icon>home</v-icon>
        </v-list-tile-action>
        <v-list-tile-content>
          <v-list-tile-title>Home</v-list-tile-title>
        </v-list-tile-content>
      </v-list-tile>
    </v-list>
  </v-navigation-drawer>
</template>

<script>
export default {
  name: 'Navigation',
  data() {
    return {
      drawer: null
    }
  }
}
</script>

<style scoped></style>
