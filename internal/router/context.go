package router

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/Tomoka64/RECIPE_Api/internal/header"
	"github.com/Tomoka64/RECIPE_Api/internal/mime"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

type Context interface {
	Get(interface{}) (interface{}, bool)
	Set(interface{}, interface{})
	Request() *http.Request
	Response() http.ResponseWriter
	Params() httprouter.Params
	ResetContext(http.ResponseWriter, *http.Request, httprouter.Params)
	SetContentType(code int, contentType string)
	NoContent(code int) error
	Redirect(code int, url string) error
	JSON(code int, i interface{}) error
	String(code int, s string) (err error)
}

type ctx struct {
	m        sync.Map
	request  *http.Request
	response http.ResponseWriter
	params   httprouter.Params
}

func NewContext() Context {
	return &ctx{}
}

func (c *ctx) ResetContext(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	c.response = res
	c.request = req
	c.params = params
}

func (c *ctx) Get(key interface{}) (interface{}, bool) {
	return c.m.Load(key)
}

func (c *ctx) Set(key, val interface{}) {
	c.m.Store(key, val)
}

func (c *ctx) Request() *http.Request {
	return c.request
}

func (c *ctx) Response() http.ResponseWriter {
	return c.response
}

func (c *ctx) Params() httprouter.Params {
	return c.params
}

func (c *ctx) SetContentType(code int, contentType string) {
	c.response.Header().Set(header.ContentType, contentType)
	c.response.WriteHeader(code)
}

func (c *ctx) NoContent(code int) error {
	c.response.WriteHeader(code)
	return nil
}

func (c *ctx) Redirect(code int, url string) error {
	if code < 300 || code > 308 {
		return errors.New("Invalid redirect status code")
	}
	c.response.Header().Set(header.Location, url)
	c.response.WriteHeader(code)
	return nil
}

func (c *ctx) JSON(code int, i interface{}) error {
	c.SetContentType(code, mime.ApplicationJSONCharsetUTF8)
	return json.NewEncoder(c.response).Encode(i)
}

func (c *ctx) String(code int, s string) (err error) {
	c.SetContentType(code, mime.ApplicationJSONCharsetUTF8)
	_, err = c.response.Write([]byte(s))
	return
}
