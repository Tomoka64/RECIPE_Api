package router

import (
	"fmt"
	"net/http"

	"github.com/Tomoka64/RECIPE_Api/internal/status"
	"github.com/Tomoka64/RECIPE_Api/internal/writer"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Errors
var (
	ErrUnsupportedMediaType        = NewHTTPError(status.UnsupportedMediaType)
	ErrNotFound                    = NewHTTPError(status.NotFound)
	ErrUnauthorized                = NewHTTPError(status.Unauthorized)
	ErrForbidden                   = NewHTTPError(status.Forbidden)
	ErrMethodNotAllowed            = NewHTTPError(status.MethodNotAllowed)
	ErrStatusRequestEntityTooLarge = NewHTTPError(status.RequestEntityTooLarge)
	ErrInternalServer              = NewHTTPError(status.InternalServerError)
	ErrValidatorNotRegistered      = errors.New("Validator not registered")
	ErrRendererNotRegistered       = errors.New("Renderer not registered")
	ErrInvalidRedirectCode         = errors.New("Invalid redirect status code")
	ErrCookieNotFound              = errors.New("Cookie not found")
)

// HTTPError represents an error that occurred while handling a request.
type HTTPError struct {
	Code    int
	Message string
}

// NewHTTPError creates a new HTTPError instance.
func NewHTTPError(code int, message ...string) *HTTPError {
	he := &HTTPError{Code: code, Message: http.StatusText(code)}
	if len(message) > 0 {
		he.Message = message[0]
	}
	return he
}

// Error makes it compatible with `error` interface.
func (he *HTTPError) Error() string {
	return fmt.Sprintf("code=%d, message=%v", he.Code, he.Message)
}

// DefaultHTTPErrorHandler render error message to response
func (r *Router) DefaultHTTPErrorHandler(err error, w http.ResponseWriter, req *http.Request) {
	var (
		code = status.InternalServerError
		msg  string
	)

	if he, ok := err.(*HTTPError); ok {
		code = he.Code
		msg = he.Message
	} else {
		msg = http.StatusText(code)
	}

	r.Logger.Error("Catching error", zap.Error(err))

	if req.Method == "HEAD" {
		err = writer.NoContent(w, code)
	} else {
		err = writer.String(w, code, msg)
	}
	if err != nil {
		r.Logger.Error("Send response error", zap.Error(err))
	}
}
