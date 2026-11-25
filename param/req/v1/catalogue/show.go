package catalogue

type ShowReq struct {
	Book *ShowBook `json:"book" binding:"required"`
}

type ShowBook struct {
	ID int64 `json:"id,string" binding:"gt=0"`
}
