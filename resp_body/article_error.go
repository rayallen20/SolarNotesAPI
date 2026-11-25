package resp_body

import "SolarNotesAPI/sys_error"

const ArticleUnavailableCode = 20401

func (r *RespBody) ArticleUnavailableError() {
	r.failure(ArticleUnavailableCode, sys_error.ArticleUnavailableErr.Error())
}
