import Vue from 'vue'
// Lib imports
import axios from 'axios'
import Cookies from 'js-cookie'

Vue.prototype.$http = axios.create()
Vue.prototype.axios = Vue.prototype.$http

Vue.prototype.$http.interceptors.request.use(function (request) {
  if (request.url.startsWith('/api') || request.url.startsWith('/daemon')) {
    request.headers[request.method].Authorization = 'Bearer ' + Cookies.get('puffer_auth') || ''
  }
  return request
}, function (error) {
  return Promise.reject(error)
})

Vue.prototype.$http.interceptors.response.use(function (response) {
  return response
}, function (error) {
  if (((error || {}).response || {}).status === 401) {
    localStorage.setItem('reauth_reason', 'session_timed_out')
    Cookies.remove('puffer_auth')
    window.location = '/auth/login'
  }
  return Promise.reject(error)
})
