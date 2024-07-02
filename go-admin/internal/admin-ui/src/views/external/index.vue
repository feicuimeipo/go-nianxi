<template>
  <el-main class="content">
    <!-- <keep-alive>
      <router-view />
    </keep-alive> -->
    <!--    style='position:absolute;z-index:4;'-->
    <iframe
      :id="name"
      :ref="name"
      :src="url"
      :name="name"
      @load="loading = false"
      allowTransparency="true"
      width="100%"
      height="100%"
      allowFullscreen="true"
      class="YuFrame1"
      frameborder="0"
      scrolling="no"
    />
  </el-main>
</template>

<script>

export default {
  name: 'Content',
  data() {
    return {
      clientHeight: 500
    }
  },
  computed: {
    cachedViews() {
      return this.$store.state.tagsView.cachedViews
    },
    key() {
      return this.$route.path
    },
    url() {
      return this.$route.meta.iframeUrl
    },
    name() {
      return this.$route.meta.name
    }
  },
  created() {
    // this.createIframe()
  },
  // methods: {
  //   // createIframe(){}
  //   setMyIframeHeight() {
  //     // 页面可视高-顶部实际高-底部实际高-多余部分高度
  //     const that = this
  //     this.$nextTick(() => {
  //       const headerHeight = document.getElementById('header').offsetHeight
  //       const footerHeight = document.getElementById('footer').offsetHeight
  //       const bodyHeight = document.body.clientHeight
  //       const outHeight = parseFloat(bodyHeight - headerHeight - footerHeight) - 5 + 'px'
  //       document.querySelector('#myIframe').style.height = outHeight
  //     })
  //   }
  // }
}
</script>

<style scoped>
.YuFrame1{
  width: 100vw;
  height: 100vh;
  top: 0;
  left: 0;

}
</style>
