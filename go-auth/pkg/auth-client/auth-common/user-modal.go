package auth_common

import (
	auth "gitee.com/go-nianxi/go-auth/pkg/auth-client/api/authentication"
)

// 返回给前端的当前用户信息
type AuthUser struct {
	ID           uint   `json:"ID"`
	Username     string `json:"username"`
	Mobile       string `json:"mobile"`
	Email        string `json:"email"`
	Avatar       string `json:"avatar"`
	Nickname     string `json:"nickname"`
	Introduction string `json:"introduction"`
	WxOpenId     string `json:"wxOpenId"`
}

func MapToAuthUser(m map[string]interface{}) *AuthUser {
	var authUser = new(AuthUser)
	for k, v := range m {
		switch k {
		case "ID":
			authUser.ID = uint(v.(float64))
		case "username":
			authUser.Username = v.(string)
		case "mobile":
			authUser.Mobile = v.(string)
		case "email":
			authUser.Email = v.(string)
		case "avatar":
			authUser.Avatar = v.(string)
		case "nickname":
			authUser.Nickname = v.(string)
		case "Introduction":
			authUser.Introduction = v.(string)
		case "WxOpenId":
			authUser.WxOpenId = v.(string)
		}
	}
	return authUser
}

func ToAuthUser(u *auth.UserInfoDTO) *AuthUser {
	dto := AuthUser{}
	dto.ID = uint(u.Id)
	dto.Username = u.Username
	dto.Mobile = u.Mobile
	dto.Email = u.Email
	dto.Avatar = u.Avatar
	dto.Nickname = u.Nickname
	dto.WxOpenId = u.WxOpenId
	dto.Introduction = u.Introduction

	return &dto
}

func ToUserInfo(u *AuthUser) *auth.UserInfoDTO {
	dto := auth.UserInfoDTO{}
	dto.Id = int64(u.ID)
	dto.Username = u.Username
	dto.Mobile = u.Mobile
	dto.Email = u.Email
	dto.Avatar = u.Avatar
	dto.Nickname = u.Nickname
	dto.WxOpenId = u.WxOpenId
	dto.Introduction = u.Introduction

	return &dto
}
