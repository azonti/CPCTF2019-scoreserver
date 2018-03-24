package router

import (
	"git.trapti.tech/CPCTF2018/scoreserver/model"
	"github.com/labstack/echo"
	"net/http"
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

//EnsureIExist Ensure I Exist
func EnsureIExist(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		me := c.Get("me").(*model.User)
		if me == nil {
			return echo.NewHTTPError(http.StatusForbidden, "please log in")
		}
		return next(c)
	}
}
