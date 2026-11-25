package catalogue

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"SolarNotesAPI/biz"
	reqCatalogue "SolarNotesAPI/param/req/v1/catalogue"
	respCatalogue "SolarNotesAPI/param/resp/v1/catalogue"
	"SolarNotesAPI/resp_body"
	"SolarNotesAPI/sys_error"

	"github.com/gin-gonic/gin"
)

// Show 返回可见书籍的公开目录树.
func Show(c *gin.Context) {
	var req reqCatalogue.ShowReq
	var body resp_body.RespBody
	if err := c.ShouldBindJSON(&req); err != nil {
		body.ValidateFailed(err)
		c.JSON(http.StatusOK, body)
		return
	}

	root, err := biz.ReadCatalogue(c.Request.Context(), req.Book.ID)
	if err != nil {
		var dbErr *sys_error.DatabaseError
		switch {
		case errors.Is(err, sys_error.BookNotFoundErr):
			body.BookNotFoundError()
		case errors.Is(err, sys_error.CatalogueNotFoundErr):
			body.CatalogueNotFoundError()
		case errors.As(err, &dbErr):
			body.DBError()
		default:
			body.InternalError()
		}
		c.JSON(http.StatusOK, body)
		return
	}

	body.Success(map[string]any{"catalogue": toResponse(root)})
	c.JSON(http.StatusOK, body)
}

func toResponse(node *biz.CatalogueNode) respCatalogue.Node {
	item := respCatalogue.Node{
		ID:   strconv.FormatInt(node.ID, 10),
		Type: node.Type,
		Name: node.Name,
	}
	if node.Type == "folder" {
		children := make([]respCatalogue.Node, 0, len(node.Children))
		for _, child := range node.Children {
			children = append(children, toResponse(child))
		}
		item.FolderFields = &respCatalogue.FolderFields{
			Intro:     node.Intro,
			CreatedAt: node.CreatedAt.Format(time.DateOnly),
			Children:  children,
		}
	}
	return item
}
