# Vue 3 + TypeScript + Vite


# 编辑器
```
对于编辑器 html 和 components/input javascript 模板，请检查 editor.html

对于 CSS 更改，请编辑 scss/editor.scss 和 scss/_builder.scss

```



# 关键文件
```
package.json # 前端项目必备
.gitignore # 一些 git 要忽略的文件和目录
vite.config.js # vite 配置

index.html # spa 入口，可以看到 script 标签使用了 type="module"
src/main.js # vue 实例化，应用从这里启动
src/App.vue # 应用容器
```


# node服务端

- 框架依然使用 koa
- 定义接口需要用到 koa-router
- 解决一下跨域问题 @koa/cors
- 全套安装 ```npm i koa koa-router @koa/cors```


# 后端框架

## koa 与 express
- express 内置了很多中间件，而 koa 可以自由组装，还可以 async/await。
- 启动后端服务：启动后端服务 node server

## vite or vue-cli 
- 相比 vue cli 的编译打包，vite 利用了浏览器原生的 module 加载，速度极快。

## vue 与 react
- vue 和 react 现在基本上分庭抗礼，都是前端必知的 mvvm 框架。

## vue3 