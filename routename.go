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
