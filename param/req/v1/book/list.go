package book

type ListReq struct {
	Planet *ListPlanet `json:"planet" binding:"required"`
}

type ListPlanet struct {
	Id *int `json:"id" binding:"required,min=1"`
}
