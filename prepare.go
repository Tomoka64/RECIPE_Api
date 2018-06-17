package main

import (
	"github.com/Tomoka64/RECIPE_Api/internal/postgres"
	"github.com/Tomoka64/RECIPE_Api/internal/redis"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	redistore "gopkg.in/boj/redistore.v1"
)

type Server struct {
	DB     *sqlx.DB
	Store  *redistore.RediStore
	Logger *zap.Logger
}

func New() *Server {
	db, err := postgres.New()
	if err != nil {
		panic(err)
	}
	store, err := redis.New()
	if err != nil {
		panic(err)
	}
	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return &Server{
		DB:     db,
		Store:  store,
		Logger: l,
	}
}
