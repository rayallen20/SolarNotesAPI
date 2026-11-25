package sys_error

import "errors"

var PlanetNotFoundErr = errors.New("所选天体不存在,请刷新天体列表")
