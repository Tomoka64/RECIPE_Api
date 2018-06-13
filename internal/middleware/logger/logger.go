package logger

import (
	"github.com/Tomoka64/RECIPE_Api/internal/router"
	"go.uber.org/zap"
)

func MiddleWare(l *zap.Logger) router.MiddleWareFunc {
	return func(h router.HandlerFunc) router.HandlerFunc {
		return func(c router.Context) error {
			if err := h(c); err != nil {
				l.Error("error found ", zap.Error(err))
				return err
			}
			return nil
		}
	}
}
