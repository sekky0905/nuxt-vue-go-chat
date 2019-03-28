<template>
  <v-toolbar app fixed clipped-left>
    <v-toolbar-side-icon @click.stop="drawer = !drawer"></v-toolbar-side-icon>
    <v-toolbar-title>Application</v-toolbar-title>
  </v-toolbar>
</template>

<script>
export default {
  name: 'Toolbar',
  data: () => ({
    drawer: true
  })
}
</script>

<style scoped></style>
