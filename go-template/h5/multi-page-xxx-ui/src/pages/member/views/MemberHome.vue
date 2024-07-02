<template xmlns="http://www.w3.org/1999/html">
  <div id="home">
      <div id="left-panel">
          <div id="btn-panel" >
              <div class="float-start">
                  <a href="#" class="btn" role="button" onclick="window.location='../home/index.html'">
                    <i class="la la-home"></i>
                  </a>
              </div>
              <div class="float-end">
                <a href="#" class="btn" role="button" >
                  <i class="la la-concierge-bell"></i>
                </a>
              </div>
          </div>

          <div id="user-panel">
              <div style="vertical-align: middle;padding-bottom: 10px">
                  <i style="font-size: xxx-large;" class="la la-user-alt"></i>
                  <label style="font-size: medium;" >Nianxl</label>
              </div>
          </div>
          <div id="menu-panel">
            <NavMenuList :navMenuList="navMenuList" @clickMenuItem="clickMenuItem" />
          </div>
      </div>
      <div id="main-panel">


      </div>
     <!--home-->
     </div>

</template>

<script setup lang="ts">
import  NavMenuList from "@/components/NavMenuList.vue"

import { NavMenu} from "@/components/model/commonModel";
import SearchBar from "@/components/SearchBar.vue";


const navMenuList = ref<NavMenu[]>( [
  {id:1,name: "会员信息",router: 'recently', icon:'la la-paw',rightIcon:null},
  {id:2,name: "帐户安全",router: 'recently', icon:'la la-portrait',rightIcon:null},
  {id:3,name: "订单中心",router: 'recently', icon:'la la-shapes',rightIcon:null},
])

const clickMenuItem = (menu:NavMenu) => {
  console.log("子子组件中传来菜单点击事件："+menu.name)
}

const search = (text: string) =>{
  console.log("search="+text);
}


const addFileMenu = () =>{
    $("#down-menu-add-file").toggle()
}
</script>

<style scoped>
#home{
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  border-style: solid;
  border-width: 0px;
}


#left-panel {
  float: left;
  position: fixed;
  border-right-width: 1px;
  width: 220px;
  width: calc(var(--home-left-panel-width));
  height: 100%;
  max-height: 100%;
  border-color: var(--bs-border-color);
  background-color: var(--bs-secondary-bg-subtle);
  float: left;
  padding: 0px 10px 0px 10px;
}

#left-panel #btn-panel{
  width: 100%;
  float: top;
  text-align: left;
  padding: 0px 0px 10px 0px;
}

#btn-panel .la{
  font-size: large
}

#user-panel {
  padding: 30px 0px 0px 0px;
}

#left-panel #menu-panel{
  padding: 30px 0px 0px 0px;
  width: 100%;
}




#main-panel {
  margin: 10px;
  padding: 10px 10px 10px calc(var(--home-left-panel-width));
  background-color: var(--bs-body-bg);
  top: var(--home-top-panel-height);
  right:  var(--home-right-panel-width);
  height: 100%;
  width: 100%;
  float: right;
}



</style>
