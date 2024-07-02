//const Cookie = require("../js-cookies/js.cookie.min");
if (Nianxi===undefined) {
    var Nianxi = {}
}
Nianxi.Http = {
    xsrfHeaderName : "Authorization", //跨域认证信息 header名
    BASE_URL :"",
    http : axios,
    //认证类型
    AUTH_TYPE : {
        BEARER:'Bearer',
        BASIC:'basic',
        AUTH1:'auth1',
        AUTH2:'auth2'
    },
    METHOD : {
        GET:'get',
        POST:'post',
        PUT:'put',
        DELETE:'delete'
    },

    Init : (baseUrl)=>{
        axios.defaults.timeout = 20000
        axios.defaults.withCredentials = false
        axios.defaults.xsrfHeaderName = Nianxi.xsrfHeaderName
        axios.defaults.xsrfCookieName = Nianxi.xsrfHeaderName
        axios.defaults.baseURL = baseUrl
    },

    //axios请求
     RequestPlugin: async function(url, method, params, config){
        switch(method){
            case Nianxi.http.METHOD.GET:
                return axios.get(url,{params,...config})
            case Nianxi.http.METHOD.POST:
                return axios.post(url,params,config)
            case Nianxi.http.METHOD.PUT:
                return axios.put(url,params,config)
            case Nianxi.http.METHOD.DELETE:
                return axios.delete(url,{params,...config})
            default:
                return axios.get(url,{params,...config})
        }
    },
    SetAuthorization: function (auth,authType=AUTH_TYPE.BEARER){
        switch(authType){
            case Nianxi.http.AUTH_TYPE.BEARER:
                Cookie.set(Nianxi.xsrfHeaderName, 'Bearer '+auth.token,{expires:auth.expireAt})
                break
            case Nianxi.http.AUTH_TYPE.BASIC:
            case Nianxi.http.AUTH_TYPE.AUTH1:
            case Nianxi.http.AUTH_TYPE.AUTH2:
            default:
                break
        }
    },
    RemoveAuthorization: function(authType=AUTH_TYPE.BEARER){
        switch(authType){
            case Nianxi.http.AUTH_TYPE.BEARER:
                Cookie.remove(Nianxi.xsrfHeaderName)
                break
            case Nianxi.http.AUTH_TYPE.BASIC:
            case Nianxi.http.AUTH_TYPE.AUTH1:
            case Nianxi.http.AUTH_TYPE.AUTH2:
            default:
                break
        }
    },
    CheckAuthorization: function (authType=AUTH_TYPE.BEARER){
        switch(authType){
            case Nianxi.http.AUTH_TYPE.BEARER:
                if(Cookie.get(Nianxi.xsrfHeaderName)){
                    return true;
                }
                break
            case Nianxi.http.AUTH_TYPE.BASIC:
            case Nianxi.http.AUTH_TYPE.AUTH1:
            case Nianxi.http.AUTH_TYPE.AUTH2:
            default:
                break
        }
        return false
    },
    //加载axios拦截器
    LoadInterceptors: function (interceptors,options){
        const {request,response} = interceptors
        //加载请求拦截器
        request.forEach(item=>{
            let {onFulfilled,onRejected} = item
            if(!onFulfilled||typeof onFulfilled !=='function'){
                onFulfilled = config =>config
            }
            if(!onRejected||typeof onRejected!=='function'){
                onRejected = error=>Promise.reject(error)
            }
            axios.interceptors.request.use(
                config=>onFulfilled(config,options),
                error=>onRejected(error,options)
            )
        })
        //加载响应拦截器
        response.forEach(item=>{
            let {onFulfilled,onRejected}=item
            if(!onFulfilled||typeof onFulfilled!=='function'){
                onFulfilled = response=>response
            }
            if(!onRejected||typeof onRejected!=='function'){
                onRejected = error=>Promise.reject(error)
            }
            axios.interceptors.response.use(
                response=>onFulfilled(response,options),
                error=>onRejected(error,options)
            )
        })
    }
}
