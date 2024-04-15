package gin

import "maps"

// route_name.go

var routeTable = make(map[string]string)

func setRouteName(relativePath, routeName string) {
	routeTable[relativePath] = routeName
}

// GetRouteNameByPath Get route name by uri path
func GetRouteNameByPath(path string) (string, bool) {
	name, b := routeTable[path]
	return name, b
}

func GetRouteTable() map[string]string {
	return maps.Clone(routeTable)
}
