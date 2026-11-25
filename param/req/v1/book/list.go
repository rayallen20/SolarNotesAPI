package book

type ListReq struct {
	Planet *ListPlanet `json:"planet" binding:"required"`
}

type ListPlanet struct {
	ID int64 `json:"id,string" binding:"gt=0"`
}
