<template>
  <div>
    <div>
      <form>
        <v-text-field
          v-model="name"
          type="text"
          :error-messages="nameErrors"
          :counter="10"
          label="Name"
          required
          @input="$v.name.$touch()"
          @blur="$v.name.$touch()"
        ></v-text-field>

        <v-text-field
          v-model="pw"
          append-icon="visibility_off"
          type="password"
          :error-messages="passwordErrors"
          label="Password"
          required
          @input="$v.pw.$touch()"
          @blur="$v.pw.$touch()"
        ></v-text-field>

        <v-btn @click="submit">submit</v-btn>
        <v-btn @click="clear">clear</v-btn>
      </form>
    </div>
    <div>
      <v-snackbar
        v-model="snackbar.isOpen"
        :color="snackbar.color"
        :multi-line="true"
        :timeout="500"
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
import { validationMixin } from 'vuelidate'
import { required, minLength, maxLength } from 'vuelidate/lib/validators'
import { mapActions, mapGetters } from 'vuex'
import { LOGIN } from '../store/action-types'

export default {
  name: 'SignUp',
  mixins: [validationMixin],

  validations: {
    name: { required, minLength: minLength(3), maxLength: maxLength(10) },
    pw: { required, minLength: minLength(10), maxLength: maxLength(32) }
  },

  data: () => ({
    name: '',
    pw: '',
    snackbar: {
      isOpen: false,
      color: '',
      text: ''
    }
  }),

  computed: {
    nameErrors() {
      const errors = []
      if (!this.$v.name.$dirty) return errors
      !this.$v.name.minLength &&
        errors.push('Name must be at least 3 characters long')
      !this.$v.name.maxLength &&
        errors.push('Name must be at most 10 characters long')
      !this.$v.name.required && errors.push('Name is required.')
      return errors
    },
    passwordErrors() {
      const errors = []
      if (!this.$v.pw.$dirty) return errors
      !this.$v.pw.minLength &&
        errors.push('Password must be at least 10 characters long')
      !this.$v.pw.maxLength &&
        errors.push('Password must be at most 32 characters long')
      !this.$v.pw.required && errors.push('Password is required.')
      return errors
    },
    ...mapGetters(['user'])
  },

  methods: {
    async submit() {
      this.$v.$touch()
      try {
        await this.LOGIN({ name: this.name, password: this.pw })
        this.snackbar.color = 'success'
        this.snackbar.text = 'success sign up'
        this.snackbar.isOpen = true
        // for snackbar
        const redirect = () => {
          this.$router.push('/threads/')
        }
        setTimeout(redirect, 500)
      } catch (error) {
        console.error(`failed to login: ${JSON.stringify(error)}`)
        if (error.response.data.code === 'AuthenticationFailure') {
          this.snackbar.color = 'error'
          this.snackbar.text = 'fail to sign up\nName or Password is invalid'
          this.snackbar.isOpen = true
        } else {
          this.snackbar.color = 'error'
          this.snackbar.text = 'fail to sign up\nsystem error occur'
          this.snackbar.isOpen = true
        }
      }
    },
    clear() {
      this.$v.$reset()
      this.name = ''
      this.pw = ''
    },
    ...mapActions([LOGIN])
  }
}
</script>

<style scoped></style>
