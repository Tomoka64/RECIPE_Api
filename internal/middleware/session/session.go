package session

import (
	"github.com/Tomoka64/RECIPE_Api/internal/router"
	"github.com/gorilla/sessions"
)

var SessionKey = "session_key"

func MiddleWare(key string, store sessions.Store) router.MiddleWareFunc {
	return func(h router.HandlerFunc) router.HandlerFunc {
		return func(c router.Context) error {
			sess, err := store.Get(c.Request(), SessionKey)
			if err != nil {
				return err
			}
			sess.Options.Secure = true
			sess.Options.HttpOnly = true
			c.Set(key, sess)
			return h(c)
		}
	}
}
