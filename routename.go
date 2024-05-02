package gin

import "maps"

// routeTable storage [[method, path]: name]
// this is not concurrent safe
// example:
var routeTable = make(map[[2]string]string)

func setRouteName(method, relativePath, routeName string) {
	routeTable[[2]string{method, relativePath}] = routeName
}

// GetRouteName Get route name by method and path.
func GetRouteName(method, path string) (string, bool) {
	name, b := routeTable[[2]string{method, path}]
	return name, b
}

// GetRouteTable clone route table.
func GetRouteTable() map[[2]string]string {
	return maps.Clone(routeTable)
}

// ApiInfo record api information: name, path, method, such as: ApiInfo{ Name:"query user list", FullPath:"/user/:id", Method:"GET"}.
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

func (a *ApiGroup) setGroup(basePath, path, name string) { //nolint:unused
	if a.Path == basePath {
		a.Group = append(a.Group, &ApiGroup{
			Path:  path,
			Name:  name,
			Group: []*ApiGroup{},
			Api:   []ApiInfo{},
		})
	} else {
		for _, group := range a.Group {
			group.setGroup(basePath, path, name)
		}
	}
}

func (a *ApiGroup) setRoute(basePath, fullPath, name, method string) { //nolint:unused
	if a.Path == basePath {
		a.Api = append(a.Api, ApiInfo{
			Name:     name,
			FullPath: fullPath,
			Method:   method,
		})
	} else {
		for _, group := range a.Group {
			group.setRoute(basePath, fullPath, name, method)
		}
	}
}
