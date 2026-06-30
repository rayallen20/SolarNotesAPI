package article

type ShowReq struct {
	Catalogue *ShowCatalogue `json:"catalogue" binding:"required"`
}

type ShowCatalogue struct {
	Id *int `json:"id" binding:"required,min=1"`
}
