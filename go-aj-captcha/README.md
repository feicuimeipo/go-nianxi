# go-aj-captcha

#### 介绍
go-aj-captcha
滑动验证码


####  后端引入
> 配置信息读入vipper,也可以没有，但是vipper不可以为空, 后期做处理
- 初始化
```go
func NewCaptcha(redis redis.UniversalClient, v *viper.Viper, logger *zap.Logger) *captcha.AJCaptcha {
	//path := getCurrentAbPathByCaller()
	resource := ""
	//resource := filepath.Join(path, "pkg", "captcha")
	//resource = filepath.Clean(resource)
	if redis != nil {
		captcha := captcha.NewAJCaptchaByRDB(cache.GetRedisClient(), v, logger, resource)
		Conf.Captcha = captcha
	} else {
		captcha := captcha.NewAJCaptcha(v, logger, resource)
		Conf.Captcha = captcha

	}

	return Conf.Captcha
}
```
- 路由
```azure

var captcha *captcha.AJCaptcha
// 二维码
func InitCaptchaRoutes(r *gin.RouterGroup, c *controller.CaptchaController) gin.IRoutes {
	r.POST("/captcha/get", GetCaptcha)
	r.POST("/captcha/check", c.CheckCaptcha)
return r
    }


//	@Tags		认证
//	@summary	获得验证码
func GetCaptcha(c *gin.Context) {
	writer := c.Writer
	request := c.Request
	captcha.HandleFunc.GetCaptcha(writer, request)
return
    }

//	@Tags		认证
//	@summary	验证验证码验证
func  CheckCaptcha(c *gin.Context) {
	writer := c.Writer
	request := c.Request
	captcha.HandleFunc.CheckCaptcha(writer, request)
return
    }

```
#### 前端引入
> 例子见web, jquery实现
> web 文件夹下lib除了验证码与jquery之外，其他均可不用
```
<script type="text/javascript" src="./libs/base/jquery.min.js" ></script>
<script src="libs/aj-captcha/js/crypto-js.js"></script>
<script src="libs/aj-captcha/js/ase.js"></script>
<script src="libs/aj-captcha/js/verify.js" ></script>

<script>
         //后端地址
          baseUrl = http://www.../ 
           //pointsVerify点选式 改为slideVerify则为滑动试
		$('#captchaPanel').pointsVerify({
			baseUrl: baseUrl,
			mode:'pop', //弹出式 ,fixed为嵌入式， 
			containerId:'loginBtn',
			imgSize : {
				width: '400px',
				height: '200px',
			},
			barSize:{
				width: '400px',
				height: '40px',
			},
			beforeCheck:function(){
				form.on('submit(login)', function(data) {

				})
				const name = $("#username").val();
				const pass = $('#password').val();
				if (name == '' || pass == '') {
					return false
				}
				return true;
			},
			ready : function() {
			},  //加载完毕的回调
			success : function(params) { //成功的回调
				console.log("验证码成功！="+params)
				//params为返回的二次验证参数 需要在接下来的实现逻辑回传服务器
				//--
				//TODO: 回传服务器

			},//验证码
			error : function() {
				return false
			} //失败的回调
		});
```


