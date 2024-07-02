

const path = require("path");

module.exports = {
  babel: {
    plugins: [["@babel/plugin-proposal-decorators", { legacy: true }]]
  },
  // webpack 配置
  webpack: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
    },
    configure: (webpackConfig,{env,paths}) =>{
      // 修改build的生成文件名称
      paths.appBuild = 'dist';
      webpackConfig.output ={
        ...webpackConfig.output,
        path:path.resolve(__dirname,'dist'),
        publicPath: '/'
      }
      return webpackConfig;
    }
  },
  devServer: {
    open: false,
    proxy: {
        "/api": {
          target: 'http://localhost:8000',
          changeOrigin: true,
          secure:false,
          pathRewrite: {
            "^/api": ""
          }
        }
      }
  }
}
