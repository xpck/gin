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
	_code        = "code"
	_err         = "err"
	_serverError = "server error"
	_badRequest  = "bad request"
	_total       = "total"
)

func (c *Context) SuccessOK() {
	c.JSON(http.StatusOK, H{
		_msg: _ok,
	})
}

func (c *Context) CreatedOK(data any) {
	c.JSON(http.StatusCreated, H{
		_msg:  _ok,
		_data: data,
	})
}

func (c *Context) SuccessData(data any) {
	c.JSON(http.StatusOK, H{
		_data: data,
		_msg:  _ok,
	})
}

func (c *Context) NoData(data any) {
	c.JSON(http.StatusOK, H{
		_msg:  empty,
		_data: data,
	})
}

func (c *Context) ServerError(code int, err error) {
	if err == nil {
		err = errors.New(_serverError)
	}
	c.JSON(http.StatusInternalServerError, H{
		_code: code,
		_err:  err.Error(),
	})

	c.Errors = append(c.Errors, &Error{
		Err:  err,
		Meta: code,
	})
}

func (c *Context) BadRequest(code int, err error) {
	if err == nil {
		err = errors.New(_badRequest)
	}
	c.JSON(http.StatusBadRequest, H{
		_code: code,
		_err:  err.Error(),
	})

	c.Errors = append(c.Errors, &Error{
		Err:  err,
		Meta: code,
	})
}

func (c *Context) List(count int, list any) {
	c.JSON(http.StatusOK, H{
		_data:  list,
		_total: count,
		_msg:   _ok,
	})
}
