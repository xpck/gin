// response.go
// wrap response.

package gin

import (
	"errors"
	"net/http"
)

const (
	empty        = "empty"
	_msg         = "msg"
	_ok          = "ok"
	_data        = "data"
	_code        = "err_code"
	_err         = "err"
	_serverError = "internal server error"
	_badRequest  = "client bad request"
	_total       = "total"
	_list        = "list"
)

// SuccessOK return {"msg": "ok"}, the status code is http.StatusOK.
func (c *Context) SuccessOK() {
	c.JSON(http.StatusOK, H{
		_msg: _ok,
	})
}

// CreatedOK return {"msg": "ok", "data": data}, the status code is http.StatusCreated
func (c *Context) CreatedOK(data any) {
	c.JSON(http.StatusCreated, H{
		_msg:  _ok,
		_data: data,
	})
}

// SuccessData return {"msg": "ok", "data": data}, the status code is http.StatusOK
func (c *Context) SuccessData(data any) {
	c.JSON(http.StatusOK, H{
		_msg:  _ok,
		_data: data,
	})
}

// NoData return {"msg": "empty", "data": data}, the status code is http.StatusNoContent
func (c *Context) NoData(data any) {
	c.JSON(http.StatusNoContent, H{
		_msg:  empty,
		_data: data,
	})
}

// ServerError return {"err_code": errCode, "err": err.Error()}, the status code is http.StatusInternalServerError
// If err is nil, the err = "internal server error", but not suggest.
// The "code" is error code, you must custom define it.
func (c *Context) ServerError(errCode int, err error) {
	if err == nil {
		err = errors.New(_serverError)
	}
	c.JSON(http.StatusInternalServerError, H{
		_code: errCode,
		_err:  err.Error(),
	})

	c.Errors = append(c.Errors, &Error{
		Err:  err,
		Meta: errCode,
	})
}

// BadRequest return {"err": err.Error()}, the status code is http.StatusBadRequest
// If err is nil, the err = "client bad request", but not suggest
func (c *Context) BadRequest(err error) {
	if err == nil {
		err = errors.New(_badRequest)
	}
	c.JSON(http.StatusBadRequest, H{
		_err: err.Error(),
	})

	c.Errors = append(c.Errors, &Error{
		Err: err,
	})
}

// BadReqStr return {"err": msg}, the status code is http.StatusBadRequest
// If err is nil, the err = "client bad request", but not suggest
func (c *Context) BadReqStr(msg string) {
	if msg == "" {
		msg = _badRequest
	}
	c.JSON(http.StatusBadRequest, H{
		_err: msg,
	})

	c.Errors = append(c.Errors, &Error{
		Err: errors.New(msg),
	})
}

// List return {"msg": "ok", "data": {"list": list, "total": count}}, the status code is http.StatusOK.
func (c *Context) List(count int, list any) {
	c.JSON(http.StatusOK, H{
		_msg: _ok,
		_data: H{
			_list:  list,
			_total: count,
		},
	})
}

// Unauthorized return {"data": data} with status http.StatusUnauthorized.
func (c *Context) Unauthorized(data any) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, H{
		_data: data,
	})
}

// Forbidden return {"err_code": errCode} with status http.StatusForbidden
func (c *Context) Forbidden(errCode int) {
	c.AbortWithStatusJSON(http.StatusForbidden, H{
		_code: errCode,
	})

	c.Errors = append(c.Errors, &Error{
		Meta: errCode,
	})
}
