package gin

var api = &ApiGroup{
	Path:  "/",
	Name:  "root",
	Group: make([]*ApiGroup, 0),
	Api:   []ApiInfo{},
}

// ApiInfo record api information, about name, path, method, such as: ApiInfo{ Name:"query user list", FullPath:"/user/:id", Method:"GET"}.
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

func (a *ApiGroup) setGroup(basePath, path, name string) {
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

func (a *ApiGroup) setRoute(method, basePath, fullPath, name string) {
	if a.Path == basePath {
		a.Api = append(a.Api, ApiInfo{
			Name:     name,
			FullPath: fullPath,
			Method:   method,
		})
	} else {
		for _, group := range a.Group {
			group.setRoute(method, basePath, fullPath, name)
		}
	}
}

func (a *ApiGroup) getRouteName(method string, fullPath string) string {
	for _, info := range a.Api {
		if info.FullPath == fullPath && info.Method == method {
			return info.Name
		}
	}
	for _, group := range a.Group {
		if name := group.getRouteName(method, fullPath); name != "" {
			return name
		}
	}

	return ""
}

func (a *ApiGroup) getGroup(path string) *ApiGroup {
	if a.Path == path {
		return a
	}
	for _, group := range a.Group {
		if result := group.getGroup(path); result != nil {
			return result
		}
	}

	return nil
}

// GetGroup return group's api info.
func GetGroup(path string) (ApiGroup, bool) {
	group := api.getGroup(path)

	if group == nil {
		return ApiGroup{}, false
	}

	return *group, true
}

// GetApiName return api name.
func GetApiName(method string, fullPath string) (string, bool) {
	name := api.getRouteName(method, fullPath)

	return name, name != ""
}
