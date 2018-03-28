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
		if err == model.ErrUnknownProvider {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get an authorization URL: %v", err))
	}
	if redirectURL := c.Request().Header.Get("Referer"); redirectURL != "" {
		redirectURLCookie := &http.Cookie{
			Name:  "redirect_url",
			Value: redirectURL,
			Path:  "/",
		}
		c.SetCookie(redirectURLCookie)
	}
	return c.Redirect(http.StatusFound, authoURL.String())
}

//AuthCallback the Method Handler of "GET /auth/:provider/callback"
func AuthCallback(c echo.Context) error {
	provider := c.Param("provider")

	query := c.Request().URL.Query()
	id, err := model.GetAuthedUserID(provider, &query)
	if err != nil {
		if err == model.ErrUnknownProvider {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the authenticated user's ID: %v", err))
	}
	user, err := model.GetUserByID(id, true)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the user record: %v", err))
	}
	if err := user.SetToken(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to set a token: %v", err))
	}
	redirectURL := os.Getenv("AUTH_CALLBACK_REDIRECT_URL")
	if cookie, err := c.Cookie("redirect_url"); err == nil {
		redirectURL = cookie.Value
		redirectURLCookie := &http.Cookie{
			Name:   "redirect_url",
			Value:  "dummy",
			MaxAge: -114514,
			Path:   "/",
		}
		c.SetCookie(redirectURLCookie)
	}
	tokenCookie := &http.Cookie{
		Name:    "token",
		Value:   user.Token,
		Expires: user.TokenExpires,
		Path:    "/",
	}
	c.SetCookie(tokenCookie)
	return c.Redirect(http.StatusFound, redirectURL)
}
