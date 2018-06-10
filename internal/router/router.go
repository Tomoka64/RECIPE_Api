package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type Router struct {
	base             *httprouter.Router
	middleware       []MiddleWareFunc
	Logger           *zap.Logger
	HTTPErrorHandler func(error, http.ResponseWriter, *http.Request)
}

type (
	MiddleWareFunc func(HandlerFunc) HandlerFunc
	HandlerFunc    func(http.ResponseWriter, *http.Request, httprouter.Params) error
)

var (
	// NotFoundHandler run if not match registerd method and path.
	NotFoundHandler = func(http.ResponseWriter, *http.Request, httprouter.Params) error {
		return ErrNotFound
	}

	// MethodNotAllowedHandler run if not match registerd method.
	MethodNotAllowedHandler = func(http.ResponseWriter, *http.Request, httprouter.Params) error {
		return ErrMethodNotAllowed
	}
)

func New(l *zap.Logger) *Router {
	r := &Router{
		base:   httprouter.New(),
		Logger: l,
	}
	r.HTTPErrorHandler = r.DefaultHTTPErrorHandler
	return r
}

func (r *Router) UseMiddleWare(middleware ...MiddleWareFunc) {
	r.middleware = append(r.middleware, middleware...)
}

var methods = []string{
	"DELETE",
	"GET",
	"HEAD",
	"OPTIONS",
	"PATCH",
	"POST",
	"PUT",
}

// NOTE: slow point
func (r *Router) checkMethodNotAllowed(path string) HandlerFunc {
	for _, m := range methods {
		h, _, _ := r.base.Lookup(m, path)
		if h != nil {
			return MethodNotAllowedHandler
		}
	}
	return NotFoundHandler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.base.ServeHTTP(w, req)
}

func (r *Router) HandlerFunc(method, path string, handler http.HandlerFunc) {
	r.base.Handler(method, path, handler)
}

func (r *Router) Handle(method, path string, handler HandlerFunc) {
	r.base.Handle(method, path, func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		h := handler
		// Chain middleware
		for i := len(r.middleware) - 1; i >= 0; i-- {
			h = r.middleware[i](h)
		}
		if err := h(w, req, p); err != nil {
			r.HTTPErrorHandler(err, w, req)
		}
	})
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (r *Router) GET(path string, handle httprouter.Handle) {
	r.base.Handle("GET", path, handle)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (r *Router) HEAD(path string, handle httprouter.Handle) {
	r.base.Handle("HEAD", path, handle)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (r *Router) OPTIONS(path string, handle httprouter.Handle) {
	r.base.Handle("OPTIONS", path, handle)
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (r *Router) POST(path string, handle httprouter.Handle) {
	r.base.Handle("POST", path, handle)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (r *Router) PUT(path string, handle httprouter.Handle) {
	r.base.Handle("PUT", path, handle)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (r *Router) PATCH(path string, handle httprouter.Handle) {
	r.base.Handle("PATCH", path, handle)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (r *Router) DELETE(path string, handle httprouter.Handle) {
	r.base.Handle("DELETE", path, handle)
}
