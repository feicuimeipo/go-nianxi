package vo

// 创建接口结构体
type CreateMenuRequest struct {
	Name          string `json:"name" form:"name"           bind:"required,min=1,max=50"`
	Title         string `json:"title" form:"title"         bind:"required,min=1,max=50"`
	Icon          string `json:"icon" form:"icon"           bind:"min=0,max=50"`
	Path          string `json:"path" form:"path"           bind:"required,min=1,max=100"`
	Redirect      string `json:"redirect" form:"redirect"   bind:"min=0,max=100"`
	Component     string `json:"component" form:"component" bind:"required,min=1,max=100"`
	Sort          uint   `json:"sort" form:"sort"             bind:"gte=1,lte=999"`
	Status        uint   `json:"status" form:"status"         bind:"oneof=1 2"`
	Hidden        uint   `json:"hidden" form:"hidden"         bind:"oneof=1 2"`
	NoCache       uint   `json:"noCache" form:"noCache"       bind:"oneof=1 2"`
	AlwaysShow    uint   `json:"alwaysShow" form:"alwaysShow" bind:"oneof=1 2"`
	Breadcrumb    uint   `json:"breadcrumb" form:"breadcrumb" bind:"oneof=1 2"`
	ActiveMenu    string `json:"activeMenu" form:"activeMenu" bind:"min=0,max=100"`
	ParentId      uint   `json:"parentId" form:"parentId"`
	ApplicationId uint   `json:"applicationId" form:"applicationId" bind:"required,min=1,max=20"`
}

// 更新接口结构体
type UpdateMenuRequest struct {
	Name          string `json:"name" form:"name"    bind:"required,min=1,max=50"`
	Title         string `json:"title" form:"title" bind:"required,min=1,max=50"`
	Icon          string `json:"icon" form:"icon" bind:"min=0,max=50"`
	Path          string `json:"path" form:"path" bind:"required,min=1,max=100"`
	Redirect      string `json:"redirect" form:"redirect" bind:"min=0,max=100"`
	Component     string `json:"component" form:"component" bind:"min=0,max=100"`
	Sort          uint   `json:"sort" form:"sort" bind:"gte=1,lte=999"`
	Status        uint   `json:"status" form:"status" bind:"oneof=1 2"`
	Hidden        uint   `json:"hidden" form:"hidden" bind:"oneof=1 2"`
	NoCache       uint   `json:"noCache" form:"noCache" bind:"oneof=1 2"`
	AlwaysShow    uint   `json:"alwaysShow" form:"alwaysShow" bind:"oneof=1 2"`
	Breadcrumb    uint   `json:"breadcrumb" form:"breadcrumb" bind:"oneof=1 2"`
	ActiveMenu    string `json:"activeMenu" form:"activeMenu" bind:"min=0,max=100"`
	ParentId      uint   `json:"parentId" form:"parentId"`
	ApplicationId uint   `json:"applicationId" form:"applicationId" bind:"required,min=1,max=20"`
}

// 删除接口结构体
type DeleteMenuRequest struct {
	MenuIds []uint `json:"menuIds" form:"menuIds"`
}
