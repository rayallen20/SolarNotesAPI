package book

import (
	"errors"
	"net/http"
	"strconv"

	"SolarNotesAPI/biz"
	reqBook "SolarNotesAPI/param/req/v1/book"
	respBook "SolarNotesAPI/param/resp/v1/book"
	"SolarNotesAPI/resp_body"
	"SolarNotesAPI/sys_error"

	"github.com/gin-gonic/gin"
)

// List 返回所选天体的公开书单.
func List(c *gin.Context) {
	var req reqBook.ListReq
	var body resp_body.RespBody
	if err := c.ShouldBindJSON(&req); err != nil {
		body.ValidateFailed(err)
		c.JSON(http.StatusOK, body)
		return
	}

	var service biz.Book
	books, err := service.List(c.Request.Context(), req.Planet.ID)
	if err != nil {
		var dbErr *sys_error.DatabaseError
		switch {
		case errors.Is(err, sys_error.PlanetNotFoundErr):
			body.PlanetNotFoundError()
		case errors.As(err, &dbErr):
			body.DBError()
		default:
			body.InternalError()
		}
		c.JSON(http.StatusOK, body)
		return
	}

	items := make([]respBook.ListItem, 0, len(books))
	for _, book := range books {
		items = append(items, respBook.ListItem{
			ID:    strconv.FormatInt(book.ID, 10),
			Title: book.Title,
		})
	}
	body.Success(map[string]any{"bookList": items})
	c.JSON(http.StatusOK, body)
}
