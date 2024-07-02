package auth_client

import (
	"encoding/base64"
	"errors"
	"gitee.com/go-nianxi/go-auth/pkg/auth-client/auth-common"
	auth_jwt "gitee.com/go-nianxi/go-auth/pkg/auth-client/client-config"
	"github.com/gin-gonic/gin"
	"time"
)

type RedirectUrl struct {
	RedirectUrl string                `json:"redirectUrl"`
	Sign        string                `json:"sign"`
	EncodeToken string                `json:"encodeToken"`
	Flag        string                `json:"flag"`
	User        *auth_common.AuthUser `json:"user"`
	Token       string                `json:"token"`
}

func (authClient *ClientContext) LoginSuccess(authUser auth_common.AuthUser, c *gin.Context) (*RedirectUrl, error) {
	time := time.Hour * 7 * 24
	//if authClient.conf.O.Jwt.ExpireTimeInHours == 0 {
	//	authClient.conf.O.Jwt.ExpireTimeInHours
	//}
	//time := time.Duration(authClient.conf.O.Jwt.ExpireTimeInHours) * time.Hour
	token, err := authClient.conf.Jwt.CreateJwt(authUser.ID, authUser.Username)
	if err != nil {
		return nil, err
	}

	bearToken := auth_common.TokenHeaderNamePrefix + token
	authClient.conf.Cache.Set(auth_common.GetAuthCacheKey(authUser.ID), authUser, time)
	c.SetCookie(auth_common.TokenCookiesName, bearToken, int(time.Seconds()), "/", authClient.conf.O.Domain, false, true)
	c.Header(auth_common.TokenHeaderName, bearToken)

	redirectUrl := c.Query("redirectUrl")
	if redirectUrl == "" {
		redirectUrl = c.Param("redirectUrl")
	}

	publicK := auth_jwt.GenPublic(authClient.conf.O.Jwt.Secret)
	encodeToken := base64.URLEncoding.EncodeToString([]byte(token))
	sign := base64.URLEncoding.EncodeToString([]byte(publicK))
	flag := base64.URLEncoding.EncodeToString([]byte(auth_common.FLAG_LOOGIN))
	return &RedirectUrl{
		RedirectUrl: redirectUrl,
		Sign:        sign,
		EncodeToken: encodeToken,
		Flag:        flag,
		Token:       token,
		User:        &authUser,
	}, nil
}

func (authClient *ClientContext) Logout(token string) {
	//验证token字符串
	claim, _ := authClient.conf.Jwt.ParserToken(token)
	if claim != nil {
		authClient.conf.Cache.Delete(auth_common.GetAuthCacheKeyString(claim.Id))
	}
}

func (authClient *ClientContext) GetCurrentUser(ctx *gin.Context) (*auth_common.AuthUser, error) {
	ctxUser, exist := ctx.Get(auth_common.ContextAuthUser)
	if exist {
		u := auth_common.MapToAuthUser(ctxUser.(map[string]interface{})) // ctxUser.(auth_common.AuthUser)
		return u, nil
	}

	authorization := ctx.GetHeader(auth_common.TokenHeaderName)
	cookiesToken, _ := ctx.Cookie(auth_common.TokenCookiesName)

	prefixLen := len(auth_common.TokenHeaderNamePrefix) //Bearer可自定义为一个常量

	authToken := ""
	if cookiesToken != "" {
		authToken = cookiesToken[prefixLen:]
	} else if authorization != "" {
		authToken = authorization[prefixLen:]
	}

	if authToken != "" {
		claim, _, err := authClient.conf.Jwt.ValidAndRefreshToken(authToken)
		if err == nil && claim != nil && (time.Now().Unix() <= claim.ExpiresAt) {
			//保存上下文数据
			if authClient.conf.Cache != nil {
				if auth, found := authClient.conf.Cache.Get(auth_common.GetAuthCacheKeyString(claim.Id)); found {
					authUser := auth.(auth_common.AuthUser)
					return &authUser, nil
				}
			} else {
				if auth, err1 := authClient.apiImpl.AuthApi.GetCurrentUserInfo(ctx, authToken); err1 == nil {
					return auth, nil
				}
			}
		}
	}

	return nil, errors.New("用户未登录")
}
