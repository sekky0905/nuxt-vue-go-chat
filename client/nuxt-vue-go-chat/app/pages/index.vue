<template>
  <div>
    <div v-if="radios === 'login'">
      <Login />
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
import Login from '~/components/Login'
import SignUp from '~/components/SignUp'
export default {
  components: {
    Login,
    SignUp
  },
  data() {
    return {
      radios: 'login'
    }
  }
}
</script>
