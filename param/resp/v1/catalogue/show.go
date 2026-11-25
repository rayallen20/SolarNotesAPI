package catalogue

type Node struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
	*FolderFields
}

// FolderFields 仅嵌入 folder 的响应, 保留空简介和空子数组.
type FolderFields struct {
	Intro     string `json:"intro"`
	CreatedAt string `json:"createdAt"`
	Children  []Node `json:"children"`
}
