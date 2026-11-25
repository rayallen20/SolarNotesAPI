package resp_body

import "SolarNotesAPI/sys_error"

const (
	BookNotFoundCode      = 20202
	CatalogueNotFoundCode = 20301
)

func (r *RespBody) BookNotFoundError() {
	r.failure(BookNotFoundCode, sys_error.BookNotFoundErr.Error())
}

func (r *RespBody) CatalogueNotFoundError() {
	r.failure(CatalogueNotFoundCode, sys_error.CatalogueNotFoundErr.Error())
}
