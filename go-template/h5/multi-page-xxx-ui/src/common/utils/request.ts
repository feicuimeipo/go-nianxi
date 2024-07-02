import axios from 'axios'
import { ElMessageBox, ElMessage  } from 'element-plus'
import store from '@/common/store'
import { getToken } from '@/common/token'
import router from "@/common/router";


const request = axios.create({
    baseURL: import.meta.env.APP_BASE_API,
    timeout: 5000 // request timeout
})


// request interceptor
request.interceptors.request.use(
        config => {
            if (store.getters.token) {
                config.headers['Authorization'] = 'Bearer ' + getToken()
            }
            return config
        },
        error => {
            console.log(error) // for debug
            return Promise.reject(error)
        }
    )

    // response interceptor
request.interceptors.response.use(
    response => {
        if (response.status === 200) {
            return Promise.resolve(response);
        } else {
            return Promise.reject(response);
        }
        // const res = response.data
        // return res
    },
    error => {
        if (error.response.status === 401) {
            if (error.response.data.message.indexOf('JWT认证失败') !== -1) {
                ElMessageBox.confirm(
                    '登录超时, 重新登录或继续停留在当前页？',
                    '登录状态已失效',
                    {
                        confirmButtonText: '重新登录',
                        cancelButtonText: '继续停留',
                        type: 'warning'
                    }
                ).then(() => {
                    store.dispatch('user/logout').then(() => {
                        location.reload()
                    })
                })
            } else {
                ElMessage({
                    showClose: true,
                    message: error.response.data.message,
                    type: 'error',
                    duration: 5 * 1000
                })
                return Promise.reject(error)
            }
        } else if (error.response.status === 403) {
            router.push({path: '/401'})
        } else {
            ElMessage({
                showClose: true,
                message: error.response.data.message || error.message,
                type: 'error',
                duration: 5 * 1000
            })
            return Promise.reject(error)
        }
    },
)


export default request

