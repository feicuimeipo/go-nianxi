<template>
  <div>
    <div v-for="item in apps" class="app-list">
      <a  href="#" @click="changeCurrentAppId(item.ID)">
        <el-row :gutter="20">
          <el-col :span="1">
               <img  src='@/assets/sidebar-logo/app.svg' width="20px" height="20px" style="pa" />
          </el-col>
          <el-col :span="15"> {{ item.desc }}</el-col>
        </el-row>
      </a>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AppItem',
  props: {
    apps: []
  },
  methods: {
    changeCurrentAppId(val) {
      this.$emit('changeCurrentAppId', val)
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
  padding-top: 10px;
}

.app-title{
  margin-bottom: 10px;
}
</style>
