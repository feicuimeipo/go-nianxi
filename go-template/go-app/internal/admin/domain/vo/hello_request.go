package vo

// 操作日志请求结构体
type HelloListRequest struct {
	Msg      string `json:"msg" form:"msg"`
	PageNum  int    `json:"pageNum" form:"pageNum"`
	PageSize int    `json:"pageSize" form:"pageSize"`
}
