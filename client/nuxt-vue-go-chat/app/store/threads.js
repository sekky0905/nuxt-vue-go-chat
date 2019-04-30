import {
  LIST_THREADS,
  SAVE_THREAD,
  EDIT_THREAD,
  DELETE_THREAD,
  LIST_THREADS_MORE,
  CHANGE_IS_DIALOG_VISIBLE
} from './action-types'
import {
  SET_THREADS,
  SET_THREAD_LIST,
  ADD_THREAD_LIST,
  ADD_THREAD,
  UPDATE_THREAD,
  REMOVE_THREAD,
  CLEAR_THREADS,
  SET_IS_DIALOG_VISIBLE
} from './mutation-types'

export const state = () => ({
  threadList: {
    threads: [],
    hasNext: false,
    cursor: ''
  },
  isDialogVisible: false
})

export const getters = {
  threadList: state => state.threadList,
  threads: state => state.threadList.threads,
  isDialogVisible: state => state.isDialogVisible
}

export const mutations = {
  [SET_THREADS](state, { threads }) {
    state.threadList.threads = threads
  },
  [SET_THREAD_LIST](state, { threadList }) {
    state.threadList = threadList
  },
  [ADD_THREAD_LIST](state, { threadList }) {
    state.threadList.threads = state.threadList.threads.concat(
      threadList.threads
    )
    state.threadList.hasNext = threadList.hasNext
    state.threadList.cursor = threadList.cursor
  },
  [ADD_THREAD](state, { thread }) {
    state.threadList.threads.push(thread)
  },
  [EDIT_THREAD](state, { thread }) {
    state.threadList.threads = state.threadList.threads.map(t =>
      t.id === thread.id ? thread : t
    )
  },
  [REMOVE_THREAD](state, { id }) {
    state.threadList.threads = state.threadList.threads.filter(t => t.id !== id)
  },
  [CLEAR_THREADS](state) {
    state.threadList = null
  },
  [SET_IS_DIALOG_VISIBLE](state, { dialogState }) {
    state.isDialogVisible = !dialogState
  }
}

export const actions = {
  async [LIST_THREADS]({ commit }) {
    const list = await this.$axios.$get('/threads?limit=20')
    commit(SET_THREAD_LIST, { threadList: list })
  },
  async [LIST_THREADS_MORE]({ commit }, { limit, cursor }) {
    const list = await this.$axios.$get(
      `/threads?limit=${limit}&cursor=${cursor}`
    )

    commit(ADD_THREAD_LIST, { threadList: list })
  },
  async [SAVE_THREAD]({ commit }, { payload }) {
    const response = await this.$axios.$post(`/threads`, payload)
    commit(ADD_THREAD, { thread: response })
  },

  async [EDIT_THREAD]({ commit }, { payload }) {
    const response = await this.$axios.$put(`/threads/${payload.id}`, payload)
    commit(UPDATE_THREAD, { thread: response })
  },

  async [DELETE_THREAD]({ commit }, { id }) {
    await this.$axios.$delete(`/threads/${id}`)
    commit(REMOVE_THREAD, id)
  },
  [CHANGE_IS_DIALOG_VISIBLE]({ commit }, { dialogState }) {
    commit(SET_IS_DIALOG_VISIBLE, { dialogState: dialogState })
  }
}
