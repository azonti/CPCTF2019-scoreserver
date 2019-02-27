package router

import (
	"fmt"
	"net/http"
	"time"

	"git.trapti.tech/CPCTF2019/scoreserver/model"
	"github.com/labstack/echo"
)

//EnsureContestStarted Ensure the Contest has Started
func EnsureContestStarted(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		now, start := time.Now(), model.StartTime()
		me := c.Get("me").(*model.User)
		if start.After(now) && !me.IsAuthor {
			c.Response().Header().Set("Retry-After", start.UTC().Format(http.TimeFormat))
			return echo.NewHTTPError(http.StatusServiceUnavailable, fmt.Sprintf("the contest has not started yet"))
		}
		return next(c)
	}
}

//EnsureContestNotFinished Ensure the Contest has not finished yet
func EnsureContestNotFinished(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		now, finish := time.Now(), model.FinishTime()
		me := c.Get("me").(*model.User)
		if !finish.After(now) && !me.IsAuthor {
			return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("the contest has finished"))
		}
		return next(c)
	}
}
