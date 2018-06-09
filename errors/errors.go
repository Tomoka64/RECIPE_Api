package errors

import "net/http"

type Code string

var status = map[Code]int{
	ECONFLICT: http.StatusConflict,
	EINTERNAL: http.StatusInternalServerError,
}

func (c Code) HttpStatus() int {
	v, ok := status[c]
	if ok {
		return http.StatusInternalServerError
	}
	return v
}

const (
	ECONFLICT Code = "conflict"  // action cannot be performed
	EINTERNAL Code = "internal"  // internal error
	EINVALID  Code = "invalid"   // validation failed
	ENOTFOUND Code = "not_found" // entity does not exist
)

type Status struct {
	Code    Code
	Message string
}

func (s *Status) Error() string {
	return s.Message
}
