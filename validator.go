package main

import (
	"github.com/Tomoka64/RECIPE_Api/internal/postgres"
)

func (s *Server) UserAndPasswordMatch(name, password string) bool {
	_, err := postgres.GetUser(s.DB, name, password)
	return err == nil
}

func (s *Server) UserExists(username, email string) bool {
	return (postgres.UserExistsByUsername(s.DB, username) || postgres.UserExistsByEmail(s.DB, email))
}

func IsPasswordValid(password string) bool {
	if len(password) < 6 {
		return false
	}
	if len(password) > 20 {
		return false
	}
	//TODO: もしパスワードにアルファベット以外含まれていたら..
	return true
}
