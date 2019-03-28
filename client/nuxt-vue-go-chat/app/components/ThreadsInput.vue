<template>
  <div>
    <v-dialog v-model="dialog" persistent max-width="600px">
      <v-card>
        <v-card-title>
          <span class="headline">Thread Title</span>
        </v-card-title>
        <v-card-text>
          <v-text-field
            v-model="title"
            type="text"
            :error-messages="titleErrors"
            :counter="10"
            label="Name"
            required
            @input="$v.title.$touch()"
            @blur="$v.title.$touch()"
          ></v-text-field>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <el-button type="primary" @click="submitForm('formData')"
            >Create</el-button
          >
          <el-button @click="closeDialogState()">Cancel</el-button>
          <el-button @click="resetForm('formData')">Reset</el-button>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import { validationMixin } from 'vuelidate'
import { required, minLength, maxLength } from 'vuelidate/lib/validators'
import { mapGetters, mapActions } from 'vuex'
import { SAVE_THREAD, CHANGE_IS_DIALOG_VISIBLE } from '../store/action-types'
export default {
  mixins: [validationMixin],
  data() {
    return {
      dialog: false,
      title: '',
      snackbar: {
        isOpen: false,
        color: '',
        text: ''
      }
    }
  },
  validations: {
    title: { required, minLength: minLength(1), maxLength: maxLength(20) }
  },
  computed: {
    titleErrors() {
      const errors = []
      if (!this.$v.title.$dirty) return errors
      !this.$v.title.minLength &&
        errors.push('Name must be at least 3 characters long')
      !this.$v.title.maxLength &&
        errors.push('Name must be at most 10 characters long')
      !this.$v.title.required && errors.push('Name is required.')
      return errors
    },
    ...mapGetters('threads', ['threads']),
    ...mapGetters(['user'])
  },
  methods: {
    async submitForm() {
      this.closeDialogState()
      const payload = {
        user: this.user,
        ...this.formData
      }
      try {
        await this.SAVE_THREAD({ payload: payload })
        this.snackbar.color = 'success'
        this.snackbar.text = `success create 【${this.title}】thread`
        this.snackbar.isOpen = true
      } catch (error) {
        this.snackbar.color = 'error'
        this.snackbar.text = 'fail sign up\nsystem error occur'
        this.snackbar.isOpen = true
      }
      this.title = ''
    },
    resetForm() {
      this.title = ''
    },
    closeDialogState() {
      this.dialog = false
    },
    ...mapActions('threads', [SAVE_THREAD, CHANGE_IS_DIALOG_VISIBLE])
  }
}
</script>
