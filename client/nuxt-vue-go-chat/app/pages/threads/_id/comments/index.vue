<template>
  <div class="comment-background">
    <div>
      <div class="fixed">
        <CommentInput />
      </div>
      <div class="fixed">
        <v-btn color="success" @click="showDialog()">Create</v-btn>
      </div>

      <div v-for="comment in showComments" :key="comment.id">
        <div v-if="isMyComment(comment.user.id)">
          <div class="my-comment">
            <div>
              <span class="comment-info">
                {{ comment.user.name }}: {{ comment.createdAt }}
              </span>
            </div>
            <span>
              <v-icon color="red darken-1" @click="removeComment(comment.id)"
                >delete_forever</v-icon
              >
            </span>
            <p>{{ comment.content }}</p>
          </div>
        </div>

        <div v-else>
          <div class="other-comment">
            <div class="chatting">
              <div>
                <span class="comment-info">
                  {{ comment.user.name }}: {{ comment.createdAt }}
                </span>
              </div>
              <div class="says">
                <p>{{ comment.content }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div v-if="existsMore">
        <infinite-loading @infinite="listMore" />
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import CommentInput from '~/components/CommentInput'
import moment from '~/plugins/moment'
import InfiniteLoading from 'vue-infinite-loading'
import {
  LIST_COMMENTS,
  LIST_COMMENTS_MORE,
  DELETE_COMMENT,
  CHANGE_IS_DIALOG_VISIBLE
} from '../../../../store/action-types'
export default {
  name: 'Index',
  components: {
    CommentInput,
    InfiniteLoading
  },
  computed: {
    existsData() {
      return this.commentList && this.commentList.length !== 0
    },
    existsMore() {
      return this.commentList.hasNext
    },
    showComments() {
      if (!this.existsData) {
        return
      }
      return this.comments.map(comment => {
        comment.createdAt = moment(comment.createdAt).format(
          'YYYY/MM/DD HH:mm:ss'
        )
        return comment
      })
    },
    dialogVisible() {
      return this.isDialogVisible
    },
    ...mapGetters(['user']),
    ...mapGetters('comments', ['comments', 'commentList', 'isDialogVisible'])
  },
  async asyncData({ store, params }) {
    try {
      await store.dispatch(`comments/${LIST_COMMENTS}`, { threadId: params.id })
    } catch (error) {
      console.error(`failed to list comments: ${JSON.stringify(error)}`)
    }
  },
  methods: {
    async removeComment(id) {
      try {
        await this.DELETE_COMMENT({ threadId: this.$route.params.id, id: id })
      } catch (error) {
        console.error(`failed to delete comment: ${JSON.stringify(error)}`)
      }
    },
    isMyComment(userId) {
      return this.user.id === userId
    },
    isOtherComment(userId) {
      return this.user.id !== userId
    },
    async listMore() {
      const comment = this.showComments[this.showComments.length - 1]
      try {
        await this.LIST_COMMENTS_MORE({
          threadId: comment.threadId,
          limit: 20,
          cursor: comment.id
        })
      } catch (error) {
        console.error(`failed to list comments more: ${JSON.stringify(error)}`)
      }
    },
    showDialog() {
      this.CHANGE_IS_DIALOG_VISIBLE({ dialogState: this.dialogVisible })
    },
    ...mapActions('comments', [
      LIST_COMMENTS_MORE,
      DELETE_COMMENT,
      CHANGE_IS_DIALOG_VISIBLE
    ])
  }
}
</script>

<style scoped>
/** comments background **/
.comment-background {
  padding: 20px 10px;
  margin: 20px auto;
  width: 95%;
  text-align: right;
  font-size: 15px;
  background: #424242;
  opacity: 0.7;
}

/** other's comments **/
.other-comment {
  width: 100%;
  margin: 10px 0;
  overflow: hidden;
}
.other-comment .chatting {
  width: 100%;
  text-align: left;
}
.says {
  display: inline-block;
  position: relative;
  margin: 0 0 0 10px;
  padding: 10px;
  max-width: 250px;
  border-radius: 12px;
  background: #bdbdbd;
}
.says:after {
  content: '';
  display: inline-block;
  position: absolute;
  top: 3px;
  left: -19px;
  border: 8px solid transparent;
  border-right: 18px solid #bdbdbd;
  -ms-transform: rotate(35deg);
  -webkit-transform: rotate(35deg);
  transform: rotate(35deg);
}
.says p {
  margin: 0;
  padding: 0;
}

/** my comments **/
.my-comment {
  margin: 10px 0;
}
.my-comment p {
  display: inline-block;
  position: relative;
  margin: 0 10px 0 0;
  padding: 8px;
  max-width: 250px;
  border-radius: 12px;
  background: #409eff;
  font-size: 15px;
  color: #edf1ee;
}
.my-comment p:after {
  content: '';
  position: absolute;
  top: 3px;
  right: -19px;
  border: 8px solid transparent;
  border-left: 18px solid #409eff;
  -ms-transform: rotate(-35deg);
  -webkit-transform: rotate(-35deg);
  transform: rotate(-35deg);
}
/* reference */
/* https://saruwakakun.com/html-css/reference/speech-bubble */
.user-name {
  color: #000;
}

/* Fix button position */
.fixed {
  position: fixed;
  top: 90%;
  right: 10%;
  z-index: 100;
}

.comment-info {
  color: #76ff03;
  border-bottom: solid 1px #76ff03;
}
</style>
