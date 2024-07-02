import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import '@/assets/home/css/index.scss';
//图标-全局引入
import svgIcon from "./components/SvgIcon.vue";
const app = createApp(App)
app.config.globalProperties.$=jQuery
app.component('svg-icon', svgIcon)
app.use(router)
   .mount('#app')
   ;

