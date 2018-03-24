package router

import (
	"git.trapti.tech/CPCTF2018/scoreserver/model"
	"github.com/labstack/echo"
)

//DetermineMe Determine Who am I
func DetermineMe(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil {
			return next(c)
		}
		me, err := model.GetUserByToken(cookie.Value)
		if err != nil {
			return next(c)
		}
		c.Set("me", me)
		return next(c)
	}
}
