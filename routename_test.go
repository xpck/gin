package gin

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin/internal/json"
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

	name, b := GetApiName(http.MethodGet, "/audit/log/metrics")
	assert.Equal(t, "query audit metrics", name)
	assert.True(t, b)
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
	initApi()
	var (
		name string
		exit bool
	)
	name, exit = GetApiName(http.MethodGet, "/hello")
	assert.Equal(t, "hello", name)
	assert.True(t, exit)

	name, exit = GetApiName(http.MethodDelete, "/base/menu/:id")
	assert.Equal(t, "delete menu by id", name)
	assert.True(t, exit)

	name, exit = GetApiName(http.MethodPut, "/base/config")
	assert.Equal(t, "put basic service config", name)
	assert.True(t, exit)

	name, exit = GetApiName(http.MethodGet, "/base/log/mysql/:id")
	assert.Equal(t, "query mysql log by id", name)
	assert.True(t, exit)

	// json api, pretty it, and it's expected!
	marshal, _ := json.Marshal(api)
	fmt.Println(string(marshal))
	//
	//group, _ := GetGroup("/audit")
	//marshal, _ = json.Marshal(group)
	//fmt.Println(string(marshal))

	logGroup, b := GetGroup("/base/log")
	assert.True(t, b)
	assert.Equal(t, "query mysql log by id", logGroup.getRouteName(http.MethodGet, "/base/log/mysql/:id"))
}

func TestGetApiMap(t *testing.T) {
	initApi()
	expect := map[[2]string]string{
		[2]string{"DELETE", "/base/menu/:id"}:   "delete menu by id",
		[2]string{"DELETE", "/base/user/:id"}:   "delete user by id",
		[2]string{"GET", "/audit/ping"}:         "ping audit",
		[2]string{"GET", "/base/log/mysql/:id"}: "query mysql log by id",
		[2]string{"GET", "/base/menu/:id"}:      "query menu by id",
		[2]string{"GET", "/base/user/:id"}:      "query user by id",
		[2]string{"GET", "/hello"}:              "hello",
		[2]string{"PUT", "/base/config"}:        "put basic service config"}
	assert.Equal(t, expect, GetApiMap())
}

func initApi() {
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
		audit.POSTEX("", "create audit", emptyHandler)
	}
	{
		engine.GETEX("hello", "hello", emptyHandler)
	}
}

func Test_setGroup(t *testing.T) {
	var apiGroup = &ApiGroup{
		Path:  "/",
		Name:  "root",
		Group: make([]*ApiGroup, 0),
		Api:   []ApiInfo{},
	}
	apiGroup.setGroup("/", "/base", "basic service")
	apiGroup.setGroup("/", "/base", "base service")
	assert.Equal(t, len(apiGroup.Group), 1)
	apiGroup.setGroup("/", "/audit", "base service")
	assert.Equal(t, 2, len(apiGroup.Group))
}

func TestGetApiList(t *testing.T) {
	initApi()
	list := GetApiList(nil)
	for _, info := range list {
		t.Log(info)
	}
}
