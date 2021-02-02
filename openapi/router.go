package openapi

import (
	"path"
	"reflect"
)

// OpFn customize operations
type OpFn func(o *Operation)

// Router is a extract of api path.
// router grouped in multiple levels can share parameters(path/query/header/etc)
type Router interface {
	Root() *OpenAPI
	WithParam(param *Param) Router
	WithPathParam(name, description string) Router
	WithTags(tags ...string) Router
	Route(path string, fn func(r Router)) Router
	// HTTP methods
	Method(method, path string, opFn OpFn)
	GET(path string, opFn OpFn)
	PUT(path string, opFn OpFn)
	POST(path string, opFn OpFn)
	DELETE(path string, opFn OpFn)
	HEAD(path string, opFn OpFn)
	PATCH(path string, opFn OpFn)
}

type router struct {
	root      *OpenAPI
	parent    *router
	path      string
	tags      []string
	params    []*Param
	paths     map[string]*Path
	subRoutes map[string]*router
}

// NewRouter create a new router
func NewRouter(root *OpenAPI) Router {
	return newRouter(root)
}

func newRouter(root *OpenAPI) *router {
	if root == nil {
		panic(ErrNoRoot)
	}
	return &router{
		root:      root,
		paths:     make(map[string]*Path),
		subRoutes: make(map[string]*router),
	}
}

// Root return document root
func (r *router) Root() *OpenAPI {
	return r.root
}

// GET short cut
func (r *router) GET(path string, opFn OpFn) {
	r.Method("get", path, opFn)
}

// PUT put
func (r *router) PUT(path string, opFn OpFn) {
	r.Method("put", path, opFn)
}

// POST post
func (r *router) POST(path string, opFn OpFn) {
	r.Method("post", path, opFn)
}

func (r *router) DELETE(path string, opFn OpFn) {
	r.Method("delete", path, opFn)
}

func (r *router) PATCH(path string, opFn OpFn) {
	r.Method("patch", path, opFn)
}

func (r *router) HEAD(path string, opFn OpFn) {
	r.Method("head", path, opFn)
}

// Method add method to router
func (r *router) Method(method, path string, opFn OpFn) {
	apiPath, exists := r.paths[path]

	// Retrieve upstream to collect things back
	retriveUpstream := func(fn func(upstream *router)) {
		upstream := r
		for upstream != nil {
			fn(upstream)
			upstream = upstream.parent
		}
	}
	if !exists {
		newPath := &Path{
			root:       r.root,
			operations: make(opMap),
		}

		pathParts := []string{path}
		retriveUpstream(func(upstream *router) {
			pathParts = append(pathParts, upstream.path)
			newPath.Parameters = append(newPath.Parameters, upstream.params...)
		})

		reverse(pathParts)
		fullPath := joinPathParts(pathParts...)
		apiPath, exists = r.root.Paths[fullPath]
		if !exists {
			r.paths[path] = newPath
			r.root.Paths[fullPath] = newPath
			apiPath = newPath
		}
	}
	tags := make([]string, 0)
	retriveUpstream(func(upstream *router) {
		tags = append(tags, upstream.tags...)
	})
	op := apiPath.AddOperation(method)
	opFn(op)
	op.Tags(tags...)
}

// WithParam add param
func (r *router) WithParam(param *Param) Router {
	r.params = append(r.params, param)
	return r
}

// WithPathParam add path param
func (r *router) WithPathParam(name, description string) Router {
	return r.WithParam(&Param{
		Name:        name,
		In:          PathParam,
		Description: description,
		Required:    true,
		Schema: &Schema{
			Type: "string",
		},
	})
}

// WithTags add tag to the path
func (r *router) WithTags(tags ...string) Router {
	r.tags = append(r.tags, tags...)
	return r
}

// Route to sub paths. Remember that the returned router is newly created **sub** router
func (r *router) Route(path string, fn func(r Router)) Router {
	sub := newRouter(r.root)
	sub.parent = r
	sub.path = path
	r.subRoutes[path] = sub
	if fn != nil {
		fn(sub)
	}

	return r
}

func joinPathParts(parts ...string) string {
	return path.Join(parts...)
}

func reverse(slice interface{}) {
	rv := reflect.ValueOf(slice)
	swap := reflect.Swapper(slice)
	length := rv.Len()
	half := length / 2
	for i := 0; i < half; i++ {
		j := length - i - 1
		swap(i, j)
	}
}
