<template>
  <div>
    <form>
      <v-textarea v-model="content" box label="Comment" auto-grow></v-textarea>
      <v-btn @click="submit">Submit</v-btn>
      <v-btn @click="clear">Clear</v-btn>
    </form>
  </div>
</template>

<script>
import { validationMixin } from 'vuelidate'
import { required, minLength, maxLength } from 'vuelidate/lib/validators'
import { mapGetters, mapActions } from 'vuex'
import { SAVE_COMMENT } from '../store/action-types'
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
    ...mapGetters(['user'])
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
        this.snackbar.color = 'error'
        this.snackbar.text = 'fail sign up\nsystem error occur'
        this.snackbar.isOpen = true
      }

      this.content = ''
    },
    clear() {
      this.$v.$reset()
      this.content = ''
    },
    ...mapActions('comments', [SAVE_COMMENT])
  }
}
</script>
