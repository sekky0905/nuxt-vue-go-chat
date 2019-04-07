<template>
  <div>
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
        <v-toolbar-side-icon
          @click.stop="drawer = !drawer"
        ></v-toolbar-side-icon>
        <v-toolbar-title style="color: #64B5F6"
          >Nuxt Go Vue Chat</v-toolbar-title
        >
      </v-toolbar>
    </div>
    <div>
      <v-snackbar
        v-model="snackbar.isOpen"
        :color="snackbar.color"
        :multi-line="true"
        :timeout="6000"
        :top="true"
      >
        {{ snackbar.text }}
        <v-btn color="pink" flat @click="snackbar.isOpen = false">
          Close
        </v-btn>
      </v-snackbar>
    </div>
  </div>
</template>

<script>
import { mapActions } from 'vuex'
import { LOGOUT } from '../store/action-types'
export default {
  name: 'Toolbar',
  data: () => ({
    drawer: false,
    snackbar: {
      isOpen: false,
      color: '',
      text: ''
    }
  }),
  methods: {
    async logout() {
      try {
        await this.LOGOUT()
        this.snackbar.color = 'success'
        this.snackbar.text = 'success sign up'
        this.snackbar.isOpen = true
        this.$router.push('/')
      } catch (error) {
        console.error(`failed to logout: ${error}`)
      }
    },
    ...mapActions([LOGOUT])
  }
}
</script>

<style scoped></style>
