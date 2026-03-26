package core_http_server

import "net/http"

type Route struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Handler http.HandlerFunc
}

/*func NewRoute(method, path string, handler http.HandlerFunc) Route {
	return Route{
		method,
		path,
		handler,
	}
}*/
