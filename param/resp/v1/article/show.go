package article

// ShowResp contains the public fields of one article.
type ShowResp struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}
