import {createStore} from 'vuex'
import {getToken, removeToken, setToken} from "@/common/token";
import {getInfo, login, logout, refreshToken} from "@/api/auth/user";
import {resetRouter} from "@/common/router";

const userState = {
    token: getToken(),
    userId:'',
    username: '',
    avatar: '',
    introduction: '',
}

export default createStore({
    state: {
        token: "",
        userId: "",
        username: "",
        avatar: "",
        introduction: "",
    },
    getters:{
        token: state => userState.token,
        userId: state => userState.userId,
        avatar: state => userState.avatar,
        name: state => userState.username,
        introduction: state => userState.introduction,
    },
    mutations: {
        SET_TOKEN: (state:any, token:string) => {
            state.token = token
        },
        SET_INTRODUCTION: (state:any, introduction:string) => {
            state.introduction = introduction
        },
        SET_NAME: (state:any, name:string) => {
            state.name = name
        },
        SET_AVATAR: (state:any, avatar:string) => {
            state.avatar = avatar
        },
        SET_USERID:(state:any,userId:number) =>{
            state.userId = userId
        }
    },
    actions: {
        // user login
        login: ({ commit }:{commit:any}, userInfo:{username: string,password:string}) =>{
            const { username, password } = userInfo
            return new Promise((resolve:any, reject) => {
                login({ username: username.trim(), password: password }).then(response => {
                    const { data } = response
                    commit('SET_TOKEN', data.token)
                    setToken(data.token)
                    resolve()
                }).catch(error => {
                    reject(error)
                })
            })
        },
        // get user info
        getInfo: ({ commit, state }:{commit:any,state: any}) => {
            return new Promise((resolve, reject) => {
                getInfo(state.token).then(response => {
                    const { data } = response
                    if (!data) {
                        reject('Verification failed, please Login again.')
                    }
                    const userInfo = data.userInfo
                    const { userId, username, avatar, introduction } = userInfo

                    commit('SET_NAME', username)
                    commit('SET_AVATAR', avatar)
                    commit('SET_INTRODUCTION', introduction)
                    commit("SET_USERID",userId)
                    resolve(userInfo)

                }).catch(error => {
                    reject(error)
                })
            })
        },

        // user logout
        logout:({ commit, state, dispatch }:{commit:any,state: any,dispatch:any})=> {
            return new Promise((resolve:any, reject) => {
                logout(state.token).then(() => {
                    commit('SET_TOKEN', '')
                    removeToken()
                    resetRouter()

                    resolve()
                }).catch(error => {
                    reject(error)
                })
            })
        },

        // refresh token
        async refreshToken({ commit }:{commit:any}) {
            // 刷新token
            const { data } = await refreshToken()
            commit('SET_TOKEN', data.token)
            setToken(data.token)
        },

        // remove token
        resetToken({ commit }:{commit:any}) {
            return new Promise((resolve:any) => {
                commit('SET_TOKEN', '')
                removeToken()
                resolve()
            })
        },
    },
    modules: {
    }
})
