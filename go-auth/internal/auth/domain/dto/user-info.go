package dto

// 返回给前端的当前用户信息
type UserDTO struct {
	ID           int64   `json:"id"`
	Username     string  `json:"username"`
	UserType     int32   `json:"userType"`
	Mobile       string  `json:"mobile"`
	Email        string  `json:"email"`
	Avatar       string  `json:"avatar"`
	Nickname     *string `json:"nickname"`
	Introduction *string `json:"introduction"`
	WxOpenId     string  `json:"wxOpenId"`
}

type UserIdDTO struct {
	ID          uint   `json:"id"`
	Mobile      string `json:"mobile"`
	UserName    string `json:"username"`
	MultiRecord bool   `json:"multiRecord"` //存在多条记录
}
