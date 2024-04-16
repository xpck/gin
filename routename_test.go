package gin

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRouteNameByPath(t *testing.T) {
	engine := New()
	var emptyHandler = func(c *Context) {}
	engine.GETEX("user", "get user", emptyHandler)
	name, exist := GetRouteName(http.MethodGet, "/user")
	assert.Equal(t, name, "get user")
	assert.Equal(t, exist, true)

	engine.Group("user").DELETEEX(":id", "delete user by id", emptyHandler)
	name, exist = GetRouteName(http.MethodDelete, "/user/:id")
	assert.Equal(t, name, "delete user by id")
	assert.Equal(t, exist, true)

	engine.PUTEX("{id}", "update info", emptyHandler)
	name, exist = GetRouteName(http.MethodPut, "/{id}")
	assert.Equal(t, name, "update info")
	assert.Equal(t, exist, true)

	engine.POSTEX("user", "create user", emptyHandler)
	name, exist = GetRouteName(http.MethodPost, "/user")
	assert.Equal(t, name, "create user")
	assert.Equal(t, exist, true)

	_, exist = GetRouteName(http.MethodPost, "user")
	assert.Equal(t, exist, false)
}
