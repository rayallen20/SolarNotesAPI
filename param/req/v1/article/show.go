package article

type ShowReq struct {
	Article *ShowArticle `json:"article" binding:"required"`
}

type ShowArticle struct {
	Id *int `json:"id" binding:"required,min=1"`
}
