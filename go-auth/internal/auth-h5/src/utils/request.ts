import axios from 'axios';
import {getToken} from "@/api/auth/token";
import {logout} from "@/api/auth/auth-provider";


const request = axios.create({
    baseURL: process.env.REACT_APP_API_URL,
    withCredentials: true,
    timeout: 50000,
    headers: {
       'X-Requested-With':'XMLHttpRequest',
       'Content-Type': 'application/json; charset=UTF-8'
    }
})


// request interceptor
request.interceptors.request.use(
    config => {
        if (getToken()) {
            config.headers['Authorization'] = 'Bearer ' + getToken()
        }
        return config
    },
    error => {
        console.log(error) // for debug
        return Promise.reject(error)
    }
)



request.interceptors.response.use(
  response => {
      if (response.data && response.data.code && (response.data.code !==0 && response.data.code!=="0")){
          console.log("\n 5.请求异常",JSON.stringify(response.data))
          return Promise.reject(response.data);
      }else {
          const res = response.data
          return Promise.resolve(res);
      }
  },
    error => {
        if (error.response) {
            const data = {
                code: error.response.code ? error.response.code : error.response.status,
                status: error.response.status,
                msg: error.response.data.msg ? error.response.data.msg : error.response.message
            }
            if (error.response.status && (error.response.status === 401 )) {
                window.location.reload()
                return Promise.reject(error.response.data);
            }else if (error.response.status && (error.response.status === 403 || error.response.status === 402)){
                console.log("\n 6 .请求异常", JSON.stringify(error))
                return Promise.reject(data);
            }else {
                console.log("\n 3 .请求异常", JSON.stringify(error))
                return Promise.reject(data);
            }
        }else{
            const data = {
                code: error.status,
                status: error.status,
                msg: error.message
            }
            //console.log("\n 3.1 .请求异常", JSON.stringify(data))s
            console.log("\n 4 .请求异常", JSON.stringify(error))
            return Promise.reject(data);
        }
    }
)

export default request
