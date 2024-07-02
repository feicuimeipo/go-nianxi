

# 代码
- 链码更新
- ```azure
go get gitee.com/nianxiaoling/go-nianxi/fabric-chaincode/internal/ccsdkapi/user
```

最佳实践
https://zhuanlan.zhihu.com/p/514426016
https://github.com/sdgmf/go-project-sample

# 参考s
### 某blog的目录结构
> xxx-blog/├── conf├── middleware├── models├── pkg├── routers└── runtime
- conf：用于存储配置文件
- middleware：应用中间件
- models：应用数据库模型
- pkg：第三方包
- routers 路由逻辑处理
- runtime 应用运行时数据


### web-ui-mini
```azure
├─common # casbin mysql zap validator 等公共资源
├─config # viper读取配置
├─controller # controller层，响应路由请求的方法
├─dto # 返回给前端的数据结构
├─middleware # 中间件
├─model # 结构体模型
├─repository # 数据库操作
├─response # 常用返回封装，如Success、Fail
├─routes # 所有路由
├─util # 工具方法
└─vo # 接收前端请求的数据结构
```

### 标准工程结构
```azure
/cmd
/internal
/pkg
/API
/WEB
/configs
/scripts
/build
/deployments
/test
/docs
```

###### /deployments

IaaS、PaaS、系统和容器编排部署配置和模板(docker-compose、kubernetes/helm、mesos、terraform、bosh)。注意，在一些存储库中(特别是使用 kubernetes 部署的应用程序)，这个目录被称为 /deploy。

###### /test

###### /docs
设计和用户文档(除了 godoc 生成的文档之外)。


