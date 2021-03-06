package router

import (
	"fmt"
	"git.trapti.tech/CPCTF2019/scoreserver/model"
	"github.com/labstack/echo"
	"net/http"
)

//Auth the Method Handler of "GET /auth/:provider"
func Auth(c echo.Context) error {
	provider := c.Param("provider")

	authnURL, err := model.GetAuthnURL(provider)
	if err != nil {
		if err == model.ErrUnknownProvider {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get an authentication URL: %v", err))
	}

	if redirectURL := c.Request().Header.Get("Referer"); redirectURL != "" {
		redirectURLCookie := &http.Cookie{
			Name:  "redirect_url",
			Value: redirectURL,
			Path:  "/",
		}
		c.SetCookie(redirectURLCookie)
	}

	return c.Redirect(http.StatusFound, authnURL.String())
}

//AuthCallback the Method Handler of "GET /auth/:provider/callback"
func AuthCallback(c echo.Context) error {
	provider := c.Param("provider")

	query := c.Request().URL.Query()
	id, err := model.GetAuthnedUserID(provider, &query)
	if err != nil && err != model.ErrDeniedAuthn {
		if err == model.ErrUnknownProvider {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the authenticated user's ID: %v", err))
	}

	if err != model.ErrDeniedAuthn {
		user, err := model.GetUserByID(id, true)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the user record: %v", err))
		}

		if err := user.PutToken(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to set a token: %v", err))
		}

		tokenCookie := &http.Cookie{
			Name:    "token",
			Value:   user.Token,
			Expires: user.TokenExpires,
			Path:    "/",
		}
		c.SetCookie(tokenCookie)
	}

	redirectURL := "/"
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

	return c.Redirect(http.StatusFound, redirectURL)
}

//Logout the Method Handler of "GET /logout"
func Logout(c echo.Context) error {
	me := c.Get("me").(*model.User)

	if err := me.DeleteToken(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to remove the token: %v", err))
	}

	tokenCookie := &http.Cookie{
		Name:   "token",
		Value:  "dummy",
		MaxAge: -114514,
		Path:   "/",
	}
	c.SetCookie(tokenCookie)

	redirectURL := c.Request().Header.Get("Referer")
	if redirectURL == "" {
		redirectURL = "/"
	}

	return c.Redirect(http.StatusFound, redirectURL)
}
