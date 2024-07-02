package auth_client

import (
	"bytes"
	"encoding/base64"
	"errors"
	"gitee.com/go-nianxi/go-auth/pkg/auth-client/auth-common"
	client_config "gitee.com/go-nianxi/go-auth/pkg/auth-client/client-config"
	"gitee.com/go-nianxi/go-common/pkg/ecode/resp"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type redirectDTO struct {
	loginURl string
}

func (authClient *ClientContext) goLoginPage(c *gin.Context) {
	protocol := "http://"
	if c.Request.TLS != nil {
		protocol = "https://"
	}

	loginUrl := authClient.conf.O.Client.LoginUrl
	returnURl := protocol + c.Request.Host + c.Request.URL.RequestURI()
	loginUrl = loginUrl + "?redirectUrl=" + returnURl

	ajax := c.GetHeader(auth_common.X_Requested_With_HEADER)
	if ajax == "" || ajax != auth_common.X_Requested_With_VALUE {
		c.Redirect(http.StatusTemporaryRedirect, loginUrl)
	} else {
		resp.WithHttpStatus(c, http.StatusUnauthorized, redirectDTO{loginURl: loginUrl}, "JWT认证失败")
	}
}

func (authClient *ClientContext) handlerQuerySign(ctx *gin.Context) {
	queryToken := ctx.Query(auth_common.TokenQueryName)
	header := ctx.GetHeader(auth_common.TokenHeaderName)
	cookies, _ := ctx.Cookie(auth_common.TokenCookiesName)
	if queryToken != "" {
		tokenByte, _ := base64.URLEncoding.DecodeString(queryToken)
		bearToken := auth_common.TokenHeaderNamePrefix + string(tokenByte)

		flag := ctx.Query("flag")
		if flag != "" {
			flagByte, _ := base64.URLEncoding.DecodeString(flag)
			flag = string(flagByte)
		}

		sign := ctx.Query("sign")
		byteSign, _ := base64.URLEncoding.DecodeString(sign)
		sign = string(byteSign)

		publicKey := client_config.GenPublic(authClient.conf.O.Jwt.Secret)
		if bytes.Equal([]byte(sign), []byte(publicKey)) || sign == publicKey {
			if flag == auth_common.FLAG_LOOGIN || (header == "" && cookies == "") {
				ctx.Header(auth_common.TokenHeaderName, bearToken)
				ctx.SetCookie(auth_common.TokenCookiesName, bearToken, int(client_config.GetExpireTimeInHours().Milliseconds()), "/", authClient.conf.O.Domain, false, true)
			}
		}

		if flag == auth_common.FLAG_LOOGIN {
			protocol := "http://"
			if ctx.Request.TLS != nil {
				protocol = "https://"
			}
			ctx.Set("flag", flag)
			returnURl := protocol + ctx.Request.Host + filterURl(ctx.Request.URL.RequestURI())
			ctx.Redirect(http.StatusMovedPermanently, returnURl)
			return
		}
	}

}

func (authClient *ClientContext) NxAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authClient.handlerQuerySign(ctx)
		isAllowPath := authClient.ignorePath(ctx)

		authorization := ctx.GetHeader(auth_common.TokenHeaderName)
		cookiesToken, cookiesErr := ctx.Cookie(auth_common.TokenCookiesName)

		prefix := auth_common.TokenHeaderNamePrefix
		prefixLen := len(prefix) //Bearer可自定义为一个常量

		var authErr error = nil
		authToken := ""
		if cookiesToken != "" && cookiesErr != nil {
			if !strings.HasPrefix(cookiesToken, prefix) {
				authErr = errors.New("cookie's token is invalid！")
			} else {
				authToken = cookiesToken[prefixLen:]
			}
		} else if authorization != "" {
			if !strings.HasPrefix(authorization, prefix) {
				authErr = errors.New("header's token is invalid！")
			} else {
				authToken = authorization[prefixLen:]
			}
		} else {
			authErr = errors.New("token is null！")
		}

		//验证token字符串
		if authToken != "" {
			claim, newToken, err := authClient.conf.Jwt.ValidAndRefreshToken(authToken)
			if err != nil || claim == nil {
				authErr = errors.New("token is invalid！")
			} else if time.Now().Unix() > claim.ExpiresAt {
				authErr = errors.New("token is expired！")
			} else {
				if authClient.conf.Cache != nil {
					authUser, found := authClient.conf.Cache.Get(auth_common.GetAuthCacheKeyString(claim.Id))
					if !found || authUser == nil {
						authErr = errors.New("token is not exist or invalid！")
					} else {
						ctx.Set(auth_common.ContextTokenName, authToken)
						ctx.Set(auth_common.ContextAuthUser, authUser)
					}
				} else {
					authUser, remoteErr := authClient.apiImpl.AuthApi.GetCurrentUserInfo(ctx, authToken)
					if remoteErr != nil || authUser == nil || authUser.ID < 0 {
						authErr = errors.New("token is not exist or invalid！")
					} else {
						ctx.Set(auth_common.ContextTokenName, authToken)
						ctx.Set(auth_common.ContextAuthUser, &authUser)
					}
				}

				//刷新新的token
				if newToken != "" {
					ctx.Set(auth_common.ContextTokenName, newToken)
					ctx.Header(auth_common.NewTokenName, newToken)
					ctx.SetCookie(auth_common.TokenCookiesName, newToken, int(client_config.GetExpireTimeInHours().Milliseconds()), "/", authClient.conf.O.Domain, false, true)
				}
			}
		}
		if authErr == nil || isAllowPath {
			ctx.Next()
		} else {
			authClient.conf.Logger.Error(authErr.Error())
			authClient.goLoginPage(ctx)
			return //ctx.Abort()
		}
	}
}
