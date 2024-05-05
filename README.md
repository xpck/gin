# Gin Web Framework
> Fork from `github.com/gin-gonic/gin`

## Feature
- Record api name.
- Record group name, and group's api.
- Add two struct to record api with using `RouterGroup.GroupEX()` and `RouterGroup.PUTEX(),RouterGroup.GETEX(),RouterGroup.DELETEEX(),RouterGroup.GETEX(),RouterGroup.PATCHEX()`
```go
package gin

type ApiInfo struct {
	Name     string `json:"name"`
	FullPath string `json:"full_path"`
	Method   string `json:"method"`
}

type ApiGroup struct {
	Path  string      `json:"path"`
	Name  string      `json:"name"`
	Group []*ApiGroup `json:"group"`
	Api   []ApiInfo   `json:"api"`
}
```
- Add two function to get information about api.(If you use builtin method, that you can't get information from the following two APIs)
  - `func GetGroup(path string) (*ApiGroup, bool)`
  - `func GetApiName(method string, fullPath string) (string, bool)`
```go
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.New()
	emptyHandler := func(*gin.Context) {}
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
	name, exit = gin.GetApiName(http.MethodPatch, "/hello")
	fmt.Println(name,exit) // "", false
	
	//assert.Equal(t, "", name)
	//assert.Equal(t, false, exit)
	
	name, exit = gin.GetApiName(http.MethodGet, "/hello")
	//assert.Equal(t, "hello", name)
	//assert.Equal(t, true, exit)
	

	name, exit = gin.GetApiName(http.MethodDelete, "/base/menu/:id")
	//assert.Equal(t, "delete menu by id", name)
	//assert.Equal(t, true, exit)

	name, exit =gin.GetApiName(http.MethodPut, "/base/config")
	//assert.Equal(t, "put basic service config", name)
	//assert.Equal(t, true, exit)

	name, exit = gin.GetApiName(http.MethodGet, "/base/log/mysql/:id")
	//assert.Equal(t, "query mysql log by id", name)
	//assert.Equal(t, true, exit)
	
	// marshal api, pretty it, and it is expected!
	//marshal, _ := json.Marshal(api)
	//fmt.Println(string(marshal))
	//
	//group, _ := GetGroup("/audit")
	//marshal, _ = json.Marshal(group)
	//fmt.Println(string(marshal))
}
```
```json
{
  "path":"/",
  "name":"root",
  "group":[
    {
      "path":"/base",
      "name":"basic service",
      "group":[
        {
          "path":"/base/menu",
          "name":"menu service",
          "group":[ ],
          "api":[
            {
              "name":"query menu by id",
              "full_path":"/base/menu/:id",
              "method":"GET"
            },
            {
              "name":"delete menu by id",
              "full_path":"/base/menu/:id",
              "method":"DELETE"
            }
          ]
        },
        {
          "path":"/base/user",
          "name":"user service",
          "group":[ ],
          "api":[
            {
              "name":"query user by id",
              "full_path":"/base/user/:id",
              "method":"GET"
            },
            {
              "name":"delete user by id",
              "full_path":"/base/user/:id",
              "method":"DELETE"
            }
          ]
        },
        {
          "path":"/base/log",
          "name":"basic log service",
          "group":[
            {
              "path":"/base/log/mysql",
              "name":"basic log service by mysql",
              "group":[ ],
              "api":[
                {
                  "name":"query mysql log by id",
                  "full_path":"/base/log/mysql/:id",
                  "method":"GET"
                }
              ]
            }
          ],
          "api":[ ]
        }
      ],
      "api":[
        {
          "name":"put basic service config",
          "full_path":"/base/config",
          "method":"PUT"
        }
      ]
    },
    {
      "path":"/audit",
      "name":"audit service",
      "group":[ ],
      "api":[
        {
          "name":"ping audit",
          "full_path":"/audit/ping",
          "method":"GET"
        }
      ]
    }
  ],
  "api":[
    {
      "name":"hello",
      "full_path":"/hello",
      "method":"GET"
    }
  ]
}
```