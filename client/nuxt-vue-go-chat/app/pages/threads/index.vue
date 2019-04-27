<template>
  <div style="width: 100%">
    <div>
      <ThreadsInput />
    </div>
    <div class="fixed">
      <v-btn color="success" @click="removeButton()">Create</v-btn>
    </div>

    <div class="list">
      <v-data-table
        :headers="headers"
        :items="showThreads"
        :no-data-text="nodata"
        :pagination.sync="pagination"
        :loading="loading"
      >
        <template slot="headerCell" slot-scope="props">
          <v-tooltip bottom>
            <template v-slot:activator="{ on }">
              <span style="color: #64B5F6" v-on="on">
                {{ props.header.text }}
              </span>
            </template>
            <span style="color: #64B5F6">
              {{ props.header.text }}
            </span>
          </v-tooltip>
        </template>
        <template v-slot:items="props">
          <tr @click="redirectToComment(props.item)">
            <td>{{ props.item.title }}</td>
            <td>{{ props.item.user.name }}</td>
            <td>{{ props.item.createdAt }}</td>
          </tr>
        </template>
      </v-data-table>
    </div>
  </div>
</template>

<script>
import moment from '~/plugins/moment'
import ThreadsInput from '~/components/ThreadsInput'
import { mapGetters, mapActions } from 'vuex'
import {
  LIST_THREADS,
  LIST_THREADS_MORE,
  CHANGE_IS_DIALOG_VISIBLE
} from '../../store/action-types'
export default {
  components: {
    ThreadsInput
  },

  data() {
    return {
      nodata: 'There is no data',
      loading: false,
      pagination: {},
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
  watch: {
    pagination: {
      handler() {
        this.listMore()
      },
      deep: true
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
        return thread
      })
    },
    dialogVisible() {
      return this.isDialogVisible
    },
    ...mapGetters('threads', ['threads', 'threadList', 'isDialogVisible'])
  },
  asyncData({ store }) {
    try {
      store.dispatch(`threads/${LIST_THREADS}`)
    } catch (error) {
      console.error(`failed to list threads: ${JSON.stringify(error)}`)
    }
  },
  methods: {
    redirectToComment(thread) {
      this.$router.push(`/threads/${thread.id}/comments`)
    },
    async listMore() {
      if (!this.threadList.hasNext) {
        console.log(
          `this.threads.hasNext ========>>${JSON.stringify(
            this.threadList.hasNext
          )}`
        )
        return
      }

      this.loading = true
      console.log('====X====')
      const lastId = this.threads[this.threads.length - 1].id
      console.log(`this.threads ========>>${JSON.stringify(this.threads)}`)

      try {
        await this.LIST_THREADS_MORE({ limit: 20, cursor: lastId })
        console.log(`this.threads = ${JSON.stringify(this.threads)}`)
      } catch (error) {
        console.error(`failed to list threads more: ${JSON.stringify(error)}`)
      }
      this.loading = false
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
  top: 90%;
  right: 10%;
  z-index: 100;
}
.list {
  z-index: 50;
}
</style>
