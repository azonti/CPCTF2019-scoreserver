package router

import (
	"fmt"
	"git.trapti.tech/CPCTF2018/scoreserver/model"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type challengeJSON struct {
	ID        string      `json:"id"`
	Author    *userJSON   `json:"author"`
	Score     int         `json:"score"`
	Caption   string      `json:"caption"`
	Hints     []*hintJSON `json:"hints"`
	Flag      string      `json:"flag"`
	Answer    string      `json:"answer"`
	WhoSolved []*userJSON `json:"who_solved"`
}

type hintJSON struct {
	ID      string `json:"id"`
	Caption string `json:"caption"`
	Penalty int    `json:"penalty"`
}

func contains(slice []string, x string) bool {
	for _, y := range slice {
		if x == y {
			return true
		}
	}
	return false
}

func newChallengeJSON(me *model.User, challenge *model.Challenge) (*challengeJSON, error) {
	json := &challengeJSON{
		ID:      challenge.ID,
		Score:   challenge.Score,
		Caption: challenge.Caption,
	}
	author, err := model.GetUserByID(challenge.AuthorID, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get the author record: %v", err)
	}
	json.Author = newUserJSON(author)
	for i := 0; i < len(challenge.Hints); i++ {
		if me.IsAuthor || contains(challenge.WhoSolvedIDs, me.ID) || contains(me.OpenedHintIDs, challenge.Hints[i].ID) {
			json.Hints = append(json.Hints, &hintJSON{ID: challenge.Hints[i].ID, Caption: challenge.Hints[i].Caption, Penalty: challenge.Hints[i].Penalty})
		}
	}
	if me.IsAuthor || contains(challenge.WhoSolvedIDs, me.ID) {
		json.Flag, json.Answer = challenge.Flag, challenge.Answer
	}
	for _, whoSolvedID := range challenge.WhoSolvedIDs {
		whoSolved, err := model.GetUserByID(whoSolvedID, false)
		if err != nil {
			return nil, fmt.Errorf("failed to get who solved record: %v", err)
		}
		json.WhoSolved = append(json.WhoSolved, newUserJSON(whoSolved))
	}
	return json, nil
}

//EnsureContestStarted Ensure the Contest has Started
func EnsureContestStarted(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		now, start := time.Now(), model.StartTime()
		me := c.Get("me").(*model.User)
		if start.After(now) && !me.IsAuthor {
			return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("the contest has not start yet"))
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

//GetChallenges the Method Handler of "GET /challenges"
func GetChallenges(c echo.Context) error {
	challenges, err := model.GetChallenges()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	me := c.Get("me").(*model.User)
	jsons := make([]*challengeJSON, len(challenges))
	for i := 0; i < len(challenges); i++ {
		json, err := newChallengeJSON(me, challenges[i])
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to parse the challenge record: %v", err))
		}
		jsons[i] = json
	}
	return c.JSON(http.StatusOK, jsons)
}

//GetChallenge the Method Handler of "GET /challenges/:challengeID"
func GetChallenge(c echo.Context) error {
	challengeID := c.Param("challengeID")
	challenge, err := model.GetChallengeByID(challengeID)
	if err != nil {
		if err == model.ErrChallengeNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	me := c.Get("me").(*model.User)
	json, err := newChallengeJSON(me, challenge)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to parse the challenge record: %v", err))
	}
	return c.JSON(http.StatusOK, json)
}

//PostChallenge the Method Handler of "POST /challenges"
func PostChallenge(c echo.Context) error {
	req := &challengeJSON{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to bind request body: %v", err))
	}
	captions, penalties := make([]string, len(req.Hints)), make([]int, len(req.Hints))
	for _, json := range req.Hints {
		idSplit := strings.Split(json.ID, ":")
		i, _ := strconv.Atoi(idSplit[1])
		captions[i] = json.Caption
		penalties[i] = json.Penalty
	}
	challenge, err := model.NewChallenge(req.Author.ID, req.Score, req.Caption, captions, penalties, req.Flag, req.Answer)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	me := c.Get("me").(*model.User)
	json, _ := newChallengeJSON(me, challenge)
	c.Response().Header().Set(echo.HeaderLocation, "/challenges/"+challenge.ID)
	return c.JSON(http.StatusCreated, json)
}

//PutChallenge the Method Handler of "PUT /challenges/:challengeID"
func PutChallenge(c echo.Context) error {
	challengeID := c.Param("challengeID")
	challenge, err := model.GetChallengeByID(challengeID)
	if err != nil {
		if err == model.ErrChallengeNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the challenge record: %v", err))
	}
	req := &challengeJSON{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to bind request body: %v", err))
	}
	captions, penalties := make([]string, len(req.Hints)), make([]int, len(req.Hints))
	for _, json := range req.Hints {
		idSplit := strings.Split(json.ID, ":")
		i, _ := strconv.Atoi(idSplit[1])
		captions[i] = json.Caption
		penalties[i] = json.Penalty
	}
	if err := challenge.Update(req.Author.ID, req.Score, req.Caption, captions, penalties, req.Flag, req.Answer); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

//DeleteChallenge the Method Handler of "DELETE /challenges/:challengeID"
func DeleteChallenge(c echo.Context) error {
	challengeID := c.Param("challengeID")
	challenge, err := model.GetChallengeByID(challengeID)
	if err != nil {
		if err == model.ErrChallengeNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the challenge record: %v", err))
	}
	if err := challenge.Delete(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

//CheckAnswer the Method Handler of "POST /challenges/:challengeID"
func CheckAnswer(c echo.Context) error {
	challengeID := c.Param("challengeID")
	challenge, err := model.GetChallengeByID(challengeID)
	if err != nil {
		if err == model.ErrChallengeNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get the challenge record: %v", err))
	}
	req := &struct {
		Flag string `form:"flag"`
	}{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to bind request body: %v", err))
	}
	if challenge.Flag != req.Flag {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("the flag is wrong"))
	}
	me := c.Get("me").(*model.User)
	if err := challenge.AddWhoSolved(me); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to add you to the list of who solved: %v", err))
	}
	return c.Redirect(http.StatusSeeOther, "/challenges/"+challengeID)
}
