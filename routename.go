package gin

var (
	api = &ApiGroup{
		Path:  "/",
		Name:  "root",
		Group: make([]*ApiGroup, 0),
		Api:   []ApiInfo{},
	}
	// storage all api, [2]string{method, fullPath}: name
	apiMap = make(map[[2]string]string)
)

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
	apiMap[[2]string{method, fullPath}] = name
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
func GetGroup(path string) (*ApiGroup, bool) {
	group := api.getGroup(path)

	if group == nil {
		return nil, false
	}

	return group, true
}

// GetApiName return api name.
func GetApiName(method, fullPath string) (string, bool) {
	name, exist := apiMap[[2]string{method, fullPath}]

	return name, exist
}

// GetApiList returns all APIs for the group. If group is nil, return all APIs(gin.api).
func GetApiList(group *ApiGroup) []ApiInfo {
	if group == nil {
		group = api
	}

	return group.getApiTable()
}

func GetApiMap() map[[2]string]string {
	return apiMap
}

func (a *ApiGroup) getApiTable() []ApiInfo {
	table := make([]ApiInfo, 0)
	table = append(table, a.Api...)
	for _, group := range a.Group {
		table = append(table, group.getApiTable()...)
	}

	return table
}
