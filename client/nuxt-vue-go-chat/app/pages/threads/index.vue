<template>
  <div>
    <div>
      <ThreadsInput />
    </div>
    <div class="fixed">
      <v-btn color="warning" @click="removeButton()">Create</v-btn>
    </div>

    <div class="list">
      <v-data-table :headers="headers" :items="showThreads" class="elevation-1">
        <template slot="headerCell" slot-scope="props">
          <v-tooltip bottom>
            <template v-slot:activator="{ on }">
              <span v-on="on">
                {{ props.header.text }}
              </span>
            </template>
            <span>
              {{ props.header.text }}
            </span>
          </v-tooltip>
        </template>
        <template v-slot:items="props">
          <td class="text-xs-right">{{ props.item.title }}</td>
          <td class="text-xs-right">{{ props.item.user.name }}</td>
          <td class="text-xs-right">{{ props.item.createdAt }}</td>
        </template>
      </v-data-table>
    </div>
    <div v-if="existsMore">
      <infinite-loading
        @infinite="listMore(showThreads[showThreads.length - 1].id)"
      />
    </div>
  </div>
</template>

<script>
import moment from '~/plugins/moment'
import ThreadsInput from '~/components/ThreadsInput'
import { mapGetters, mapActions } from 'vuex'
import InfiniteLoading from 'vue-infinite-loading'
import {
  LIST_THREADS,
  LIST_THREADS_MORE,
  CHANGE_IS_DIALOG_VISIBLE
} from '../../store/action-types'
export default {
  components: {
    ThreadsInput,
    InfiniteLoading
  },

  data() {
    return {
      headers: [
        {
          text: 'Title',
          align: 'left',
          sortable: false,
          value: 'title'
        },
        {
          text: 'createdBy',
          align: 'left',
          sortable: false,
          value: 'user.name'
        },
        {
          text: 'createdAt',
          align: 'left',
          sortable: false,
          value: 'createdAt'
        }
      ]
    }
  },
  computed: {
    existsData() {
      return this.threadList && this.threadList.length !== 0
    },
    existsMore() {
      return this.threadList.hasNext
    },
    showThreads() {
      if (!this.existsData) {
        return
      }
      return this.threads.map(thread => {
        thread.createdAt = moment(thread.createdAt).format(
          'YYYY/MM/DD HH:mm:ss'
        )
        console.log(`thread => ${JSON.stringify(thread)}`)
        return thread
      })
    },
    dialogVisible() {
      return this.isDialogVisible
    },
    ...mapGetters('threads', ['threads', 'threadList', 'isDialogVisible'])
  },
  async asyncData({ store }) {
    try {
      await store.dispatch(`threads/${LIST_THREADS}`)
    } catch (e) {
      console.log(`threads のe==> ${JSON.stringify(e)}`)
    }
  },
  methods: {
    handleClick(thread) {
      this.$router.push(`/threads/${thread.id}`)
    },
    async listMore(id) {
      try {
        await this.LIST_THREADS_MORE({ limit: 20, cursor: id })
      } catch (error) {
        console.error(`error has occurred => ${JSON.stringify(error)}`)
      }
    },
    removeButton() {
      this.CHANGE_IS_DIALOG_VISIBLE({ dialogState: this.dialogVisible })
    },
    ...mapActions('threads', [
      LIST_THREADS,
      LIST_THREADS_MORE,
      CHANGE_IS_DIALOG_VISIBLE
    ])
  }
}
</script>

<style scoped>
/* Fix button position */
.fixed {
  position: fixed;
  top: 40%;
  right: 30px;
  z-index: 100;
}
.list {
  z-index: 50;
}
</style>
