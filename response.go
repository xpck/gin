package gin

import (
	"errors"
	"net/http"
)

const (
	_ok    = "ok"
	_data  = "data"
	_err   = "err"
	_msg   = "msg"
	_total = "total"

	_errParameter = "invalid parameter"
	_noContent    = "no content"
)

func (c *Context) ServerError(msg string, err error) {
	if err == nil {
		err = errors.New(msg)
	}
	c.JSON(http.StatusOK, H{
		_msg: msg,
		_err: err.Error(),
	})

	c.Errors = append(c.Errors, &Error{
		Err:  err,
		Type: ErrorTypeAny,
		Meta: msg,
	})
}

// NoContent return empty struct or nil
func (c *Context) NoContent(data any) {
	c.JSON(http.StatusOK, H{
		_msg:  _noContent,
		_data: data,
	})
}

func (c *Context) BadRequest(err error) {
	if err == nil {
		err = errors.New(_errParameter)
	}
	c.JSON(http.StatusBadRequest, H{
		_msg: err.Error(),
		_err: _errParameter,
	})

	c.Errors = append(c.Errors, &Error{
		Err:  err,
		Type: ErrorTypeAny,
		Meta: _errParameter,
	})
}

func (c *Context) BadReqStr(msg string) {
	c.JSON(http.StatusBadRequest, H{
		_msg: msg,
		_err: _errParameter,
	})

	c.Errors = append(c.Errors, &Error{
		Err:  errors.New(msg),
		Type: ErrorTypeAny,
		Meta: _errParameter,
	})
}

func (c *Context) SuccessOK() {
	c.JSON(http.StatusOK, H{
		_msg: _ok,
	})
}

func (c *Context) SuccessData(data any) {
	c.JSON(http.StatusOK, H{
		_msg:  _ok,
		_data: data,
	})
}

func (c *Context) SuccessTotal(list any, total int) {
	c.JSON(http.StatusOK, H{
		_msg:   _ok,
		_data:  list,
		_total: total,
	})
}

func (c *Context) CreatedOK(data any) {
	c.JSON(http.StatusCreated, H{
		_msg:  _ok,
		_data: data,
	})
}
