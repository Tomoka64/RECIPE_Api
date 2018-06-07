package redis

import (
	"os"

	redistore "gopkg.in/boj/redistore.v1"
)

type Store = *redistore.RediStore

func New() (Store, error) {
	store, err := redistore.NewRediStore(10, "tcp", ":"+os.Getenv("REDIS_PORT"), "", []byte(os.Getenv("REDIS_SECRETKEY")))
	if err != nil {
		return nil, err
	}
	store.SetMaxAge(10 * 24 * 3600) //30days
	return store, nil
}
