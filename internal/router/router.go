package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type Router struct {
	base             *httprouter.Router
	ctx              Context
	middleware       []MiddleWareFunc
	Logger           *zap.Logger
	HTTPErrorHandler func(error, Context)
}

type (
	MiddleWareFunc func(HandlerFunc) HandlerFunc
	HandlerFunc    func(Context) error
)

var (
	// NotFoundHandler run if not match registerd method and path.
	NotFoundHandler = func(Context) error {
		return ErrNotFound
	}

	// MethodNotAllowedHandler run if not match registerd method.
	MethodNotAllowedHandler = func(Context) error {
		return ErrMethodNotAllowed
	}
)

func New(l *zap.Logger) *Router {
	r := &Router{
		base:   httprouter.New(),
		ctx:    NewContext(),
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
		r.ctx.ResetContext(w, req, p)
		if err := h(r.ctx); err != nil {
			r.HTTPErrorHandler(err, r.ctx)
		}
	})
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (r *Router) GET(path string, handle HandlerFunc) {
	r.Handle("GET", path, handle)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (r *Router) HEAD(path string, handle HandlerFunc) {
	r.Handle("HEAD", path, handle)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (r *Router) OPTIONS(path string, handle HandlerFunc) {
	r.Handle("OPTIONS", path, handle)
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (r *Router) POST(path string, handle HandlerFunc) {
	r.Handle("POST", path, handle)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (r *Router) PUT(path string, handle HandlerFunc) {
	r.Handle("PUT", path, handle)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (r *Router) PATCH(path string, handle HandlerFunc) {
	r.Handle("PATCH", path, handle)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (r *Router) DELETE(path string, handle HandlerFunc) {
	r.Handle("DELETE", path, handle)
}
