package main

import (
	"html/template"

	"github.com/Tomoka64/RECIPE_Api/internal/postgres"
	"github.com/Tomoka64/RECIPE_Api/internal/redis"
	"github.com/Tomoka64/RECIPE_Api/logger"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	redistore "gopkg.in/boj/redistore.v1"
)

type Server struct {
	DB     *sqlx.DB
	Store  *redistore.RediStore
	Tpl    *template.Template
	Logger *zap.Logger
}

func New() *Server {
	tpl := template.Must(template.ParseFiles("templates/*"))
	db, err := postgres.New()
	if err != nil {
		panic(err)
	}
	store, err := redis.New()
	if err != nil {
		panic(err)
	}
	l, err := logger.New()
	if err != nil {
		panic(err)
	}
	return &Server{
		DB:     db,
		Store:  store,
		Tpl:    tpl,
		Logger: l,
	}
}
