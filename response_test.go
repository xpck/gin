package gin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext_SuccessOK(t *testing.T) {
	r := New()
	r.GET("/", func(c *Context) {
		c.SuccessOK()
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, `{"msg":"ok"}`, resp.Body.String())
}

func TestContext_Forbidden(t *testing.T) {
	r := New()
	r.GET("/", func(c *Context) {
		c.Forbidden(1001)
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusForbidden, resp.Code)
	assert.Equal(t, `{"err_code":1001}`, resp.Body.String())
}
