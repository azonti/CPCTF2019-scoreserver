package router

import (
	"fmt"
	"git.trapti.tech/CPCTF2018/scoreserver/model"
	"github.com/labstack/echo"
	"net/http"
)

//Auth Method Handler of "GET /auth/:provider"
func Auth(c echo.Context) error {
	provider := c.Param("provider")
	authURL, err := model.GetAuthURL(provider)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get authorization URL: %v", err))
	}
	return c.Redirect(http.StatusFound, authURL.String())
}

//AuthCallback Method Handler of "GET /auth/:provider/callback"
func AuthCallback(c echo.Context) error {
	provider := c.Param("provider")

	token, err := model.Login(provider, c.Request().URL.Query())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to login: %v", err))
	}

	return c.String(http.StatusOK, token)
}
