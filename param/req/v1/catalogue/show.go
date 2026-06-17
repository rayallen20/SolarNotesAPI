package catalogue

type ShowReq struct {
	Book *ShowBook `json:"book" binding:"required"`
}

type ShowBook struct {
	Id *int `json:"id" binding:"required"`
}
