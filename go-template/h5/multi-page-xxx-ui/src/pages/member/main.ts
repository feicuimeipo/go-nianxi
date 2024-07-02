import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from '@/common/store'
//图标-全局引入
import svgIcon from "./components/SvgIcon.vue";
import '@/assets/member/css/index.scss';

const app = createApp(App)
app.config.globalProperties.$=jQuery
app.component('svg-icon', svgIcon)
app.use(router)
    .use(store)
   .mount('#root')
   ;

