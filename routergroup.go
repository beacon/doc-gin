// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"net/http"

	"github.com/beacon/doc-gin/openapi"
	"github.com/gin-gonic/gin"
)

// IRouter defines all router handle interface includes single and group router.
type IRouter interface {
	IRoutes
	Group(string, DocRouteFunc, ...HandlerFunc) *RouterGroup
}

// Context wrapper for gin.Context, to avoid user importing two "gin" packages
type Context struct {
	*gin.Context
}

// HandlerFunc def
type HandlerFunc func(*Context)

// DocRouteFunc to deal with doc
type DocRouteFunc func(openapi.Router)

// DocOpFunc to deal with operations
type DocOpFunc func(*openapi.Operation)

// IRoutes defines all router handle interface.
type IRoutes interface {
	Handle(string, string, DocOpFunc, ...HandlerFunc) IRoutes
	Any(string, DocOpFunc, ...HandlerFunc) IRoutes
	GET(string, DocOpFunc, ...HandlerFunc) IRoutes
	POST(string, DocOpFunc, ...HandlerFunc) IRoutes
	DELETE(string, DocOpFunc, ...HandlerFunc) IRoutes
	PATCH(string, DocOpFunc, ...HandlerFunc) IRoutes
	PUT(string, DocOpFunc, ...HandlerFunc) IRoutes
	OPTIONS(string, DocOpFunc, ...HandlerFunc) IRoutes
	HEAD(string, DocOpFunc, ...HandlerFunc) IRoutes
}

// RouterGroup is used internally to configure router, a RouterGroup is associated with
// a prefix and an array of handlers (middleware).
type RouterGroup struct {
	engine    *Engine
	docRouter openapi.Router
	*gin.RouterGroup
	root bool
}

var _ IRouter = &RouterGroup{}

// Group creates a new router group. You should add all the routes that have common middlewares or the same path prefix.
// For example, all the routes that use a common middleware for authorization could be grouped.
func (g *RouterGroup) Group(relativePath string, routeFn DocRouteFunc, handlers ...HandlerFunc) *RouterGroup {
	subGroup := &RouterGroup{
		engine:      g.engine,
		RouterGroup: g.RouterGroup.Group(relativePath, toGinHandlers(handlers)...),
	}
	if routeFn != nil && g.docRouter != nil {
		subGroup.docRouter = g.docRouter.Route(relativePath, routeFn)
	}
	return subGroup
}

func toGinHandlers(handlers []HandlerFunc) []gin.HandlerFunc {
	gh := make([]gin.HandlerFunc, len(handlers))
	for i, h := range handlers {
		gh[i] = func(c *gin.Context) {
			h(&Context{c})
		}
	}
	return gh
}

// Handle of any methods
func (g *RouterGroup) Handle(method, relativePath string, opFn DocOpFunc, handlers ...HandlerFunc) IRoutes {
	if opFn != nil && g.docRouter != nil {
		g.docRouter.Method(method, relativePath, openapi.OpFn(opFn))
	}
	g.RouterGroup.Handle(method, relativePath, toGinHandlers(handlers)...)
	return g
}

// POST is a shortcut for router.Handle("POST", path, handle).
func (g *RouterGroup) POST(relativePath string, docFn DocOpFunc, handlers ...HandlerFunc) IRoutes {
	return g.Handle(http.MethodPost, relativePath, docFn, handlers...)
}

// GET is a shortcut for router.Handle("GET", path, handle).
func (g *RouterGroup) GET(relativePath string, docFn DocOpFunc, handlers ...HandlerFunc) IRoutes {
	return g.Handle(http.MethodGet, relativePath, docFn, handlers...)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle).
func (g *RouterGroup) DELETE(relativePath string, docFn DocOpFunc, handlers ...HandlerFunc) IRoutes {
	return g.Handle(http.MethodDelete, relativePath, docFn, handlers...)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle).
func (g *RouterGroup) PATCH(relativePath string, docFn DocOpFunc, handlers ...HandlerFunc) IRoutes {
	return g.Handle(http.MethodPatch, relativePath, docFn, handlers...)
}

// PUT is a shortcut for router.Handle("PUT", path, handle).
func (g *RouterGroup) PUT(relativePath string, docFn DocOpFunc, handlers ...HandlerFunc) IRoutes {
	return g.Handle(http.MethodPut, relativePath, docFn, handlers...)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle).
func (g *RouterGroup) OPTIONS(relativePath string, docFn DocOpFunc, handlers ...HandlerFunc) IRoutes {
	return g.Handle(http.MethodOptions, relativePath, docFn, handlers...)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle).
func (g *RouterGroup) HEAD(relativePath string, docFn DocOpFunc, handlers ...HandlerFunc) IRoutes {
	return g.Handle(http.MethodHead, relativePath, docFn, handlers...)
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (g *RouterGroup) Any(relativePath string, docFn DocOpFunc, handlers ...HandlerFunc) IRoutes {
	g.Handle(http.MethodGet, relativePath, docFn, handlers...)
	g.Handle(http.MethodPost, relativePath, docFn, handlers...)
	g.Handle(http.MethodPut, relativePath, docFn, handlers...)
	g.Handle(http.MethodPatch, relativePath, docFn, handlers...)
	g.Handle(http.MethodHead, relativePath, docFn, handlers...)
	g.Handle(http.MethodOptions, relativePath, docFn, handlers...)
	g.Handle(http.MethodDelete, relativePath, docFn, handlers...)
	g.Handle(http.MethodConnect, relativePath, docFn, handlers...)
	g.Handle(http.MethodTrace, relativePath, docFn, handlers...)
	return g.returnObj()
}

func (g *RouterGroup) returnObj() IRoutes {
	if g.root {
		return g.engine
	}
	return g
}
