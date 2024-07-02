package vo

type HelloRequest struct {
	Id uint `form:"id"       json:"id"         example:"123"        binding:"required,min=2,max=20"`
}
