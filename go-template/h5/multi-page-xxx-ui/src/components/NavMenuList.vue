<template>

    <div v-for="item in navMenuList" class="navMenu"  @click="click(item)" >
        <div id="navMenuItem" class="link-btn">
             <i :class=item.icon></i>
             <label class="link-secondary">{{item.name}}</label>
         </div>
         <label class="rightIcon">
            <i :class=item.rightIcon></i>
        </label>
    </div>

</template>

<script setup lang="ts">
import {nextTick, onMounted, ref,defineEmits} from "vue";
import {NavMenu} from "@/components/model/commonModel";
import $ from "jquery";
interface Props {
  navMenuList: NavMenu[]|null
}
const props = defineProps<Props>()
const navMenuList = ref(props.navMenuList)

const emit = defineEmits<{
  (e:'clickMenuItem',menu:NavMenu):void
}>()

const click = (menu: NavMenu) =>{
  emit('clickMenuItem', menu)
}

onMounted(()=>{
  nextTick(()=>{
      $(".navMenu").each(function (i,element){
        const that = this;
        $(element).on("mousemove", function () {
          $(element).find("label").removeClass("link-secondary")
          $(element).addClass("mousemove")
          //$(element).find("span").addClass("mousemove")
        }),
            $(element).on("mouseleave", function () {
              $(element).find("label").addClass("link-secondary")
              $(element).removeClass("mousemove")
              // $(element).css("background-color","red")
            })
      })

  })
})

</script>

<style scoped lang="scss">
.navMenu {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}

.navMenu .la{
  font-size: 18px;
}

#navMenuItem,.rightIcon{
  vertical-align: middle;
  height: 35px;
  letter-spacing: 3px;
}

.mousemove{
  color: var(--nav-menu-font-color) !important;;
  background-color: var(--nav-menu-hover-color) !important;
  border-color: var(--bs-btn-hover-border-color);
}




</style>
