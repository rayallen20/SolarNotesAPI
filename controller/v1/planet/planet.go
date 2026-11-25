package planet

import (
	"errors"
	"net/http"
	"strconv"

	"SolarNotesAPI/biz"
	respPlanet "SolarNotesAPI/param/resp/v1/planet"
	"SolarNotesAPI/resp_body"
	"SolarNotesAPI/sys_error"

	"github.com/gin-gonic/gin"
)

// List 返回太阳系场景所需的天体简介。
func List(c *gin.Context) {
	var body resp_body.RespBody
	var service biz.Planet
	planets, err := service.List(c.Request.Context())
	if err != nil {
		var dbErr *sys_error.DatabaseError
		if errors.As(err, &dbErr) {
			body.DBError()
		} else {
			body.InternalError()
		}
		c.JSON(http.StatusOK, body)
		return
	}

	list := respPlanet.ListResp{
		Planets: make([]respPlanet.ListPlanetResp, 0, len(planets)),
	}
	for _, planet := range planets {
		list.Planets = append(list.Planets, respPlanet.ListPlanetResp{
			ID:    strconv.FormatInt(planet.ID, 10),
			Name:  planet.Name,
			Title: planet.Title,
			Intro: planet.Intro,
		})
	}
	body.Success(map[string]any{"planets": list.Planets})
	c.JSON(http.StatusOK, body)
}
