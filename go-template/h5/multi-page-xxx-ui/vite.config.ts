import {defineConfig, loadEnv} from 'vite'
import path,{join,resolve} from 'path'
import { getPlugins } from './build/getPlugins'
import { input } from './build/getPages'


const {VITE_APP_URL} = process.env

export default defineConfig(({mode})=>{
  const plugins = getPlugins(mode)
  return {
      envPrefix: "APP_",
      root: './src/pages',
      publicDir: resolve(__dirname, 'public'),
      base: './',
      plugins,
      define: {
        'process.env': loadEnv(mode, process.cwd())
      },
      resolve: {
        alias: {
            '@': resolve(__dirname, 'src'),
            '@libs': resolve(__dirname, 'public/libs'),
        }
      },
      server: {
        host: 'localhost',        // 指定服务器应该监听哪个 IP 地址
        port: 5173,               // 端口
        strictPort: false, // 若端口已被占用,尝试下移一格端口
        open: false,
        proxy: {
          '/api': {
            target: process.env.VITE_API_URL,
            ws: true,
            changeOrigin: true,
            rewrite: (path) => path.replace(/^\/api/, ''),
          },
        },
      },

      build: {
        outDir: resolve(process.cwd(),'dist'),// 指定输出路径（相对于 项目根目录)
        assetsDir: 'static', // 静态文件目录
        // 默认情况下 若 outDir 在 root 目录下， 则 Vite 会在构建时清空该目录。
        emptyOutDir: true,
        sourcemap: false, // 构建后是否生成 source map 文件
        chunkSizeWarningLimit: 1500, // 规定触发警告的 chunk(文件块) 大小
        minify: 'esbuild',
        rollupOptions: {  // 自定义底层的 Rollup 打包配置
          input,
          output:{
            entryFileNames: 'static/js/[name]-[hash].js',
            chunkFileNames: 'static/js/[name]-[hash].js',
            assetFileNames:'static/[ext]/[name]-[hash].[ext]',
            compact: true,
            manualChunks: (id: string) => {
              if(id.includes("node_modules")) {
                 return id.toString().split('node_modules/')[1].split('/')[0].toString(); // 拆分多个vendors
              }
            }
          },
        },
      },//builder
    }//return
})
