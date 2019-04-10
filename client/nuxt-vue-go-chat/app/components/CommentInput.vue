<template>
  <div>
    <div>
      <v-dialog v-model="dialogVisible" persistent max-width="600px">
        <v-form>
          <v-card>
            <v-card-title>
              <span class="headline">Comment</span>
            </v-card-title>
            <v-card-text>
              <v-textarea
                v-model="content"
                box
                label="Comment"
                auto-grow
                :error-messages="contentError"
                @input="$v.content.$touch()"
                @blur="$v.content.$touch()"
              >
              </v-textarea>
            </v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="info" @click="submit">Create</v-btn>
              <v-btn color="error" @click="closeDialogState()">Cancel</v-btn>
              <v-btn color="warning" @click="clear()">Clear</v-btn>
            </v-card-actions>
          </v-card>
        </v-form>
      </v-dialog>
    </div>
    <div>
      <v-snackbar
        v-model="snackbar.isOpen"
        :color="snackbar.color"
        :multi-line="true"
        :timeout="500"
        :right="true"
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
import { mapGetters, mapActions } from 'vuex'
import { SAVE_COMMENT, CHANGE_IS_DIALOG_VISIBLE } from '../store/action-types'

export default {
  name: 'CommentInput',
  mixins: [validationMixin],
  validations: {
    content: { required, minLength: minLength(1), maxLength: maxLength(200) }
  },

  data: () => ({
    content: '',
    snackbar: {
      isOpen: false,
      color: '',
      text: ''
    }
  }),
  computed: {
    dialogVisible() {
      return this.isDialogVisible
    },
    contentError() {
      const errors = []
      if (!this.$v.content.$dirty) return errors
      !this.$v.content.minLength &&
        errors.push('Content must be at least 1 character long')
      !this.$v.content.maxLength &&
        errors.push('Content must be at most 200 characters long')
      !this.$v.content.required && errors.push('Content is required.')
      return errors
    },
    ...mapGetters(['user']),
    ...mapGetters('comments', ['isDialogVisible'])
  },
  methods: {
    async submit(f) {
      this.$v.$touch()

      const payload = {
        threadId: Number(this.$route.params.id),
        content: this.content,
        user: {
          id: this.user.id,
          name: this.user.name
        }
      }

      try {
        await this.SAVE_COMMENT({ payload: payload })
        this.snackbar.color = 'success'
        this.snackbar.text = 'success create comment'
        this.snackbar.isOpen = true
      } catch (error) {
        console.error(`failed to save comment: ${JSON.stringify(error)}`)
        this.snackbar.color = 'error'
        this.snackbar.text = 'fail to sign up\nsystem error occur'
        this.snackbar.isOpen = true
      }

      this.closeDialogState()
    },
    clear() {
      this.$v.$reset()
      this.content = ''
    },
    closeDialogState() {
      this.clear()
      this.CHANGE_IS_DIALOG_VISIBLE({ dialogState: this.dialogVisible })
    },
    ...mapActions('comments', [SAVE_COMMENT, CHANGE_IS_DIALOG_VISIBLE])
  }
}
</script>
