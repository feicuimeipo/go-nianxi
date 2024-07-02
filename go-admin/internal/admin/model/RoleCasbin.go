package model

// 角色权限规则 -临时表(用于比对导入)s
type RoleCasbin struct {
	Keyword string `json:"v0"` // 角色关键字
	Path    string `json:"v1"` // 访问路径
	Method  string `json:"v2"` // 请求方式
	BaseUrl string `json:"v3"` // 请求方式
	v4      string `json:"v4"` // 请求方式
	v5      string `json:"v5"` // 请求方式
	ptype   string `json:"ptype"`
	id      uint64 `json:"id"`
}
