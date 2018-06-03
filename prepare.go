package main

import (
	"github.com/Tomoka64/RECIPE_Api/internal/postgres"
	"github.com/Tomoka64/RECIPE_Api/internal/redis"
	"github.com/alecthomas/template"
	"github.com/jmoiron/sqlx"
	redistore "gopkg.in/boj/redistore.v1"
)

type Server struct {
	DB         *sqlx.DB
	Store      *redistore.RediStore
	StackTrace bool
	Tpl        *template.Template
}

func New() *Server {
	db, err := postgres.New()
	tpl := template.Must(template.ParseFiles("templates/*"))
	if err != nil {
		panic(err)
	}
	store, err := redis.New()
	if err != nil {
		panic(err)
	}
	return &Server{
		DB:         db,
		Store:      store,
		StackTrace: false,
		Tpl:        tpl,
	}
}
