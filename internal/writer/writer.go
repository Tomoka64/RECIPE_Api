package writer

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Tomoka64/RECIPE_Api/internal/header"
	"github.com/Tomoka64/RECIPE_Api/internal/mime"
)

func SetContentType(w http.ResponseWriter, code int, contentType string) {
	w.Header().Set(header.ContentType, contentType)
	w.WriteHeader(code)
}

func NoContent(w http.ResponseWriter, code int) error {
	w.WriteHeader(code)
	return nil
}

func Redirect(w http.ResponseWriter, code int, url string) error {
	if code < 300 || code > 308 {
		return errors.New("Invalid redirect status code")
	}
	w.Header().Set(header.Location, url)
	w.WriteHeader(code)
	return nil
}

func JSON(w http.ResponseWriter, code int, i interface{}) error {
	SetContentType(w, code, mime.ApplicationJSONCharsetUTF8)
	return json.NewEncoder(w).Encode(i)
}

func String(w http.ResponseWriter, code int, s string) (err error) {
	SetContentType(w, code, mime.ApplicationJSONCharsetUTF8)
	_, err = w.Write([]byte(s))
	return
}
