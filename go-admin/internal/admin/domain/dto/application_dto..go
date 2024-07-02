package dto

type ApplicationTreeDto struct {
	ID       uint                  `json:"ID"`
	Desc     string                `json:"desc"`
	Children []*ApplicationTreeDto `json:"children"`
	Icon     string                `json:"icon"`
}
