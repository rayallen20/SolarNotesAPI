package planet

type ListPlanetResp struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`
	Intro string `json:"intro"`
}

type ListResp struct {
	Planets []ListPlanetResp `json:"planets"`
}
