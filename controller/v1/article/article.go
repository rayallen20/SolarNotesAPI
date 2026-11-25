package article

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"SolarNotesAPI/biz"
	reqArticle "SolarNotesAPI/param/req/v1/article"
	respArticle "SolarNotesAPI/param/resp/v1/article"
	"SolarNotesAPI/resp_body"
	"SolarNotesAPI/sys_error"

	"github.com/gin-gonic/gin"
)

func Show(readMarkdown func(string) (string, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req reqArticle.ShowReq
		var body resp_body.RespBody
		if err := c.ShouldBindJSON(&req); err != nil {
			body.ValidateFailed(err)
			c.JSON(http.StatusOK, body)
			return
		}
		a, err := biz.ReadArticle(c.Request.Context(), req.Catalogue.ID, readMarkdown)
		if err != nil {
			var dbErr *sys_error.DatabaseError
			switch {
			case errors.Is(err, sys_error.ArticleUnavailableErr):
				body.ArticleUnavailableError()
			case errors.As(err, &dbErr):
				body.DBError()
			default:
				body.InternalError()
			}
			c.JSON(http.StatusOK, body)
			return
		}
		body.Success(map[string]any{"article": respArticle.ShowResp{
			ID:        strconv.FormatInt(a.ID, 10),
			Content:   a.Content,
			CreatedAt: a.CreatedAt.Format(time.DateOnly),
		}})
		c.JSON(http.StatusOK, body)
	}
}
