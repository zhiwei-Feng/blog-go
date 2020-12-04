import axios from 'axios'
import router from '../router'

axios.defaults.withCredentials = true
export const request = (config) => {
  return axios(config)
}

// 请求前设置header
axios.interceptors.request.use(
  config => {
    if (localStorage.getItem('token')) {
      config.headers.Authorization = `Bearer ${localStorage.getItem("token")}`;
    }

    return config
  },
  error => {
    return Promise.reject(error)
  })
// 请求完成后 拦截器
axios.interceptors.response.use(
  response => {
    console.log(response)
    if (response.data.code === 1003) {
      router.replace({
        path: '/'
      })
      localStorage.removeItem('token')
    }
    return response
  },
  error => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          localStorage.removeItem('token')
          router.replace({
            path: '/',
            query: {redirect: router.currentRoute.fullPath} // 登录成功后 跳转当前页面
          })
      }
    }
  }
)
// 这句一定要写
export default axios
