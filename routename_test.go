package gin

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_setRouteName(t *testing.T) {
	api.setGroup("/", "/base", "basic service")
	api.setGroup("/base", "/base/menu", "menu")

	api.setRoute(http.MethodGet, "/base/menu", "/base/menu/:id", "query menu by id")
	api.setRoute(http.MethodDelete, "/base/menu", "/base/menu/:id", "delete menu by id")

	assert.Equal(t, "query menu by id", api.getRouteName(http.MethodGet, "/base/menu/:id"))
	assert.Equal(t, "delete menu by id", api.getRouteName(http.MethodDelete, "/base/menu/:id"))

	api.setGroup("/", "/audit", "audit service")
	api.setGroup("/audit", "/audit/log", "audit log")

	api.setRoute(http.MethodGet, "/audit/log", "/audit/log", "query audit log")
	assert.Equal(t, "query audit log", api.getRouteName(http.MethodGet, "/audit/log"))

	api.setRoute(http.MethodGet, "/audit/log", "/audit/log/metrics", "query audit metrics")
	assert.Equal(t, "query audit metrics", api.getRouteName(http.MethodGet, "/audit/log/metrics"))
}

func TestGetGroup(t *testing.T) {
	api.setGroup("/", "/base", "basic service")
	api.setGroup("/", "/audit", "audit service")
	api.setGroup("/audit", "/audit/log", "audit log")
	tests := []struct {
		name      string
		path      string
		groupName string
		want      bool
	}{
		{name: "root", path: "/", groupName: "root", want: true},
		{name: "api", path: "/api", groupName: "", want: false},
		{name: "basic service", path: "/base", groupName: "basic service", want: true},
		{name: "audit service", path: "/audit", groupName: "audit service", want: true},
		{name: "audit log", path: "/audit/log", groupName: "audit log", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group, exist := GetGroup(tt.path)
			assert.Equal(t, tt.groupName, group.Name)
			assert.Equal(t, tt.want, exist)
		})
	}
}

func TestGetApiName(t *testing.T) {
	engine := New()
	emptyHandler := func(*Context) {}
	{
		base := engine.GroupEX("base", "basic service")
		base.PUTEX("config", "put basic service config", emptyHandler)
		{
			menuGroup := base.GroupEX("menu", "menu service")
			menuGroup.GETEX(":id", "query menu by id", emptyHandler)
			menuGroup.DELETEEX(":id", "delete menu by id", emptyHandler)
		}

		{
			userGroup := base.GroupEX("user", "user service")
			userGroup.GETEX(":id", "query user by id", emptyHandler)
			userGroup.DELETEEX(":id", "delete user by id", emptyHandler)
		}

		{
			logGroup := base.GroupEX("log", "basic log service")
			mysqlGroup := logGroup.GroupEX("mysql", "basic log service by mysql")
			mysqlGroup.GETEX(":id", "query mysql log by id", emptyHandler)
		}
	}
	{
		audit := engine.GroupEX("audit", "audit service")
		audit.GETEX("ping", "ping audit", emptyHandler)
	}
	{
		engine.GETEX("hello", "hello", emptyHandler)
	}
	var (
		name string
		exit bool
	)
	name, exit = GetApiName(http.MethodGet, "/hello")
	assert.Equal(t, "hello", name)
	assert.Equal(t, true, exit)

	name, exit = GetApiName(http.MethodPatch, "/hello")
	assert.Equal(t, "", name)
	assert.Equal(t, false, exit)

	name, exit = GetApiName(http.MethodDelete, "/base/menu/:id")
	assert.Equal(t, "delete menu by id", name)
	assert.Equal(t, true, exit)

	name, exit = GetApiName(http.MethodPut, "/base/config")
	assert.Equal(t, "put basic service config", name)
	assert.Equal(t, true, exit)

	name, exit = GetApiName(http.MethodGet, "/base/log/mysql/:id")
	assert.Equal(t, "query mysql log by id", name)
	assert.Equal(t, true, exit)
	// marshal api, pretty it, and it is expected!
	//marshal, _ := json.Marshal(api)
	//fmt.Println(string(marshal))
	//
	//group, _ := GetGroup("/audit")
	//marshal, _ = json.Marshal(group)
	//fmt.Println(string(marshal))

	logGroup, b := GetGroup("/base/log")
	assert.Equal(t, true, b)
	assert.Equal(t, "query mysql log by id", logGroup.getRouteName(http.MethodGet, "/base/log/mysql/:id"))
}
