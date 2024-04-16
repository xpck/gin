# Gin Web Framework
> Fork from `github.com/gin-gonic/gin`

## Change
- Add some method to record route name.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.New()
	engine.GETEX("user", "query user", func(c *gin.Context) {})
	engine.POSTEX("user","create user", func(c *gin.Context) {})
	
	routeName, exist := gin.GetRouteName(http.MethodGet, "/user")
	fmt.Println(routeName == "query user")
    fmt.Println(exist == true)
}
```