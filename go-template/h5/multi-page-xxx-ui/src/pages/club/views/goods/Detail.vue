 <script setup lang="ts" name="product">
    import { useRouter, useRoute } from "vue-router";
    import { computed, reactive } from "vue";
    
    // defineProps,父组件传值入子组件 - 带默认值
    const props = withDefaults(defineProps<{id:number}>(),{
        id: -1,
    }) 
    
    
    // useRouter()方法执行后，返回当前项目中的路由器对象
    let $router = useRouter();   
    // useRoute()方法执行后，返回当前路由信息
    let $route = useRoute();
    console.log("id:=", $route.params.id)
    console.log("props。id:=",props.id)
    // 商品数组
    let goodsList = reactive([
      {
        id: 1,
        name: "手机",
        price: 5999,
        color: "白色",
      },
      {
        id: 2,
        name: "电脑",
        price: 4999,
        color: "红色",
      },
      {
        id: 3,
        name: "电视",
        price: 3999,
        color: "黑色",
      },
    ]);
    
    // 指定的商品对象
    const goods = computed(() => {     
       //return goodsList.find((r) => r.id == props.id);  
        return goodsList.find((r) => r.id == $route.params.id);  
    });
    
    

    const goBack = () =>{
      // 返回列表页的三种方式
      // $router.push('/home')
      // $router.back()
      $router.go(-1);
    }
 </script>

<template>
  <div class="goods">
 <h2>商品信息 </h2>
 <ul>
   <li>商品编号：{{ goods?.id }}</li>
   <li>商品名称：{{ goods?.name }}</li>
   <li>商品价格：{{ goods?.price }}</li>
   <li>商品颜色：{{ goods?.color }}</li>
   <li>
     <button @click=" goBack();">返回</button>
   </li>
 </ul>
</div>
</template>

 <style scoped>
 </style>