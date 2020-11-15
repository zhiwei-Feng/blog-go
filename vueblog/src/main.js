import Vue from 'vue'
import App from './App'
import router from './router'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
// import './styles/element-variables.scss'
import 'font-awesome/css/font-awesome.min.css'
import './utils/filter_utils.js'
import axios from 'axios'

Vue.use(ElementUI)
Vue.config.productionTip = false
window.bus = new Vue();
new Vue({
  el: '#app',
  router,
  template: '<App/>',
  components: {App}
})
axios.interceptors.request.use(
  config => {
    if (localStorage.getItem("token")) {  // 判断是否存在token，如果存在的话，则每个http header都加上token
      config.headers.Authorization = `Bearer ${localStorage.getItem("token")}`;
    }
    return config;
  },
  err => {
    return Promise.reject(err);
  });
