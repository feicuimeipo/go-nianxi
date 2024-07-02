<template>
  <div>
    <div class="input-group mb-3 search" >
      <div class="search-item"><i class="la la-search" :style="form.iconStyle"></i></div>
      <div class="search-item"><input type="text" id="input-search"  :style="form.inputStyle"  :placeholder="form.placeholder"  :aria-describedby="form.placeholder" @mouseenter="mouseenter"></div>
      <!--搜点出最近文档-->
    </div>
  </div>
</template>

<script lang="ts" setup>

interface Props {
  width: string|null
  placeholder: string|null,
  iconSize: string|null,
}
//搜索文件、文件夹、用户
const props = withDefaults(defineProps<Props>(),{
  width: "350px",
  placeholder: "搜索文件、文件夹、用户",
  iconSize: "20px"
})


const form = ref({
  inputStyle: "width:" + props.width,
  placeholder: props.placeholder,
  iconStyle: "font-size: "+props.iconSize
})


const emit = defineEmits<{
  (e:'search',text:string):void
}>()

const mouseenter = () =>{
  const text = $("#input-search").text();
  emit("search",text);
}


// :style="search"
</script>

<style scoped>
.search{
  border:1px solid var(--bs-border-color);
  margin-top: 3px;
  padding-top: 5px;
  border-radius: 5px;
  height: calc(var(--home-top-panel-height)/1.2);
}

.search-item{
  padding-left: 10px;
}

.search .la{
  font-size: 20px;
  color: var(--bs-gray);
}

.search-item input{
  border: 0px solid var(--bs-border-color);
  outline: none;
  width: 100px;
  background-color: transparent;
}
</style>
