<template>
  <div class="application">
    <el-tag class="app-title">应用列表</el-tag>
    <div v-for="item in appsTree" class="app-list">
        <h3> {{ item.desc }}</h3>
      <AppItem :apps="item.children" @changeCurrentAppId="changeCurrentAppId" />
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import variables from '@/styles/variables.scss'
import store from '@/store'
import AppItem from './AppItem'

export default {
  components: {
    AppItem
  },
  data() {
    return {
      appsTree: []
    }
  },
  created() {
    this.getAppsData()
  },
  computed: {
    ...mapGetters([
      'permission_apps',
      'currentAppId'
    ]),
    variables() {
      return variables
    }
  },
  methods: {
    getAppsData() {
      const apps = store.getters.permission_apps
      this.appsTree = apps
    },
    changeCurrentAppId(val) {
      const uname = this.$store.getters.name
      const data = uname + '_' + val
      this.$store.dispatch('user/setCurrentAppId', data)

      setTimeout(() => {
          location.reload()
      }, 1000)
    }
  }

}
</script>

<style scoped>
.application{
  padding: 24px;
  font-size: 14px;
  line-height: 1.5;
  word-wrap: break-word;
}
.img{
  width: 20px;
  height: 20px;
}
.app-list{
  padding-top: 5px;
}

.app-title{
  margin-bottom: 10px;
}
</style>
