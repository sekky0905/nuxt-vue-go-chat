import Vue from 'vue'

export default ({ app, $axios, redirect }) => {
  $axios.onError(error => {
    if (error.response.status === 401) {
      // store.commit('setUser', { user: null, isLoggedIn: false })
      redirect('/')
    }
    return Promise.reject(error)
  })
}
