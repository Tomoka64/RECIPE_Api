package router

import (
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"
)

type Group struct {
	prefix           string
	base             *httprouter.Router
	middleware       []MiddleWareFunc
	ctx              Context
	HTTPErrorHandler func(error, Context)
}

func (r *Router) Group(prefix string) *Group {
	return &Group{
		base:             r.base,
		prefix:           prefix,
		ctx:              r.ctx,
		middleware:       r.middleware,
		HTTPErrorHandler: r.HTTPErrorHandler,
	}
}

/*

g := r.Group("/")

g.Get("", func() {}) // We can access '/'
g.Get("/admin") // We can access '/admin'


g := r.Group("/hello")

g.Use(
	middleware1,
	middleware2,
	...
)
g.Get("/world") // '/hello/world'

*/

func (g *Group) UseMiddleWare(middleware ...MiddleWareFunc) {
	g.middleware = append(g.middleware, middleware...)
}

func (g *Group) Handle(method, suffix string, handler HandlerFunc) {
	p := path.Join(g.prefix, suffix)
	g.base.Handle(method, p, func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		h := handler
		// Chain middleware
		for i := len(g.middleware) - 1; i >= 0; i-- {
			h = g.middleware[i](h)
		}
		g.ctx.ResetContext(w, req, p)
		if err := h(g.ctx); err != nil {
			g.HTTPErrorHandler(err, g.ctx)
		}
	})
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (g *Group) GET(path string, handle HandlerFunc) {
	g.Handle("GET", path, handle)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (g *Group) HEAD(path string, handle HandlerFunc) {
	g.Handle("HEAD", path, handle)
}

// OPTION shortcut for router.Handle("OPTIONS", path, handle)
func (g *Group) OPTIONS(path string, handle HandlerFunc) {
	g.Handle("OPTIONS", path, handle)
}

// POST shortcut for router.Handle("POST", path, handle)
func (g *Group) POST(path string, handle HandlerFunc) {
	g.Handle("POST", path, handle)
}

// PUT shortcut for router.Handle("PUT", path, handle)
func (g *Group) PUT(path string, handle HandlerFunc) {
	g.Handle("PUT", path, handle)
}

// PATCH shortcut for router.Handle("PATCH", path, handle)
func (g *Group) PATCH(path string, handle HandlerFunc) {
	g.Handle("PATCH", path, handle)
}

// DELETE shortcut for router.Handle("DELETE", path, handle)
func (g *Group) DELETE(path string, handle HandlerFunc) {
	g.Handle("DELETE", path, handle)
}
