'use strict'
const webpack = require('webpack')
const CompressionWebpackPlugin = require('compression-webpack-plugin')
let pages = require('./build/getPages.js')
const pagesArray = require('./build/pages.js')
const productionGzipExtensions = ['js', 'css']
const path = require('path')

const defaultSettings = require('./build/settings.js')
const title = defaultSettings.title || 'Go Nianxi' // page title

function resolve(dir) {
  return path.join(__dirname, dir)
}

// If your port is set to 80,
// use administrator privileges to execute the command line.
// For example, Mac: sudo npm run
// You can change the port by the following method:
// port = 8080 npm run dev OR npm run dev --port = 8080
const port = process.env.port || process.env.npm_config_port || 8080 // dev port

// All configuration item explanations can be find in https://cli.vuejs.org/config/
module.exports = {
  /**
   * You will need to set publicPath if you plan to deploy your site under a sub path,
   * for example GitHub Pages. If you plan to deploy your site to https://foo.github.io/bar/,
   * then publicPath should be set to "/bar/".
   * In most cases please use '/' !!!
   * Detail: https://cli.vuejs.org/config/#publicpath
   */
  runtimeCompiler: true,
  publicPath: '/',
  pages,
  productionSourceMap: false,
  // context: path.resolve(__dirname, '../'),
  // outputDir: 'dist',
  // indexPath: 'index.html', // 指定生成的index.html输出路径
  // assetsDir: 'static',
  lintOnSave: process.env.NODE_ENV === 'development',
  devServer: {
    port: port,
    open: false,
    overlay: {
      warnings: false,
      errors: true
    },
    proxy: {
      '/api': {
        target: process.env.VUE_APP_BASE_API,
        changeOrigin: true,
        pathRewrite: {
          '^/api': 'api'
        }
      }
    }
  },
  configureWebpack: {
    // provide the app's title in webpack's name field, so that
    // it can be accessed in index.html to inject the correct title.
    name: title,
    resolve: {
      alias: {
        '@': resolve('src')
      }
    },
    plugins: [
      new CompressionWebpackPlugin({
        algorithm: 'gzip',
        test: new RegExp('\\.(' + productionGzipExtensions.join('|') + ')$'),
        threshold: 10240,
        minRatio: 0.8
      }),
      new webpack.optimize.LimitChunkCountPlugin({
        maxChunks: 5,
        minChunkSize: 100
      })
    ],
    optimization: {
      splitChunks: {
        chunks: 'all',
        minSize: 20000, // 允许新拆出 chunk 的最小体积，也是异步 chunk 公共模块的强制拆分体积
        maxAsyncRequests: 6, // 每个异步加载模块最多能被拆分的数量
        maxInitialRequests: 6, // 每个入口和它的同步依赖最多能被拆分的数量
        enforceSizeThreshold: 50000, // 强制执行拆分的体积阈值并忽略其他限制
        cacheGroups: {
          libs: {
            // 第三方库
            name: 'chunk-libs',
            test: /[\\/]node_modules[\\/]/,
            priority: 10
            // chunks: "initial" // 只打包初始时依赖的第三方
          },
          vant: {
            // vant 单独拆包
            name: 'chunk-vant',
            test: /[\\/]node_modules[\\/]vant[\\/]/,
            priority: 20 // 权重要大于 libs
          },
          jsencrypt: {
            // jsencrypt 单独拆包
            name: 'chunk-jsencrypt',
            test: /[\\/]node_modules[\\/]jsencrypt[\\/]/,
            priority: 30 // 权重要大于 libs
          },
          elementUI: {
            name: 'chunk-elementUI', // split elementUI into a single package
            priority: 20, // the weight needs to be larger than libs and app or it will be packaged into libs or app
            test: /[\\/]node_modules[\\/]_?element-ui(.*)/ // in order to adapt to cnpm
          },
          // svgIcon: {
          //   // svg 图标
          //   name: 'chunk-svgIcon',
          //   test(module) {
          //     // `module.resource` 是文件的绝对路径
          //     // 用`path.sep` 代替 / or \，以便跨平台兼容
          //     // const path = require('path') // path 一般会在配置文件引入，此处只是说明 path 的来源，实际并不用加上
          //     return (
          //       module.resource &&
          //       module.resource.endsWith('.svg') &&
          //       module.resource.includes(`${path.sep}icons${path.sep}`)
          //     )
          //   },
          //   priority: 30
          // },
          commons: {
            // 公共模块包
            name: `chunk-commons`,
            minChunks: 2,
            priority: 0,
            reuseExistingChunk: true
          }
        }
      }
    }
  },
  chainWebpack(config) {
    // config.plugin('preload').tap(() => [
    //   {
    //     rel: 'preload',
    //     fileBlacklist: [/\.map$/, /hot-update\.js$/, /runtime\..*\.js$/],
    //     include: 'initial'
    //   }
    // ])

    // when there are many pages, it will cause too many meaningless requests
    // config.plugins.delete('prefetch')
    pagesArray.forEach(item => {
      let icons = `src/views/${item.pagePath}/icons`
      if (item.pagePath === 'index') {
        icons = `src/icons`
      }
      config.module.rule('svg').exclude.add(resolve(icons)).end()
      config.module
        .rule('icons')
        .test(/\.svg$/)
        .include.add(resolve(icons))
        .end()
        .use('svg-sprite-loader')
        .loader('svg-sprite-loader')
        .options({
          symbolId: 'icon-[name]'
        })
        .end()
    })
    config
      .when(process.env.NODE_ENV !== 'development',
        config => {
          config
            .plugin('ScriptExtHtmlWebpackPlugin')
            .after('html')
            .use('script-ext-html-webpack-plugin', [{
            // `runtime` must same as runtimeChunk name. default is `runtime`
              inline: /runtime\..*\.js$/
            }])
            .end()
          // config
          //   .optimization.splitChunks({
          //     chunks: 'all',
          //     minSize: 20000, // 允许新拆出 chunk 的最小体积，也是异步 chunk 公共模块的强制拆分体积
          //     maxAsyncRequests: 6, // 每个异步加载模块最多能被拆分的数量
          //     maxInitialRequests: 6, // 每个入口和它的同步依赖最多能被拆分的数量
          //     enforceSizeThreshold: 50000, // 强制执行拆分的体积阈值并忽略其他限制
          //     cacheGroups: {
          //       libs: {
          //         name: 'chunk-libs',
          //         test: /[\\/]node_modules[\\/]/,
          //         priority: 10,
          //         chunks: 'initial' // only package third parties that are initially dependent
          //       },
          //       vant: {
          //         // vant 单独拆包
          //         name: 'chunk-vant',
          //         test: /[\\/]node_modules[\\/]vant[\\/]/,
          //         priority: 20 // 权重要大于 libs
          //       },
          //       jsencrypt: {
          //         // jsencrypt 单独拆包
          //         name: 'chunk-jsencrypt',
          //         test: /[\\/]node_modules[\\/]jsencrypt[\\/]/,
          //         priority: 30 // 权重要大于 libs
          //       },
          //       elementUI: {
          //         name: 'chunk-elementUI', // split elementUI into a single package
          //         priority: 20, // the weight needs to be larger than libs and app or it will be packaged into libs or app
          //         test: /[\\/]node_modules[\\/]_?element-ui(.*)/ // in order to adapt to cnpm
          //       },
          //       commons: {
          //         name: 'chunk-commons',
          //         test: resolve('src/components'), // can customize your rules
          //         minChunks: 3, //  minimum common number
          //         priority: 5,
          //         reuseExistingChunk: true
          //       }
          //     }
          //   })
          // https:// webpack.js.org/configuration/optimization/#optimizationruntimechunk
          // config.optimization.runtimeChunk('single')
          config.optimization.minimizer('terser').tap(args => {
            args[0].terserOptions.compress.drop_console = true
            args[0].terserOptions.compress.warnings = false
            args[0].terserOptions.compress.drop_debugger = true
            return args
          })
        }
      )
  }
}
