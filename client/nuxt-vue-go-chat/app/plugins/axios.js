export default ({ app, $axios, redirect }) => {
  $axios.onError(error => {
    if (error.response.status === 401) {
      console.error(`failed to authenticate: ${JSON.stringify(error)}`)
      redirect('/')
    }
    return Promise.reject(error)
  })
}
