package mongodb

import (
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
)

func MiddlewareMongoDB(sessionkey string, session *mgo.Session) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		s := session.Clone()
		// defer s.Close()
		return func(c echo.Context) error {
			c.Set(sessionkey, s.DB(""))

			return next(c)
		}
	}
}
