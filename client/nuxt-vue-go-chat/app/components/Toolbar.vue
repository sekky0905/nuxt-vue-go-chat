<template>
  <div>
    <v-navigation-drawer v-model="drawer" clipped fixed app>
      <v-list dense>
        <v-list-tile @click="logout()">
          <v-list-tile-action>
            <v-icon>exit_to_app</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Dashboard</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
      </v-list>
    </v-navigation-drawer>
    <v-toolbar app fixed clipped-left>
      <v-toolbar-side-icon @click.stop="drawer = !drawer"></v-toolbar-side-icon>
      <v-toolbar-title style="color: #64B5F6">Nuxt Go Vue Chat</v-toolbar-title>
    </v-toolbar>
  </div>
</template>

<script>
export default {
  name: 'Toolbar',
  data: () => ({
    drawer: false
  }),
  methods: {
    logout() {
      // TODO implement logput
      console.log('logout')
    }
  }
}
</script>

<style scoped></style>
