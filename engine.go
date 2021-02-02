package gin

import (
	"net/http"

	"github.com/beacon/doc-gin/openapi"
	"github.com/gin-gonic/gin"
)

// Engine with doc/gin engine
type Engine struct {
	*gin.Engine
	openAPI  *openapi.OpenAPI
	docRoute openapi.Router
}

// NewEngine new engine
func NewEngine(enableOpenAPI bool) *Engine {
	g := &Engine{
		Engine: gin.New(),
	}
	if enableOpenAPI {
		openAPI, _ := openapi.New("3.0.0", openapi.Info{})
		g.openAPI = openAPI
		g.docRoute = openapi.NewRouter(g.openAPI)
	}
	return g
}

// Doc do something to OpenAPI document
func (g *Engine) Doc(fn func(*openapi.OpenAPI)) *openapi.OpenAPI {
	if g.openAPI != nil && fn != nil {
		fn(g.openAPI)
	}

	return g.openAPI
}

var _ IRouter = (*Engine)(nil)

// Group creates a new router group. You should add all the routes that have common middlewares or the same path prefix.
// For example, all the routes that use a common middleware for authorization could be grouped.
func (g *Engine) Group(relativePath string, routeFn DocRouteFunc, handlers ...HandlerFunc) *RouterGroup {
	subGroup := &RouterGroup{
		engine:      g,
		RouterGroup: g.Engine.Group(relativePath, toGinHandlers(handlers)...),
	}
	if routeFn != nil && g.docRoute != nil {
		subGroup.docRouter = g.docRoute.Route(relativePath, routeFn)
	}

	return &RouterGroup{
		engine: g,
	}
}

// Handle some method
func (g *Engine) Handle(method, relativePath string, docFn DocOpFunc, handlers ...HandlerFunc) IRoutes {
	g.Engine.Handle(method, relativePath, toGinHandlers(handlers)...)
	return g
}

// POST is a shortcut for router.Handle("POST", path, handle).
func (g *Engine) POST(relativePath string, docFn DocOpFunc, handlers ...HandlerFunc) IRoutes {
	return g.Handle(http.MethodPost, relativePath, docFn, handlers...)
}

// GET is a shortcut for router.Handle("GET", path, handle).
func (g *Engine) GET(relativePath string, docFn DocOpFunc, handlers ...HandlerFunc) IRoutes {
	return g.Handle(http.MethodGet, relativePath, docFn, handlers...)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle).
func (g *Engine) DELETE(relativePath string, docFn DocOpFunc, handlers ...HandlerFunc) IRoutes {
	return g.Handle(http.MethodDelete, relativePath, docFn, handlers...)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle).
func (g *Engine) PATCH(relativePath string, docFn DocOpFunc, handlers ...HandlerFunc) IRoutes {
	return g.Handle(http.MethodPatch, relativePath, docFn, handlers...)
}

// PUT is a shortcut for router.Handle("PUT", path, handle).
func (g *Engine) PUT(relativePath string, docFn DocOpFunc, handlers ...HandlerFunc) IRoutes {
	return g.Handle(http.MethodPut, relativePath, docFn, handlers...)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle).
func (g *Engine) OPTIONS(relativePath string, docFn DocOpFunc, handlers ...HandlerFunc) IRoutes {
	return g.Handle(http.MethodOptions, relativePath, docFn, handlers...)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle).
func (g *Engine) HEAD(relativePath string, docFn DocOpFunc, handlers ...HandlerFunc) IRoutes {
	return g.Handle(http.MethodHead, relativePath, docFn, handlers...)
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (g *Engine) Any(relativePath string, docFn DocOpFunc, handlers ...HandlerFunc) IRoutes {
	g.Handle(http.MethodGet, relativePath, docFn, handlers...)
	g.Handle(http.MethodPost, relativePath, docFn, handlers...)
	g.Handle(http.MethodPut, relativePath, docFn, handlers...)
	g.Handle(http.MethodPatch, relativePath, docFn, handlers...)
	g.Handle(http.MethodHead, relativePath, docFn, handlers...)
	g.Handle(http.MethodOptions, relativePath, docFn, handlers...)
	g.Handle(http.MethodDelete, relativePath, docFn, handlers...)
	g.Handle(http.MethodConnect, relativePath, docFn, handlers...)
	g.Handle(http.MethodTrace, relativePath, docFn, handlers...)
	return g
}
