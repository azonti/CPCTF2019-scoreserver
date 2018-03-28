package router

import (
	"fmt"
	"git.trapti.tech/CPCTF2018/scoreserver/model"
	"github.com/labstack/echo"
	"net/http"
)

type userJSON struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	IconURL           string `json:"icon_url"`
	TwitterScreenName string `json:"twitter_screen_name"`
	IsAuthor          bool   `json:"is_author"`
}

func newUserJSON(user *model.User) *userJSON {
	json := &userJSON{
		ID:                user.ID,
		Name:              user.Name,
		IconURL:           user.IconURL,
		TwitterScreenName: user.TwitterScreenName,
		IsAuthor:          user.IsAuthor,
	}
	return json
}

//DetermineMe Determine Who am I
func DetermineMe(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				c.Set("me", model.Nobody)
				return next(c)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the token cookie: %v", err))
		}
		me, err := model.GetUserByToken(cookie.Value)
		if err != nil {
			if err == model.ErrUserNotFound {
				c.Set("me", model.Nobody)
				return next(c)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get me: %v", err))
		}
		c.Set("me", me)
		return next(c)
	}
}

//EnsureIExist Ensure I Exist
func EnsureIExist(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		me := c.Get("me").(*model.User)
		if me.ID == model.Nobody.ID {
			return echo.NewHTTPError(http.StatusForbidden, "please log in")
		}
		return next(c)
	}
}

//EnsureIAmAuthor Ensure I am an Author
func EnsureIAmAuthor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		me := c.Get("me").(*model.User)
		if !me.IsAuthor {
			return echo.NewHTTPError(http.StatusForbidden, "you are not an author")
		}
		return next(c)
	}
}
