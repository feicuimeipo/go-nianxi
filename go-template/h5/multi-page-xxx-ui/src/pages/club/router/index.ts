
//1. 引用创建路由模型
import {createRouter,createWebHashHistory,RouteRecordRaw} from 'vue-router'

//默认的开始界面
//①利用利用redirect单独写一个路由
//②router的children的初始页面，可以用redirect定义初始页面,也可以按①的方式
//2.配置系统所有路由页面
const routes: Array<RouteRecordRaw> =[
    {
        path: '/',
        redirect:'home'
    },
    {
        path: '/home',
        name: 'home',
        component: () => import('../views/ClubHome.vue')
    },
   ]

//3. 创建路由实例
const router = createRouter({
    history: createWebHashHistory(), //使用Hash模型
    routes
})

//4. 声明，为路由提供外部引用的入口 ,不要写成{router}否则会警告
export default router
