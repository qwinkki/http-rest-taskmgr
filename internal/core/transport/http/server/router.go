package core_http_server

import (
	"fmt"
	"net/http"
)

type ApiVersion string

const (
	APIVersionV1 = ApiVersion("v1")
	ApiVersionV2 = ApiVersion("v2")
	ApiVersionV3 = ApiVersion("v3")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion ApiVersion
}

func NewAPIVersionRouter(apiVersion ApiVersion) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
	}
}

func (r *APIVersionRouter) RegisterRoutes(route ...Route) {
	for _, route := range route {
		pattern := fmt.Sprintf("/%s%s", r.apiVersion, route.Path)
		r.Handle(pattern, route.Handler)
	}
}
