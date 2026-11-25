package article

type ShowReq struct {
	Catalogue *ShowCatalogue `json:"catalogue" binding:"required"`
}

type ShowCatalogue struct {
	ID int64 `json:"id,string" binding:"gt=0"`
}
