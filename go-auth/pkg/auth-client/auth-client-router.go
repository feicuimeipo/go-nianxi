package auth_client

import (
	auth_common "gitee.com/go-nianxi/go-auth/pkg/auth-client/auth-common"
	"gitee.com/go-nianxi/go-common/pkg/ecode/resp"
	"github.com/gin-gonic/gin"
	"time"
)

func (authClient *ClientContext) InitAuthClientRouter(r *gin.RouterGroup) {
	r.GET("/auth/client/logout", authClient.AuthLogout)
}

func (authClient *ClientContext) AuthLogout(ctx *gin.Context) {
	token := ctx.GetHeader(auth_common.TokenHeaderName)
	cookies, _ := ctx.Cookie(auth_common.TokenCookiesName)

	if token == "" {
		token = cookies
	}
	if cookies != "" {
		authClient.conf.Logger.Info("登录！")
	}
	ctx.Header(auth_common.TokenHeaderName, "")
	ctx.SetCookie(auth_common.TokenCookiesName, token, time.Now().Minute(), "/", authClient.conf.O.Domain, false, true)
	resp.Success(ctx, nil)
}
