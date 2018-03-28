package router

import (
	"fmt"
	"git.trapti.tech/CPCTF2018/scoreserver/model"
	"github.com/labstack/echo"
	"net/http"
	"os"
)

//Auth the Method Handler of "GET /auth/:provider"
func Auth(c echo.Context) error {
	provider := c.Param("provider")
	authoURL, err := model.GetAuthoURL(provider)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get an authorization URL: %v", err))
	}
	return c.Redirect(http.StatusFound, authoURL.String())
}

//AuthCallback the Method Handler of "GET /auth/:provider/callback"
func AuthCallback(c echo.Context) error {
	provider := c.Param("provider")

	query := c.Request().URL.Query()
	id, err := model.GetAuthedUserID(provider, &query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the authenticated user's ID: %v", err))
	}
	user, err := model.GetUserByID(provider, id, true)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the user record: %v", err))
	}
	if err := user.SetToken(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to set a token: %v", err))
	}
	cookie := &http.Cookie{
		Name:    "token",
		Value:   user.Token,
		Expires: user.TokenExpires,
		Path:    "/",
	}
	c.SetCookie(cookie)
	return c.Redirect(http.StatusFound, os.Getenv("AUTH_CALLBACK_REDIRECT_URL"))
}
