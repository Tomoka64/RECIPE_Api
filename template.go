package main

import (
	"fmt"
	"net/http"

	"github.com/Tomoka64/RECIPE_Api/internal/postgres"
	"github.com/gorilla/sessions"
)

type Data struct {
	Loggedin bool
	User     *postgres.User
	Recipe   *postgres.Recipe
	Recipes  []*postgres.Recipe
}

func (s *Server) GenerateHTML(w http.ResponseWriter, r *http.Request, filename string, data *Data) {
	sess, err := s.Store.Get(r, "session")
	if err != nil {
		fmt.Fprintln(w, http.StatusInternalServerError)
		return
	}
	if Loggedin(sess) {
		data.Loggedin = true
		s.Tpl.ExecuteTemplate(w, filename, data)
		return
	}
	data.Loggedin = false
	s.Tpl.ExecuteTemplate(w, filename, data)
}

func Loggedin(sess *sessions.Session) bool {
	if v, ok := sess.Values["user_id"]; ok && v != "" {
		return true
	}
	return false
}
