
import { createRouter, createWebHistory } from 'vue-router'
export const constantRoutes = [
    {
        path: "/",
        name: "default",
        redirect:  '/login'
    },
    {
        path: '/login',
        component: () => import('@/common/view/LoginPage.vue'),
        hidden: true
    },
    {
        path: '/404',
        component: () => import('@/common/view/ErrorPag404.vue'),
        hidden: true
    },
    {
        path: '/401',
        component: () => import('@/common/view/ErrorPag401.vue'),
        hidden: true
    }
]

/**
 * asyncRoutes
 * the routes that need to be dynamically loaded based on user roles
 */
const myCreateRouter  = () => createRouter({
    history: createWebHistory(), //设置为history模式
    // history: createWebHashHistory(),// Hash模式
    routes: constantRoutes
})

const router = myCreateRouter()

export const resetRouter = () => function (){
    router.beforeEach((to, from, next) => {
        next(`/login?redirect=${to.fullPath}`)
    })
}


export default router
