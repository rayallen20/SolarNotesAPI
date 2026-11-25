package resp_body

import "SolarNotesAPI/sys_error"

const PlanetNotFoundCode = 20102

func (r *RespBody) PlanetNotFoundError() {
	r.failure(PlanetNotFoundCode, sys_error.PlanetNotFoundErr.Error())
}
