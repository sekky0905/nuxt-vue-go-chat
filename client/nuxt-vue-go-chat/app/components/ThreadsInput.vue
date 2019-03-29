<template>
  <div>
    <v-dialog v-model="dialogVisible" persistent max-width="600px">
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
          <v-btn color="info" @click="submitForm('formData')">Create</v-btn>
          <v-btn color="error" @click="closeDialogState()">Cancel</v-btn>
          <v-btn color="warning" @click="resetForm('formData')">Reset</v-btn>
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
    dialogVisible() {
      return this.isDialogVisible
    },
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
    ...mapGetters('threads', ['threads', 'isDialogVisible']),
    ...mapGetters(['user'])
  },
  methods: {
    async submitForm() {
      this.closeDialogState()
      const payload = {
        user: this.user,
        ...this.title
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
      this.CHANGE_IS_DIALOG_VISIBLE({ dialogState: this.dialogVisible })
    },
    ...mapActions('threads', [SAVE_THREAD, CHANGE_IS_DIALOG_VISIBLE])
  }
}
</script>
