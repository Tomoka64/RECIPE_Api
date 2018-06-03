package main

import (
	"github.com/Tomoka64/RECIPE_Api/internal/postgres"
)

func (s *Server) UserAndPasswordMatch(name, password string) bool {
	_, err := postgres.GetUser(s.DB, name, password)
	return err == nil
}

func (s *Server) UserExists(arg string) bool {
	return (postgres.UserExistsByUsername(s.DB, arg) || postgres.UserExistsByEmail(s.DB, arg))
}
