package vo

// 获取接口列表结构体
type ApplicationListRequest struct {
	AppName  string `json:"method" form:"method"`
	Alias    string `json:"alias" form:"alias"`
	Title    string `json:"title" form:"title"`
	BaseUrl  string `json:"baseUrl" form:"baseUrl"`
	TypeId   uint   `json:"typeId" form:"TypeId" bind:"required,min=1,max=2"`
	PageNum  uint   `json:"pageNum" form:"pageNum"`
	PageSize uint   `json:"pageSize" form:"pageSize"`
}

// 创建接口结构体
type CreateApplicationRequest struct {
	AppName      string  `json:"appName"      form:"appName" bind:"required,min=1,max=20"`
	Alias        string  `json:"alias"        form:"alias"   bind:"required,min=1,max=20"`
	Title        string  `json:"title"        form:"title"   bind:"required min=0,max=100"`
	BaseUrl      string  `json:"baseUrl"      form:"baseUrl" bind:"required min=0,max=120"`
	Icon         string  `json:"icon"        form:"icon"     bind:"min=0,max=255"`
	Introduction *string `json:"introduction" form:"introduction" bind:"min=0,max=255"`
}

// 更新接口结构体
type UpdateApplicationRequest struct {
	AppName      string  `json:"appName"      form:"appName" bind:"required,min=1,max=20"`
	Alias        string  `json:"alias"        form:"alias"   bind:"required,min=1,max=20"`
	Title        string  `json:"title"        form:"title"   bind:"required min=0,max=100"`
	BaseUrl      string  `json:"baseUrl"      form:"baseUrl" bind:"required min=0,max=120"`
	Introduction *string `json:"introduction" form:"introduction" bind:"min=0,max=255"`
	Icon         string  `json:"icon"        form:"icon"     bind:"min=0,max=255"`
}

// 批量删除接口结构体
type DeleteApplicationRequest struct {
	Ids []uint `json:"ids" form:"ids"`
}
