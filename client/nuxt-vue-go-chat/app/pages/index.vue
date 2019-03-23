<template>
  <div>
    <div v-if="radios === 'login'">
      TODO LOGIN Component
    </div>
    <div v-else>
      <SignUp />
    </div>
    <div>
      <v-radio-group v-model="radios" :mandatory="false">
        <v-radio label="login" value="login">Login</v-radio>
        <v-radio label="signUp" value="signUp">Sign Up</v-radio>
      </v-radio-group>
    </div>
  </div>
</template>

<script>
import SignUp from '~/components/SignUp'
export default {
  components: {
    SignUp
  },
  data() {
    return {
      radios: 'login'
    }
  }
}
</script>
