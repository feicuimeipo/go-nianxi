<template>
      <li v-for="item in menuList"  class="downMenuTree">
          <div v-if="item.child!=null && item.child.length>0 " class="box dropend"><!--左拉菜单-->
                <a class="dropdown-item" href="#" role="button"  aria-expanded="false" data-bs-toggle="dropdown" >
                  <i :class="item.icon"></i><span class="label">{{ item.name }} </span><span class="arrow">></span>
                </a>
                <ul class="dropdown-menu" :id="'dropdown-menu-'+item.Id" data-popper-placement="right-start">
                   <TopMenuChild :parent-id=item.id :menuList=item.child @clickMenuItem="clickMenuItem"/>
                </ul>
          </div>
        <div v-else-if="(item.divider && (item.name!='#' && item.name!=''))">
            <hr  class="dropdown-divider" />
           <label class="dividerLabel">{{item.name}}</label>
        </div>
        <hr v-else-if="(item.divider)" class="dropdown-divider" />
        <a v-else class="dropdown-item" href="#"  @click="click(item)">
          <i :class="item.icon"></i> {{ item.name }}
        </a>
      </li>
</template>

<script setup lang="ts">
import {defineEmits, nextTick, onMounted, ref} from 'vue'
import  TopMenuChild from "./DownMenuTreeChild.vue"
import {DownMenuTree} from "@/components/model/commonModel";
import $ from "jquery";

// <label class="addFileLabel">创建文件</label>-->

interface Props {
  menuList: DownMenuTree[]
}
const props = defineProps<Props>()
const menuList = ref(props.menuList)

const emit = defineEmits<{(e:'clickMenuItem',menu:DownMenuTree):void}>()

const click = (menu: DownMenuTree) =>{
  emit('clickMenuItem', menu)
}

const clickMenuItem = (menu: DownMenuTree) =>{
  emit('clickMenuItem', menu)
}


onMounted(()=>{
  nextTick(()=>{
    $(".downMenuTree").find("a").each(function (i,element){
          $(element).on("mousemove", function () {
            $(element).find("span").removeClass("link-secondary")
            $(element).addClass("mousemove")
            //$(element).find("span").addClass("mousemove")
          }),
          $(element).on("mouseleave", function () {
            $(element).find("span").addClass("link-secondary")
            $(element).removeClass("mousemove")
            // $(element).css("background-color","red")
          })
    })

  })
})
</script>

<style scoped>

.downMenuTree .la{
  font-size: 15px;
}

 .box {
   width: 100%;
   text-align: justify;
   padding: 0px;
 }

 .box:after {
   width: 100%;
   display: inline-block;
   overflow: hidden;
   padding: 0px
 }
 .box span {
   display: inline-block;
   text-align: left;
   width: 95%;
 }

.mousemove{
  color: var(--nav-menu-font-color) !important;;
  background-color: var(--nav-menu-hover-color) !important;
  border-color: var(--bs-btn-hover-border-color);
}

.dividerLabel{
  font-size: small;
  color: #CCCCCC;
  padding: 2px;
}



</style>
