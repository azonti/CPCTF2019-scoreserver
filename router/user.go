package router

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"git.trapti.tech/CPCTF2019/scoreserver/model"
	"github.com/labstack/echo"
)

type userJSON struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	IconURL           string `json:"icon_url"`
	TwitterScreenName string `json:"twitter_screen_name"`
	IsAuthor          bool   `json:"is_author"`
	IsOnsite          bool   `json:"is_onsite"`
	Score             int    `json:"score"`
	WebShellPass      string `json:"web_shell_pass"`
}

func newUserJSON(me *model.User, user *model.User) *userJSON {
	canISeePass := me.ID == user.ID || me.IsAuthor
	json := &userJSON{
		ID:                user.ID,
		Name:              user.Name,
		IconURL:           user.IconURL,
		TwitterScreenName: user.TwitterScreenName,
		IsAuthor:          user.IsAuthor,
		IsOnsite:          user.IsOnsite,
		Score:             user.Score,
		WebShellPass:      map[bool]string{true: user.WebShellPass}[canISeePass],
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

//EnsureINotExist Ensure I do Not Exist
func EnsureINotExist(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		me := c.Get("me").(*model.User)
		if me.ID != model.Nobody.ID {
			return echo.NewHTTPError(http.StatusForbidden, "you have already logged in")
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

//GetUsers the Method Handler of "GET /users"
func GetUsers(c echo.Context) error {
	users, err := model.GetUsers()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	me := c.Get("me").(*model.User)
	jsons := make([]*userJSON, len(users))
	for i := 0; i < len(users); i++ {
		jsons[i] = newUserJSON(me, users[i])
	}

	return c.JSON(http.StatusOK, jsons)
}

//GetUser the Method Handler of "GET /users/:userID"
func GetUser(c echo.Context) error {
	userID := c.Param("userID")
	me := c.Get("me").(*model.User)

	user, err := model.GetUserByID(userID, false)
	if err != nil {
		if err == model.ErrUserNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	json := newUserJSON(me, user)

	return c.JSON(http.StatusOK, json)
}

//GetMe the Method Handler of "GET /users/me"
func GetMe(c echo.Context) error {
	me := c.Get("me").(*model.User)

	json := newUserJSON(me, me)

	return c.JSON(http.StatusOK, json)
}

//CheckCode the Method Handler of "POST /users/me"
func CheckCode(c echo.Context) error {
	me := c.Get("me").(*model.User)

	req := &struct {
		Code string `form:"code"`
	}{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to bind request body: %v", err))
	}

	now, finish := time.Now(), model.FinishTime()

	switch req.Code {
	case os.Getenv("AUTHOR_CODE"):
		if err := me.MakeMeAuthor(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusNoContent)
	case os.Getenv("ONSITE_CODE"):
		if !finish.After(now) && !me.IsAuthor {
			return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("the contest has finished"))
		}
		if err := me.MakeMeOnsite(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusNoContent)
	case "rwsc":
		if !me.IsOnsite {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("you are not onsite"))
		}
		if err := me.RecreateWebShellContainer(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusNoContent)
	}
	switch {
	case strings.HasPrefix(req.Code, "hint:"):
		if !finish.After(now) && !me.IsAuthor {
			return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("the contest has finished"))
		}
		partedCode := strings.Split(req.Code, ":")
		if len(partedCode) != 3 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid hint code"))
		}
		hintID := strings.Join(partedCode[1:], ":")
		cnt := 0
		for _, openedHint := range me.OpenedHints {
			if hintID == openedHint.ID {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("already opened"))
			}
			if strings.HasPrefix(openedHint.ID, partedCode[1]+":") {
				cnt++
			}
		}
		if strconv.Itoa(cnt) != partedCode[2] {
			return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("you cannot open this hint yet"))
		}
		if err := me.OpenHint(hintID); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusNoContent)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("the code is wrong"))
	}
}

//GetSolvedChallenges the Method Handler of "GET /users/:userID/solved"
func GetSolvedChallenges(c echo.Context) error {
	userID := c.Param("userID")
	me := c.Get("me").(*model.User)

	user, err := model.GetUserByID(userID, false)
	if err != nil {
		if err == model.ErrUserNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the user record: %v", err))
	}

	solvedMap, openedMap, foundMap := makeSolvedOpenedFoundMaps(me)
	jsons := make([]*challengeJSON, len(user.SolvedChallenges))
	for i, challenge := range user.SolvedChallenges {
		jsons[i] = newChallengeJSON(me, challenge, solvedMap, openedMap, foundMap)
	}

	return c.JSON(http.StatusOK, jsons)
}

//GetLastSolvedChallenge the Method Handler of "GET /user/:userID/solved/last"
func GetLastSolvedChallenge(c echo.Context) error {
	userID := c.Param("userID")
	me := c.Get("me").(*model.User)

	user, err := model.GetUserByID(userID, false)
	if err != nil {
		if err == model.ErrUserNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the user record: %v", err))
	}

	if user.LastSolvedChallengeID == "" {
		return c.NoContent(http.StatusNoContent)
	}

	solvedMap, openedMap, foundMap := makeSolvedOpenedFoundMaps(me)
	json := newChallengeJSON(me, user.LastSolvedChallenge, solvedMap, openedMap, foundMap)

	c.Response().Header().Set(echo.HeaderLastModified, user.LastSolvedTime.UTC().Format(http.TimeFormat))
	return c.JSON(http.StatusOK, json)
}

//GetLastSeenChallenge the Method Handler of "GET /user/:userID/lastseen"
func GetLastSeenChallenge(c echo.Context) error {
	userID := c.Param("userID")
	me := c.Get("me").(*model.User)

	user, err := model.GetUserByID(userID, false)
	if err != nil {
		if err == model.ErrUserNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the user record: %v", err))
	}

	if user.LastSeenChallengeID == "" {
		return c.NoContent(http.StatusNoContent)
	}

	solvedMap, openedMap, foundMap := makeSolvedOpenedFoundMaps(me)
	json := newChallengeJSON(me, user.LastSeenChallenge, solvedMap, openedMap, foundMap)

	return c.JSON(http.StatusOK, json)
}
