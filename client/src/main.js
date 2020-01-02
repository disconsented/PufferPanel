// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
// Components
import './components'
import Vuetify from 'vuetify'
// Plugins
import './plugins'
// Application imports
import App from './App'
import i18n from '@/i18n'
import router from '@/router'
import VueNativeSocket from 'vue-native-websocket'
import vuetify from './plugins/vuetify'

Vue.use(Vuetify)
Vue.use(VueNativeSocket, 'ws://localhost:1234', {
  connectManually: true,
  reconnection: true,
  reconnectionDelay: 5000,
  format: 'json'
})

Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  i18n,
  router,
  vuetify,
  render: h => h(App)
}).$mount('#app')
