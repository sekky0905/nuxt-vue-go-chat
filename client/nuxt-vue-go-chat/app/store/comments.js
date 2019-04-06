import {
  SET_COMMENTS,
  SET_COMMENT_LIST,
  ADD_COMMENT_LIST,
  ADD_COMMENT,
  UPDATE_COMMENT,
  REMOVE_COMMENT,
  CLEAR_COMMENTS
} from './mutation-types'

import {
  LIST_COMMENTS,
  LIST_COMMENTS_MORE,
  SAVE_COMMENT,
  EDIT_COMMENT,
  DELETE_COMMENT
} from './action-types'

export const state = () => ({
  commentList: {
    comments: [],
    hasNext: false,
    cursor: ''
  }
})

export const getters = {
  commentList: state => state.commentList,
  comments: state => state.commentList.comments
}

export const mutations = {
  [SET_COMMENTS](state, { comments }) {
    state.commentList.comments = comments
  },
  [SET_COMMENT_LIST](state, { commentList }) {
    state.commentList = commentList
  },
  [ADD_COMMENT_LIST](state, { commentList }) {
    state.commentList.comments = state.commentList.comments.concat(
      commentList.comments
    )
    state.commentList.hasNext = commentList.hasNext
    state.commentList.cursor = commentList.cursor
  },
  [ADD_COMMENT](state, { comment }) {
    state.commentList.comments.push(comment)
  },
  [UPDATE_COMMENT](state, { comment }) {
    state.commentList.comments = state.commentList.comments.map(t =>
      t.id === comment.id ? comment : t
    )
  },
  [REMOVE_COMMENT](state, { id }) {
    state.commentList.comments = state.commentList.comments.filter(t => {
      return t.id !== id
    })
  },
  [CLEAR_COMMENTS](state) {
    state.commentList = {
      comments: [],
      hasNext: false,
      cursor: ''
    }
  }
}

export const actions = {
  async [LIST_COMMENTS]({ commit }, { threadId }) {
    console.log(`threads/${threadId}/comments`)
    commit(CLEAR_COMMENTS)
    const list = await this.$axios.$get(`threads/${threadId}/comments`)
    if (!list.comments) {
      return
    }
    commit(SET_COMMENT_LIST, { commentList: list })
  },
  async [LIST_COMMENTS_MORE]({ commit }, { threadId, limit, cursor }) {
    const list = await this.$axios.$get(
      `/threads/${threadId}/comments?limit=${limit}&cursor=${cursor}`
    )
    if (!list.comments) {
      return
    }
    commit(ADD_COMMENT_LIST, { commentList: list })
  },
  async [SAVE_COMMENT]({ commit }, { payload }) {
    const response = await this.$axios.$post(
      `threads/${payload.threadId}/comments`,
      payload
    )
    commit(ADD_COMMENT, { comment: response })
  },
  async [EDIT_COMMENT]({ commit }, { payload }) {
    const response = await this.$axios.$put(
      `threads/${payload.threadId}/comments/${payload.id}`,
      payload
    )
    commit(UPDATE_COMMENT, { comment: response })
  },

  async [DELETE_COMMENT]({ commit }, { threadId, id }) {
    await this.$axios.$delete(`threads/${threadId}/comments/${id}`)
    commit(REMOVE_COMMENT, { id: id })
  }
}
